package repository

import (
	"github.com/abushaista/lms-backend/internal/domain"
	"gorm.io/gorm"
)

type GormCategoryRepository struct {
	db *gorm.DB
}

// Delete implements domain.CategoryRepository.
func (g *GormCategoryRepository) Delete(id uint) error {
	return g.db.Delete(&domain.Category{}, id).Error
}

// Create implements domain.CategoryRepository.
func (g *GormCategoryRepository) Save(category *domain.Category) error {
	if category.ID != 0 {
		return g.db.Save(category).Error
	}
	return g.db.Create(category).Error
}

// GetAll implements domain.CategoryRepository.
func (g *GormCategoryRepository) GetAll() ([]*domain.Category, error) {
	var categories []*domain.Category
	if err := g.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetByFilterAll implements domain.CategoryRepository.
func (g *GormCategoryRepository) GetByFilterAll(page int, limit int, filter string) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64
	query := g.db.Model(&domain.Category{})
	if filter != "" {
		query.Where("name LIKE ?", "%"+filter+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

// GetByID implements domain.CategoryRepository.
func (g *GormCategoryRepository) GetByID(id uint) (*domain.Category, error) {
	var category domain.Category
	if err := g.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func NewGormCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &GormCategoryRepository{db: db}
}
