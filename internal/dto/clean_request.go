package dto

type CleanUpRequest struct {
	Url       string `json:"url" validate:"required"`
	Operation string `json:"operation" validate:"required"`
}
