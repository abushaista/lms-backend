package repository

import (
	"errors"

	"github.com/abushaista/lms-backend/internal/domain"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

// CreateUser implements domain.UserRepository.
func (g *GormUserRepository) CreateUser(u *domain.User) (string, error) {
	err := g.db.Create(&u).Error
	if err != nil {
		return "", err
	}
	return u.ID.String(), nil
}

// GetByUsername implements domain.UserRepository.
func (g *GormUserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := g.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func NewGormUserRepository(db *gorm.DB) domain.UserRepository {
	return &GormUserRepository{db: db}
}
