package usecase

import (
	"errors"
	"time"

	"github.com/abushaista/lms-backend/internal/domain"
	"github.com/abushaista/lms-backend/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo      domain.UserRepository
	validator *validator.Validate
	jwtSecret string
}

func NewAuthUseCase(r domain.UserRepository, secret string) *AuthUseCase {
	return &AuthUseCase{
		repo:      r,
		validator: validator.New(),
		jwtSecret: secret,
	}
}

func (uc *AuthUseCase) Create(req dto.CreateUserRequest) (*domain.User, error) {
	existing, err := uc.repo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &domain.User{
		Username: req.Username,
		Password: string(hash),
	}
	id := uuid.New()
	u.ID = id
	_, err = uc.repo.CreateUser(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *AuthUseCase) Login(req dto.LoginRequest) (string, error) {
	user, err := uc.repo.GetByUsername(req.Username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
