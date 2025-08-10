package dto

type CategoryRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required"`
}
