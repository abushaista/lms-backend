package http

import (
	"net/http"

	"github.com/abushaista/lms-backend/delivery/utils"
	"github.com/abushaista/lms-backend/internal/dto"
	"github.com/abushaista/lms-backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type CleanUpHandler struct {
	uc         *usecase.CleanUpUseCase
	validate   *validator.Validate
	rootLogger zerolog.Logger
}

func NewCleanUpHandler(e *echo.Group, uc *usecase.CleanUpUseCase, logger zerolog.Logger) {
	h := &CleanUpHandler{
		uc:         uc,
		validate:   validator.New(),
		rootLogger: logger,
	}
	e.POST("/clean", h.CleanUpUrl)
}

func (h *CleanUpHandler) CleanUpUrl(c echo.Context) error {
	logger := utils.WithRequestLogger(h.rootLogger, c)
	var req dto.CleanUpRequest
	if err := c.Bind(&req); err != nil {
		logger.Error().Err(err).Msg("invalid payload")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if err := c.Validate(&req); err != nil {
		logger.Error().Err(err).Msg("invalid payload")
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}

	processedUrl := req.Url

	if req.Operation == "canonical" || req.Operation == "all" {
		url, err := h.uc.CanonicalURL(processedUrl)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.FormatValidationErrors(err))
		}
		processedUrl = url
	}

	if req.Operation == "redirection" || req.Operation == "all" {
		url, err := h.uc.Redirection(processedUrl)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.FormatValidationErrors(err))
		}
		processedUrl = url
	}

	return c.JSON(http.StatusOK, echo.Map{"processed_url": processedUrl})
}
