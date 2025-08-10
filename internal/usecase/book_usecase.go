package usecase

import (
	"github.com/abushaista/lms-backend/internal/domain"
	"github.com/abushaista/lms-backend/internal/dto"
	"github.com/go-playground/validator/v10"
)

type BookUseCase struct {
	repo      domain.BookRepository
	validator *validator.Validate
}

func NewBookUsecase(r domain.BookRepository) *BookUseCase {
	return &BookUseCase{
		repo:      r,
		validator: validator.New(),
	}
}

func (uc *BookUseCase) CreateBook(req dto.CreateBookRequest) (*domain.Book, error) {
	if err := uc.validator.Struct(req); err != nil {
		return nil, err
	}
	book := domain.Book{
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		Year:          req.Year,
		Summary:       req.Summary,
		CoverImageURL: req.CoverImage,
		CategoryID:    req.CategoryID,
		Available:     req.Available,
	}
	id, err := uc.repo.Save(&book)
	// Save to DB
	if err != nil {
		return nil, err
	}
	book.ID = id
	return &book, nil
}

func (uc *BookUseCase) UpdateBook(req dto.UpdateBookRequest) (*domain.Book, error) {
	if err := uc.validator.Struct(req); err != nil {
		return nil, err
	}
	book := domain.Book{
		ID:            req.ID,
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		Year:          req.Year,
		Summary:       req.Summary,
		CoverImageURL: req.CoverImage,
		CategoryID:    req.CategoryID,
	}
	_, err := uc.repo.Save(&book)

	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (uc *BookUseCase) GetByFilterAll(page, limit int, filter domain.BookFilter) ([]*domain.Book, int64, error) {
	data, total, err := uc.repo.GetAll(page, limit, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (uc *BookUseCase) DeleteBook(id int64) error {
	return uc.repo.Delete(id)
}

func (uc *BookUseCase) GetByID(id int64) (*domain.Book, error) {
	return uc.repo.GetByID(id)
}
