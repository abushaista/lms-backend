package usecase

import (
	"github.com/abushaista/lms-backend/internal/domain"
	"github.com/abushaista/lms-backend/internal/dto"
	"github.com/go-playground/validator/v10"
)

type CategoryUseCase struct {
	repo      domain.CategoryRepository
	validator *validator.Validate
}

func NewCategoryUseCase(r domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		repo:      r,
		validator: validator.New(),
	}
}

func (uc *CategoryUseCase) Save(req dto.CategoryRequest) (*domain.Category, error) {
	if err := uc.validator.Struct(req); err != nil {
		return nil, err
	}

	category := domain.Category{
		ID:   req.ID,
		Name: req.Name,
	}
	if err := uc.repo.Save(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (uc *CategoryUseCase) GetByID(id uint) (*domain.Category, error) {
	return uc.repo.GetByID(id)
}

func (uc *CategoryUseCase) GetByFilterAll(page, limit int, filter string) ([]*domain.Category, int64, error) {
	data, total, err := uc.repo.GetByFilterAll(page, limit, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (uc *CategoryUseCase) Delete(id uint) error {
	return uc.repo.Delete(id)
}

func (uc *CategoryUseCase) GetAll() ([]*domain.Category, error) {
	return uc.repo.GetAll()
}
