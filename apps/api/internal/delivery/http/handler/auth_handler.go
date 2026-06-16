package handler

import (
	"announcement-api/internal/delivery/http/dto"
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC *usecase.AuthUsecase
}

func NewAuthHandler(authUC *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type RegisterRequest struct {
	Username string      `json:"username" binding:"required"`
	Password string      `json:"password" binding:"required,min=6"`
	Role     entity.Role `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary     Daftarkan akun baru
// @Description Create new user
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       body  body      dto.RegisterRequest   true  "Registration data"
// @Success     201   {object}  dto.RegisterResponse
// @Failure     400   {object}  dto.ErrorResponse     "Validation failed"
// @Failure     409   {object}  dto.ErrorResponse     "Username has been taken"
// @Failure     500   {object}  dto.ErrorResponse     "Internal server error"
// @Router      /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.authUC.Register(c.Request.Context(), req.Username, req.Password, entity.Role(req.Role))
	if err != nil {
		if err == usecase.ErrUsernameTaken {
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to register"})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     string(user.Role),
	})
}

// @Summary     Login user
// @Description Authenticate user and create JWT token.
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       body  body      dto.LoginRequest  true  "Login credentials"
// @Success     200   {object}  dto.AuthResponse
// @Failure     400   {object}  dto.ErrorResponse  "Validation failed"
// @Failure     401   {object}  dto.ErrorResponse  "Credential not valid"
// @Failure     500   {object}  dto.ErrorResponse  "Failed to create token"
// @Router      /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	token, user, err := h.authUC.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Message: "Login success",
		Token:   token,
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Role:     string(user.Role),
		},
	})
}
