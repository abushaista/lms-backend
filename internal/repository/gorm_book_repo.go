package repository

import (
	"errors"

	"github.com/abushaista/lms-backend/internal/domain"
	"gorm.io/gorm"
)

type GormBookRepository struct {
	db *gorm.DB
}

// Create implements domain.BookRepository.
func (g *GormBookRepository) Save(b *domain.Book) (int64, error) {
	if b.ID == 0 {
		res := g.db.Create(&b)
		return b.ID, res.Error
	}
	return b.ID, g.db.Save(b).Error
}

// Delete implements domain.BookRepository.
func (g *GormBookRepository) Delete(id int64) error {
	return g.db.Delete(&domain.Book{}, id).Error
}

// GetAll implements domain.BookRepository.
func (g *GormBookRepository) GetAll(page, limit int, filter domain.BookFilter) ([]*domain.Book, int64, error) {
	var books []*domain.Book
	var total int64
	query := g.db.Model(&domain.Book{})

	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Author != "" {
		query = query.Where("author LIKE ?", "%"+filter.Author+"%")
	}
	if filter.Year != 0 {
		query = query.Where("year = ?", filter.Year)
	}

	if filter.Category != 0 {
		query.Where("categoryId = ?", filter.Category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

// GetByID implements domain.BookRepository.
func (g *GormBookRepository) GetByID(id int64) (*domain.Book, error) {
	var book domain.Book
	if err := g.db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func NewGormBookRepository(db *gorm.DB) domain.BookRepository {
	return &GormBookRepository{db: db}
}
