package domain

import "time"

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryRepository interface {
	Create(category *Category) error
	GetAll() ([]Category, error)
	GetByFilterAll(page, limit int, filter string) ([]Category, int64, error)
	GetByID(id uint) (*Category, error)
}
