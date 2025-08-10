package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Books     []Book         `json:"books"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CategoryRepository interface {
	Save(category *Category) error
	GetAll() ([]*Category, error)
	GetByFilterAll(page, limit int, filter string) ([]*Category, int64, error)
	GetByID(id uint) (*Category, error)
	Delete(id uint) error
}
