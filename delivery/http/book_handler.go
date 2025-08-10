package http

import (
	"net/http"
	"strconv"

	"github.com/abushaista/lms-backend/delivery/utils"
	"github.com/abushaista/lms-backend/internal/domain"
	"github.com/abushaista/lms-backend/internal/dto"
	"github.com/abushaista/lms-backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	uc       *usecase.BookUseCase
	validate *validator.Validate
}

func NewBookHandler(e *echo.Group, uc *usecase.BookUseCase) {
	h := &BookHandler{
		uc:       uc,
		validate: validator.New(),
	}
	e.POST("/book", h.CreateBook)
	e.GET("/book/:id", h.GetByID)
	e.GET("/books", h.GetByFilterAll)
	e.DELETE("book/:id", h.Delete)
	e.PUT("/book/:id", h.UpdateBook)
}

func (h *BookHandler) CreateBook(c echo.Context) error {
	var req dto.CreateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}
	book, err := h.uc.CreateBook(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) GetByFilterAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	filter := domain.BookFilter{
		Title:   c.QueryParam("title"),
		Author:  c.QueryParam("author"),
		Summary: c.QueryParam("summary"),
	}

	category, err := strconv.Atoi(c.QueryParam("category"))
	if err == nil {
		filter.Category = category
	}

	year, err := strconv.Atoi(c.QueryParam("year"))
	if err == nil {
		filter.Year = year
	}

	books, total, err := h.uc.GetByFilterAll(page, limit, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  books,
		"total": total,
		"page":  page,
	})
}

func (h *BookHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	data, err := h.uc.GetByID(int64(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

func (h *BookHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.uc.DeleteBook(int64(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "book deleted"})
}

func (h *BookHandler) UpdateBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req dto.UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	req.ID = int64(id)
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}
	book, err := h.uc.UpdateBook(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, book)
}
