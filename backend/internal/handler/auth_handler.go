package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/service"
)

// ExtendedAuthHandler handles all authentication-related HTTP requests.
type ExtendedAuthHandler struct {
	authService service.ExtendedAuthService
}

// NewExtendedAuthHandler creates a new ExtendedAuthHandler instance.
func NewExtendedAuthHandler(authService service.ExtendedAuthService) *ExtendedAuthHandler {
	return &ExtendedAuthHandler{authService: authService}
}

// LogoutRequest represents a logout request.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UserResponse represents the current user response.
type UserResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	TwoFAEnabled bool   `json:"two_fa_enabled"`
}

// OAuthRequest represents an OAuth callback request.
type OAuthRequest struct {
	Provider       string `json:"provider" binding:"required"`
	ProviderUserID string `json:"provider_user_id" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	AccessToken    string `json:"access_token" binding:"required"`
	RefreshToken   string `json:"refresh_token"`
}

// OAuthResponse represents an OAuth login response.
type OAuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

// TwoFASetupResponse represents a 2FA setup response.
type TwoFASetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

// TwoFAVerifyRequest represents a 2FA verification request.
type TwoFAVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// LoginWith2FARequest represents a login request with 2FA code.
type LoginWith2FARequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// Register handles user registration.
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *ExtendedAuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.Name)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to register user"})
		return
	}

	// Log audit event with request context
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	_ = h.authService.LogAuditEvent(&user.ID, model.AuditActionRegister, ipAddress, userAgent, "", true)

	c.JSON(http.StatusCreated, RegisterResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	})
}

// Login handles user login.
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 428 {object} ErrorResponse "2FA required"
// @Router /api/v1/auth/login [post]
func (h *ExtendedAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == service.Err2FARequired {
			c.JSON(http.StatusPreconditionRequired, gin.H{
				"error":        "2FA verification required",
				"requires_2fa": true,
			})
			return
		}
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to login"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// LoginWith2FA handles user login with 2FA code.
// @Summary Login with 2FA
// @Description Authenticate user with 2FA code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginWith2FARequest true "Login credentials with 2FA code"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login/2fa [post]
func (h *ExtendedAuthHandler) LoginWith2FA(c *gin.Context) {
	var req LoginWith2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.ValidateLoginWith2FA(req.Email, req.Password, req.Code)
	if err != nil {
		if err == service.ErrInvalidCredentials || err == service.Err2FAInvalidCode {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
			return
		}
		if err == service.Err2FANotEnabled {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to login"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Refresh handles token refresh.
// @Summary Refresh access token
// @Description Generate a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *ExtendedAuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	accessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken: accessToken,
	})
}

// Logout handles user logout.
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LogoutRequest true "Refresh token to invalidate"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *ExtendedAuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// GetCurrentUser returns the current authenticated user.
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/me [get]
func (h *ExtendedAuthHandler) GetCurrentUser(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id"})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		Name:         user.Name,
		Role:         user.Role,
		TwoFAEnabled: user.TwoFAEnabled,
	})
}

// GoogleOAuth handles Google OAuth callback.
// @Summary Google OAuth login
// @Description Authenticate or register user via Google OAuth
// @Tags auth
// @Accept json
// @Produce json
// @Param request body OAuthRequest true "Google OAuth data"
// @Success 200 {object} OAuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/oauth/google [post]
func (h *ExtendedAuthHandler) GoogleOAuth(c *gin.Context) {
	h.handleOAuth(c, model.OAuthProviderGoogle)
}

// GitHubOAuth handles GitHub OAuth callback.
// @Summary GitHub OAuth login
// @Description Authenticate or register user via GitHub OAuth
// @Tags auth
// @Accept json
// @Produce json
// @Param request body OAuthRequest true "GitHub OAuth data"
// @Success 200 {object} OAuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/auth/oauth/github [post]
func (h *ExtendedAuthHandler) GitHubOAuth(c *gin.Context) {
	h.handleOAuth(c, model.OAuthProviderGitHub)
}

func (h *ExtendedAuthHandler) handleOAuth(c *gin.Context, expectedProvider model.OAuthProvider) {
	var req OAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Validate provider matches endpoint
	if model.OAuthProvider(req.Provider) != expectedProvider {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid provider for this endpoint"})
		return
	}

	info := &service.OAuthUserInfo{
		Provider:       expectedProvider,
		ProviderUserID: req.ProviderUserID,
		Email:          req.Email,
		Name:           req.Name,
		AvatarURL:      req.AvatarURL,
		AccessToken:    req.AccessToken,
		RefreshToken:   req.RefreshToken,
	}

	user, accessToken, refreshToken, err := h.authService.HandleOAuthLogin(info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to authenticate with " + string(expectedProvider)})
		return
	}

	c.JSON(http.StatusOK, OAuthResponse{
		User: UserResponse{
			ID:           user.ID.String(),
			Email:        user.Email,
			Name:         user.Name,
			Role:         user.Role,
			TwoFAEnabled: user.TwoFAEnabled,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Setup2FA sets up 2FA for the current user.
// @Summary Setup 2FA
// @Description Generate TOTP secret and backup codes for 2FA
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} TwoFASetupResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/2fa/setup [post]
func (h *ExtendedAuthHandler) Setup2FA(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	setup, err := h.authService.Setup2FA(userID)
	if err != nil {
		if err == service.Err2FAAlreadyEnabled {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to setup 2FA"})
		return
	}

	c.JSON(http.StatusOK, TwoFASetupResponse{
		Secret:      setup.Secret,
		QRCodeURL:   setup.QRCodeURL,
		BackupCodes: setup.BackupCodes,
	})
}

// Verify2FA verifies and enables 2FA for the current user.
// @Summary Verify 2FA
// @Description Verify TOTP code to enable 2FA
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TwoFAVerifyRequest true "TOTP code"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/2fa/verify [post]
func (h *ExtendedAuthHandler) Verify2FA(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req TwoFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.authService.Verify2FA(userID, req.Code); err != nil {
		if err == service.Err2FAInvalidCode {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if err == service.Err2FANotEnabled {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "2FA setup not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to verify 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FA disables 2FA for the current user.
// @Summary Disable 2FA
// @Description Disable 2FA for the current user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TwoFAVerifyRequest true "TOTP code for verification"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/2fa/disable [post]
func (h *ExtendedAuthHandler) Disable2FA(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req TwoFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.authService.Disable2FA(userID, req.Code); err != nil {
		if err == service.Err2FAInvalidCode {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if err == service.Err2FANotEnabled {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to disable 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// Helper to extract user ID from context
func (h *ExtendedAuthHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, service.ErrInvalidToken
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		return uuid.Nil, service.ErrInvalidToken
	}

	return uuid.Parse(userIDStr)
}

// RegisterExtendedAuthRoutes registers all authentication routes.
func (h *ExtendedAuthHandler) RegisterExtendedAuthRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	auth := rg.Group("/auth")
	{
		// Public routes
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/login/2fa", h.LoginWith2FA)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/oauth/google", h.GoogleOAuth)
		auth.POST("/oauth/github", h.GitHubOAuth)

		// Protected routes
		protected := auth.Group("")
		protected.Use(authMiddleware)
		{
			protected.POST("/logout", h.Logout)
			protected.GET("/me", h.GetCurrentUser)
			protected.POST("/2fa/setup", h.Setup2FA)
			protected.POST("/2fa/verify", h.Verify2FA)
			protected.POST("/2fa/disable", h.Disable2FA)
		}
	}
}
