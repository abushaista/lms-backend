package dto

type CategoryRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
