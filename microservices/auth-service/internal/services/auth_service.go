package services

import (
    "context"
    "errors"
    "fmt"
    "time"
    "strings"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"

    "auth-service/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uuid.UUID) (models.User, error)
    FindByEmail(email string) (*models.User, error)
    FindByUsername(username string) (*models.User, error)
    Update(user *models.User) error
    CreateProfile(profile *models.UserProfile) error
    UpdateProfile(profile *models.UserProfile) error
    AssignRole(userID, roleID uuid.UUID) error
    RevokeRole(userID, roleID uuid.UUID) error
}

type RoleRepository interface {
    FindByName(name string) (*models.Role, error)
}

type PermissionRepository interface{}

type Cache interface {
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Get(ctx context.Context, key string, dest interface{}) error
    Delete(ctx context.Context, key string) error
}

type Event struct {
    Type      string                 `json:"type"`
    Data      map[string]interface{} `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
}

type KafkaProducer interface {
    Publish(ctx context.Context, topic string, event Event) error
}

type AuthService struct {
    userRepo       UserRepository
    roleRepo       RoleRepository
    permissionRepo PermissionRepository
    cache          Cache
    kafkaProducer  KafkaProducer
    config         *AuthConfig
}

type JWTClaims struct {
    UserID      string   `json:"user_id"`
    Email       string   `json:"email"`
    Username    string   `json:"username"`
    Roles       []string `json:"roles"`
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}

type AuthConfig struct {
    JWTSecret        []byte
    JWTAccessExpiry  time.Duration
    JWTRefreshExpiry time.Duration
    KafkaTopics      struct{
        UserRegistered string
        UserUpdated    string
        UserDeleted    string
        RoleAssigned   string
        RoleRevoked    string
    }
}

func NewAuthService(userRepo UserRepository, roleRepo RoleRepository, permissionRepo PermissionRepository, cache Cache, kafkaProducer KafkaProducer, cfg *AuthConfig) *AuthService {
    return &AuthService{
        userRepo:       userRepo,
        roleRepo:       roleRepo,
        permissionRepo: permissionRepo,
        cache:          cache,
        kafkaProducer:  kafkaProducer,
        config:         cfg,
    }
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.TokenPair, error) {
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))
    req.Username = strings.TrimSpace(req.Username)
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	existingUser, _ = s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("user already exists with this username")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
    user := &models.User{
        ID:        uuid.New(),
        Email:     req.Email,
        Username:  req.Username,
        Password:  string(hashedPassword),
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Phone:     req.Phone,
        IsActive:  true,
    }

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default role (customer)
    customerRole, err := s.roleRepo.FindByName(models.RoleCustomer)
    if err == nil && customerRole != nil {
        if err := s.userRepo.AssignRole(user.ID, customerRole.ID); err != nil {
            return nil, fmt.Errorf("failed to assign role: %w", err)
        }
    }

	// Create user profile
    profile := &models.UserProfile{ ID: uuid.New(), UserID: user.ID }
	if err := s.userRepo.CreateProfile(profile); err != nil {
		return nil, fmt.Errorf("failed to create user profile: %w", err)
	}

    // Generate tokens
tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Cache user data
    if err := s.cacheUserData(user); err != nil {
        // Log error but don't fail the registration
        fmt.Printf("Failed to cache user data: %v\n", err)
    }

	// Publish user registered event
    event := Event{
        Type: "user.registered",
        Data: map[string]interface{}{
            "user_id":    user.ID,
            "email":      user.Email,
            "username":   user.Username,
            "first_name": user.FirstName,
            "last_name":  user.LastName,
            "created_at": user.CreatedAt,
        },
        Timestamp: time.Now(),
    }

    if err := s.kafkaProducer.Publish(ctx, s.config.KafkaTopics.UserRegistered, event); err != nil {
        // Log error but don't fail the registration
        fmt.Printf("Failed to publish user registered event: %v\n", err)
    }

	return tokenPair, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.TokenPair, *models.UserResponse, error) {
    // Normalize email and find user
    email := strings.ToLower(strings.TrimSpace(req.Email))
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return nil, nil, errors.New("invalid credentials")
    }

	// Check if user is active
    if !user.IsActive {
        return nil, nil, errors.New("account is deactivated")
    }

	// Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, nil, errors.New("invalid credentials")
    }

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(user); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	// Generate tokens
    tokenPair, err := s.generateTokenPair(user)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
    }

	// Cache user data
    if err := s.cacheUserData(user); err != nil {
        // Log error but don't fail login
        fmt.Printf("Failed to cache user data: %v\n", err)
    }

    return tokenPair, s.userToResponse(user), nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenPair, error) {
	// Validate refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if token is refresh token type
	if claims.Subject != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Get user from database
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

    // Generate new token pair
    tokenPair, err := s.generateTokenPair(&user)
    if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	// Remove user data from cache
    cacheKey := fmt.Sprintf("user:%s", userID)
    if s.cache != nil { if err := s.cache.Delete(ctx, cacheKey); err != nil {
        // Log error but don't fail logout
        fmt.Printf("Failed to delete user from cache: %v\n", err)
    } }

	return nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.UserResponse, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%s", userID)
	var user models.User
    if s.cache != nil { if err := s.cache.Get(ctx, cacheKey, &user); err == nil {
        // Found in cache, return response
        return s.userToResponse(&user), nil
    } }

	// Get from database
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err = s.userRepo.FindByID(userUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Cache for future use
    if err := s.cacheUserData(&user); err != nil {
        fmt.Printf("Failed to cache user data: %v\n", err)
    }

	return s.userToResponse(&user), nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(&user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Update profile if exists
	if user.Profile != nil {
		if req.Bio != "" {
			user.Profile.Bio = req.Bio
		}
		if req.DateOfBirth != "" {
			if dob, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
				user.Profile.DateOfBirth = &dob
			}
		}
		if req.Gender != "" {
			user.Profile.Gender = req.Gender
		}
		if req.Country != "" {
			user.Profile.Country = req.Country
		}
		if req.City != "" {
			user.Profile.City = req.City
		}
		if req.Timezone != "" {
			user.Profile.Timezone = req.Timezone
		}
		if req.Language != "" {
			user.Profile.Language = req.Language
		}

		if err := s.userRepo.UpdateProfile(user.Profile); err != nil {
			return nil, fmt.Errorf("failed to update profile: %w", err)
		}
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", userID)
    if s.cache != nil { if err := s.cache.Delete(ctx, cacheKey); err != nil {
        fmt.Printf("Failed to delete user from cache: %v\n", err)
    } }

	// Publish user updated event
    event := Event{
        Type: "user.updated",
        Data: map[string]interface{}{
            "user_id":     user.ID,
            "email":       user.Email,
            "username":    user.Username,
            "first_name":  user.FirstName,
            "last_name":   user.LastName,
            "phone":       user.Phone,
            "updated_at":  user.UpdatedAt,
        },
        Timestamp: time.Now(),
    }

    if err := s.kafkaProducer.Publish(ctx, s.config.KafkaTopics.UserUpdated, event); err != nil {
        fmt.Printf("Failed to publish user updated event: %v\n", err)
    }

	return s.userToResponse(&user), nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(&user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) AssignRole(ctx context.Context, userID string, roleName string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return errors.New("user not found")
	}

	role, err := s.roleRepo.FindByName(roleName)
	if err != nil {
		return errors.New("role not found")
	}

	if err := s.userRepo.AssignRole(user.ID, role.ID); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", userID)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to delete user from cache: %v\n", err)
	}

	// Publish role assigned event
    event := Event{
        Type: "role.assigned",
        Data: map[string]interface{}{
            "user_id":    user.ID,
            "role_name":  role.Name,
            "assigned_at": time.Now(),
        },
        Timestamp: time.Now(),
    }

    if err := s.kafkaProducer.Publish(ctx, s.config.KafkaTopics.RoleAssigned, event); err != nil {
        fmt.Printf("Failed to publish role assigned event: %v\n", err)
    }

	return nil
}

func (s *AuthService) RevokeRole(ctx context.Context, userID string, roleName string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return errors.New("user not found")
	}

	role, err := s.roleRepo.FindByName(roleName)
	if err != nil {
		return errors.New("role not found")
	}

	if err := s.userRepo.RevokeRole(user.ID, role.ID); err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", userID)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to delete user from cache: %v\n", err)
	}

	// Publish role revoked event
    event := Event{
        Type: "role.revoked",
        Data: map[string]interface{}{
            "user_id":     user.ID,
            "role_name":   role.Name,
            "revoked_at":  time.Now(),
        },
        Timestamp: time.Now(),
    }

    if err := s.kafkaProducer.Publish(ctx, s.config.KafkaTopics.RoleRevoked, event); err != nil {
        fmt.Printf("Failed to publish role revoked event: %v\n", err)
    }

	return nil
}

// Helper methods

func (s *AuthService) generateTokenPair(user *models.User) (*models.TokenPair, error) {
	// Get user roles and permissions
	roles, permissions, err := s.getUserRolesAndPermissions(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles and permissions: %w", err)
	}

	// Create access token
	accessToken, err := s.createAccessToken(user, roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Create refresh token
	refreshToken, err := s.createRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWTAccessExpiry).Unix(),
	}, nil
}

func (s *AuthService) createAccessToken(user *models.User, roles, permissions []string) (string, error) {
    claims := &JWTClaims{
        UserID:      user.ID.String(),
        Email:       user.Email,
        Username:    user.Username,
        Roles:       roles,
        Permissions: permissions,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   user.ID.String(),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWTAccessExpiry)),
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.config.JWTSecret)
}

func (s *AuthService) createRefreshToken(user *models.User) (string, error) {
    claims := &JWTClaims{
        UserID: user.ID.String(),
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   "refresh",
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWTRefreshExpiry)),
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.config.JWTSecret)
}

func (s *AuthService) validateToken(tokenString string) (*JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.config.JWTSecret, nil
    })

	if err != nil {
		return nil, err
	}

    claims, ok := token.Claims.(*JWTClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token claims")
    }

    return claims, nil
}

func (s *AuthService) getUserRolesAndPermissions(userID uuid.UUID) ([]string, []string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}

	var roles []string
	var permissions []string

	// Get roles
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
		// Get permissions from roles
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Name)
		}
	}

	// Get direct permissions
	for _, perm := range user.Permissions {
		permissions = append(permissions, perm.Name)
	}

	// Remove duplicates from permissions
	permissionMap := make(map[string]bool)
	for _, perm := range permissions {
		permissionMap[perm] = true
	}
	permissions = []string{}
	for perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return roles, permissions, nil
}

func (s *AuthService) cacheUserData(user *models.User) error {
    if s.cache == nil { return nil }
    ctx := context.Background()
    cacheKey := fmt.Sprintf("user:%s", user.ID.String())
    return s.cache.Set(ctx, cacheKey, user, 1*time.Hour)
}

func (s *AuthService) userToResponse(user *models.User) *models.UserResponse {
	var roles []models.RoleResponse
	for _, role := range user.Roles {
		roles = append(roles, models.RoleResponse{
			ID:          role.ID.String(),
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return &models.UserResponse{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Phone:           user.Phone,
		AvatarURL:       user.AvatarURL,
		IsEmailVerified: user.IsEmailVerified,
		IsPhoneVerified: user.IsPhoneVerified,
		IsActive:        user.IsActive,
		LastLoginAt:     user.LastLoginAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Roles:           roles,
		Profile:         user.Profile,
	}
}