package domain

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID            int64          `gorm:"primaryKey" json:"id"`
	Title         string         `json:"title"`
	Author        string         `json:"author"`
	ISBN          string         `json:"isbn"`
	Year          int            `json:"year"`
	CategoryID    uint           `json:"category_id"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Summary       string         `json:"summary"`
	Available     bool           `gorm:"default:true;not null" json:"available"`
	CoverImageURL string         `json:"cover_image_url"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type BookRepository interface {
	Save(b *Book) (int64, error)
	GetAll(page, limit int, filter BookFilter) ([]*Book, int64, error)
	GetByID(id int64) (*Book, error)
	Delete(id int64) error
}
