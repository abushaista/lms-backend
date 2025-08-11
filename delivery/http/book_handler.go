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
	"github.com/rs/zerolog"
)

type BookHandler struct {
	uc         *usecase.BookUseCase
	validate   *validator.Validate
	rootLogger zerolog.Logger
}

func NewBookHandler(e *echo.Group, uc *usecase.BookUseCase, logger zerolog.Logger) {
	h := &BookHandler{
		uc:         uc,
		validate:   validator.New(),
		rootLogger: logger,
	}
	e.POST("/books", h.CreateBook)
	e.GET("/books/:id", h.GetByID)
	e.GET("/books", h.GetByFilterAll)
	e.DELETE("books/:id", h.Delete)
	e.PUT("/books/:id", h.UpdateBook)
}

// CreateBook godoc
// @Summary      Create a new book
// @Description  Add a new book to the library
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateBookRequest  true  "Create Book Payload"
// @Success      201   {object}  domain.Book
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]string
// @Router       /books [post]
func (h *BookHandler) CreateBook(c echo.Context) error {
	logger := utils.WithRequestLogger(h.rootLogger, c)
	var req dto.CreateBookRequest
	if err := c.Bind(&req); err != nil {
		logger.Warn().Err(err).Msg("bind books")
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

// GetByFilterAll godoc
// @Summary      Get books by filters with pagination
// @Description  Retrieve list of books filtered by title, author, summary, category, and year with pagination support
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"      default(1)
// @Param        limit     query     int     false  "Items per page"   default(10)
// @Param        title     query     string  false  "Filter by book title"
// @Param        author    query     string  false  "Filter by author name"
// @Param        summary   query     string  false  "Filter by book summary"
// @Param        category  query     int     false  "Filter by category ID"
// @Param        year      query     int     false  "Filter by publication year"
// @Success      200       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]string
// @Router       /books [get]
func (h *BookHandler) GetByFilterAll(c echo.Context) error {
	logger := utils.WithRequestLogger(h.rootLogger, c)
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
		logger.Error().Err(err).Msg("500 internal server error")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  books,
		"total": total,
		"page":  page,
	})
}

// GetByID godoc
// @Summary      Get a book by ID
// @Description  Retrieve a single book by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  domain.Book
// @Failure      500  {object}  map[string]string
// @Router       /books/{id} [get]
func (h *BookHandler) GetByID(c echo.Context) error {
	logger := utils.WithRequestLogger(h.rootLogger, c)
	id, _ := strconv.Atoi(c.Param("id"))

	data, err := h.uc.GetByID(int64(id))
	if err != nil {
		logger.Error().Err(err).Msg("500 internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}

// Delete godoc
// @Summary      Delete a book by ID
// @Description  Delete a book from the system by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /books/{id} [delete]
func (h *BookHandler) Delete(c echo.Context) error {
	logger := utils.WithRequestLogger(h.rootLogger, c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error().Err(err).Msg("invalid id")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.uc.DeleteBook(int64(id)); err != nil {
		logger.Error().Err(err).Msg("500 internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "book deleted"})
}

// UpdateBook godoc
// @Summary      Update a book by ID
// @Description  Update the details of an existing book by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int                 true  "Book ID"
// @Param        body  body      dto.UpdateBookRequest  true  "Update Book Payload"
// @Success      200   {object}  domain.Book
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]string
// @Router       /books/{id} [put]
func (h *BookHandler) UpdateBook(c echo.Context) error {
	logger := utils.WithRequestLogger(h.rootLogger, c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req dto.UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		logger.Warn().Err(err).Msg("invalid payload")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}
	req.ID = int64(id)
	if err := c.Validate(&req); err != nil {
		logger.Warn().Err(err).Msg("invalid payload")
		return c.JSON(http.StatusBadRequest, utils.FormatValidationErrors(err))
	}
	book, err := h.uc.UpdateBook(req)
	if err != nil {
		logger.Warn().Err(err).Msg("500 internal server error")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, book)
}
