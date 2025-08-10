package http

import (
	"net/http"
	"strconv"

	"github.com/abushaista/lms-backend/delivery/utils"
	"github.com/abushaista/lms-backend/internal/dto"
	"github.com/abushaista/lms-backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	uc       *usecase.CategoryUseCase
	validate *validator.Validate
}

func NewCategoryHandler(e *echo.Group, uc *usecase.CategoryUseCase) {
	h := &CategoryHandler{
		uc:       uc,
		validate: validator.New(),
	}
	e.POST("/category", h.CreateCategory)
	e.PUT("/category/:id", h.UpdateCategory)
	e.DELETE("category/:id", h.Delete)
	e.GET("/categories", h.GetByFilterAll)
	e.GET("/category/:id", h.GetByID)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req dto.CategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}
	req.ID = 0
	category, err := h.uc.Save(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	var req dto.CategoryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	req.ID = uint(id)
	category, err := h.uc.Save(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.uc.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted"})
}

func (h *CategoryHandler) GetByFilterAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	filter := c.QueryParam("filter")
	categories, total, err := h.uc.GetByFilterAll(page, limit, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  categories,
		"total": total,
		"page":  page,
	})
}

func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	category, err := h.uc.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, category)
}
