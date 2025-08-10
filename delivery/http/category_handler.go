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
	e.POST("/categories", h.CreateCategory)
	e.PUT("/categories/:id", h.UpdateCategory)
	e.DELETE("categories/:id", h.Delete)
	e.GET("/categories", h.GetByFilterAll)
	e.GET("/categories/:id", h.GetByID)
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Add a new category to the system
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CategoryRequest  true  "Category payload"
// @Success      201   {object}  domain.Category
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]string
// @Router       /categories [post]
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

// UpdateCategory godoc
// @Summary      Update a category by ID
// @Description  Update category details by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "Category ID"
// @Param        body  body      dto.CategoryRequest  true  "Category payload"
// @Success      200   {object}  domain.Category
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]string
// @Router       /categories/{id} [put]
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

// DeleteCategory godoc
// @Summary      Delete a category by ID
// @Description  Remove a category from the system by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories/{id} [delete]
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

// GetByFilterAll godoc
// @Summary      Get categories by filter with pagination
// @Description  Retrieve a list of categories filtered by a search string with pagination
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        page    query     int     false  "Page number"     default(1)
// @Param        limit   query     int     false  "Items per page"  default(10)
// @Param        filter  query     string  false  "Filter string"
// @Success      200     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]string
// @Router       /categories [get]
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

// GetByID godoc
// @Summary      Get a category by ID
// @Description  Retrieve a single category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  domain.Category
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /categories/{id} [get]
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
