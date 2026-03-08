package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    models "auth-service/internal/models"
    utils "auth-service/internal/utils"
    "auth-service/internal/services"
    "auth-service/internal/repositories"
)

type AuthController struct {
    authService *services.AuthService
    addrRepo    *repositories.AddressRepository
}

func NewAuthController(authService *services.AuthService, addrRepo *repositories.AddressRepository) *AuthController {
    return &AuthController{ authService: authService, addrRepo: addrRepo }
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} utils.Response{data=models.TokenPair}
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	tokenPair, err := c.authService.Register(ctx, &req)
	if err != nil {
		if err.Error() == "user already exists with this email" || err.Error() == "user already exists with this username" {
			utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User registered successfully", tokenPair)
}

// Login godoc
// @Summary Login user
// @Description Login user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} utils.Response{data=models.TokenPair}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
    var req models.LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
        return
    }

    tokenPair, userResp, err := c.authService.Login(ctx, &req)
    if err != nil {
        if err.Error() == "invalid credentials" {
            utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid credentials", nil)
            return
        }
		if err.Error() == "account is deactivated" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Account is deactivated", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to login", err.Error())
		return
	}

    utils.SuccessResponse(ctx, http.StatusOK, "Login successful", map[string]interface{}{
        "access_token":  tokenPair.AccessToken,
        "refresh_token": tokenPair.RefreshToken,
        "expires_at":    tokenPair.ExpiresAt,
        "user":          userResp,
    })
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token request" Example("{"refresh_token": "your-refresh-token"}")
// @Success 200 {object} utils.Response{data=models.TokenPair}
// @Failure 401 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	tokenPair, err := c.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err.Error() == "invalid refresh token" || err.Error() == "invalid token type" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid refresh token", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to refresh token", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Token refreshed successfully", tokenPair)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and invalidate tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/logout [post]
// @Security BearerAuth
func (c *AuthController) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if err := c.authService.Logout(ctx, userID.(string)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to logout", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Logout successful", nil)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 404 {object} utils.Response
// @Router /auth/profile [get]
// @Security BearerAuth
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	profile, err := c.authService.GetProfile(ctx, userID.(string))
	if err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get profile", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", profile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.UpdateProfileRequest true "Profile update request"
// @Success 200 {object} utils.Response{data=models.UserResponse}
// @Failure 400 {object} utils.Response
// @Router /auth/profile [put]
// @Security BearerAuth
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	profile, err := c.authService.UpdateProfile(ctx, userID.(string), &req)
	if err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profile updated successfully", profile)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change user password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.ChangePasswordRequest true "Password change request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/change-password [post]
// @Security BearerAuth
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := c.authService.ChangePassword(ctx, userID.(string), &req); err != nil {
		if err.Error() == "incorrect old password" {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Incorrect old password", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to change password", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Password changed successfully", nil)
}

// Address endpoints
func (c *AuthController) ListAddresses(ctx *gin.Context) {
    userIDVal, ok := ctx.Get("user_id"); if !ok { utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil); return }
    uid, err := uuid.Parse(userIDVal.(string)); if err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error()); return }
    items, err := c.addrRepo.ListByUser(uid); if err != nil { utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to list addresses", err.Error()); return }
    utils.SuccessResponse(ctx, http.StatusOK, "Addresses", items)
}

func (c *AuthController) CreateAddress(ctx *gin.Context) {
    userIDVal, ok := ctx.Get("user_id"); if !ok { utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil); return }
    uid, err := uuid.Parse(userIDVal.(string)); if err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error()); return }
    var req models.Address
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error()); return }
    req.UserID = uid
    if err := c.addrRepo.Create(&req); err != nil { utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create address", err.Error()); return }
    utils.SuccessResponse(ctx, http.StatusCreated, "Address created", req)
}

func (c *AuthController) UpdateAddress(ctx *gin.Context) {
    userIDVal, ok := ctx.Get("user_id"); if !ok { utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil); return }
    uid, err := uuid.Parse(userIDVal.(string)); if err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error()); return }
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr); if err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid address ID", err.Error()); return }
    var req models.Address
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error()); return }
    req.ID = id; req.UserID = uid
    if err := c.addrRepo.Update(&req); err != nil { utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update address", err.Error()); return }
    utils.SuccessResponse(ctx, http.StatusOK, "Address updated", req)
}

func (c *AuthController) DeleteAddress(ctx *gin.Context) {
    userIDVal, ok := ctx.Get("user_id"); if !ok { utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil); return }
    uid, err := uuid.Parse(userIDVal.(string)); if err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error()); return }
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr); if err != nil { utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid address ID", err.Error()); return }
    if err := c.addrRepo.Delete(uid, id); err != nil { utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete address", err.Error()); return }
    utils.SuccessResponse(ctx, http.StatusOK, "Address deleted", nil)
}

// AssignRole godoc
// @Summary Assign role to user
// @Description Assign a role to a user (Admin only)
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body map[string]string true "Role assignment request" Example("{"role": "admin"}")
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/users/{user_id}/roles [post]
// @Security BearerAuth
func (c *AuthController) AssignRole(ctx *gin.Context) {
	userIDParam := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := c.authService.AssignRole(ctx, userID.String(), req.Role); err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", nil)
			return
		}
		if err.Error() == "role not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Role not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to assign role", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Role assigned successfully", nil)
}

// RevokeRole godoc
// @Summary Revoke role from user
// @Description Revoke a role from a user (Admin only)
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body map[string]string true "Role revocation request" Example("{"role": "admin"}")
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /auth/users/{user_id}/roles [delete]
// @Security BearerAuth
func (c *AuthController) RevokeRole(ctx *gin.Context) {
	userIDParam := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := c.authService.RevokeRole(ctx, userID.String(), req.Role); err != nil {
		if err.Error() == "user not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", nil)
			return
		}
		if err.Error() == "role not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Role not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to revoke role", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Role revoked successfully", nil)
}
