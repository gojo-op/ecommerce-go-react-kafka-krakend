package services

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "golang.org/x/crypto/bcrypt"

    models "auth-service/internal/models"
)

// Mock implementations
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateProfile(profile *models.UserProfile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateProfile(profile *models.UserProfile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockUserRepository) AssignRole(userID, roleID uuid.UUID) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

type MockRoleRepository struct {
    mock.Mock
}

func (m *MockRoleRepository) FindByName(name string) (*models.Role, error) {
    args := m.Called(name)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Role), args.Error(1)
}

type MockCache struct {
    mock.Mock
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

type MockKafkaProducer struct {
    mock.Mock
}

func (m *MockKafkaProducer) Publish(ctx context.Context, topic string, event Event) error {
    args := m.Called(ctx, topic, event)
    return args.Error(0)
}

func (m *MockKafkaProducer) Close() error {
    args := m.Called()
    return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockCache := new(MockCache)
	mockKafka := new(MockKafkaProducer)
	
    cfg := &AuthConfig{
        JWTSecret:        []byte("test-secret"),
        JWTAccessExpiry:  15 * time.Minute,
        JWTRefreshExpiry: 7 * 24 * time.Hour,
    }

    service := &AuthService{
        userRepo:      mockUserRepo,
        roleRepo:      mockRoleRepo,
        cache:         mockCache,
        kafkaProducer: mockKafka,
        config:        cfg,
    }

	ctx := context.Background()
	req := &models.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	t.Run("Successful registration", func(t *testing.T) {
		// Mock expectations
		mockUserRepo.On("FindByEmail", req.Email).Return(nil, errors.New("not found"))
		mockUserRepo.On("FindByUsername", req.Username).Return(nil, errors.New("not found"))
		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
		
		customerRole := &models.Role{
			ID:   uuid.New(),
			Name: models.RoleCustomer,
		}
		mockRoleRepo.On("FindByName", models.RoleCustomer).Return(customerRole, nil)
		
		mockUserRepo.On("AssignRole", mock.AnythingOfType("uuid.UUID"), customerRole.ID).Return(nil)
		mockUserRepo.On("CreateProfile", mock.AnythingOfType("*models.UserProfile")).Return(nil)
		mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*models.User"), 1*time.Hour).Return(nil)
		mockKafka.On("Publish", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("kafka.Event")).Return(nil)

		// Execute
		tokenPair, err := service.Register(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, tokenPair)
		assert.NotEmpty(t, tokenPair.AccessToken)
		assert.NotEmpty(t, tokenPair.RefreshToken)

		// Verify expectations
		mockUserRepo.AssertExpectations(t)
		mockRoleRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
		mockKafka.AssertExpectations(t)
	})

	t.Run("User already exists with email", func(t *testing.T) {
		existingUser := &models.User{ID: uuid.New(), Email: req.Email}
		mockUserRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

		// Execute
		tokenPair, err := service.Register(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, tokenPair)
		assert.Equal(t, "user already exists with this email", err.Error())

		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockCache := new(MockCache)
	mockKafka := new(MockKafkaProducer)
	
    cfg := &AuthConfig{
        JWTSecret:        []byte("test-secret"),
        JWTAccessExpiry:  15 * time.Minute,
        JWTRefreshExpiry: 7 * 24 * time.Hour,
    }

	service := &AuthService{
		userRepo:      mockUserRepo,
		roleRepo:      mockRoleRepo,
		cache:         mockCache,
		kafkaProducer: mockKafka,
		config:        cfg,
	}

	ctx := context.Background()
	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

    t.Run("Successful login", func(t *testing.T) {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
        user := &models.User{
            ID:       uuid.New(),
            Email:    req.Email,
            Password: string(hashed),
            IsActive: true,
        }

        mockUserRepo.On("FindByEmail", req.Email).Return(user, nil)
        mockUserRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
        mockUserRepo.On("FindByID", user.ID).Return(*user, nil)
        mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*models.User"), 1*time.Hour).Return(nil)

        // Execute
        tokenPair, userResp, err := service.Login(ctx, req)

        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, tokenPair)
        assert.NotNil(t, userResp)
        assert.NotEmpty(t, tokenPair.AccessToken)
        assert.NotEmpty(t, tokenPair.RefreshToken)

        mockUserRepo.AssertExpectations(t)
    })

	t.Run("User not found", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", req.Email).Return(nil, errors.New("not found"))

		// Execute
        tokenPair, userResp, err := service.Login(ctx, req)

		// Assert
		assert.Error(t, err)
        assert.Nil(t, tokenPair)
        assert.Nil(t, userResp)
		assert.Equal(t, "invalid credentials", err.Error())

		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
    cfg := &AuthConfig{
        JWTSecret: []byte("test-secret"),
    }

	service := &AuthService{
		config: cfg,
	}

	t.Run("Valid token", func(t *testing.T) {
		// Create a valid token
		claims := &models.JWTClaims{
			UserID: "test-user-id",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "test-subject",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(cfg.JWTSecret)
		assert.NoError(t, err)

		// Validate
		validatedClaims, err := service.validateToken(tokenString)
		assert.NoError(t, err)
		assert.NotNil(t, validatedClaims)
		assert.Equal(t, claims.UserID, validatedClaims.UserID)
	})

	t.Run("Invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		
		validatedClaims, err := service.validateToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
	})
}