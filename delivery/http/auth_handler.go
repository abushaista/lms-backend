package http

import (
	"net/http"

	"github.com/abushaista/lms-backend/delivery/utils"
	"github.com/abushaista/lms-backend/internal/dto"
	"github.com/abushaista/lms-backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	uc       *usecase.AuthUseCase
	validate *validator.Validate
}

func NewAuthHandler(e *echo.Group, uc *usecase.AuthUseCase) {
	h := &AuthHandler{uc: uc,
		validate: validator.New()}
	e.POST("/register", h.Create)
	e.POST("/login", h.Login)
}

// Register godoc
// @Summary Register a user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User registration"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/register [post]
func (h AuthHandler) Create(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}
	_, err := h.uc.Create(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"message": "user registered"})
}

// Login godoc
// @Summary Login user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "User login"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/login [post]
func (h AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}
	token, err := h.uc.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
