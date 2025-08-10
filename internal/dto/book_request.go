package dto

type CreateBookRequest struct {
	Title      string `json:"title" validate:"required"`
	Author     string `json:"author" validate:"required"`
	ISBN       string `json:"isbn" validate:"required"`
	Year       int    `json:"year" validate:"required"`
	Summary    string `json:"summary" validate:"required"`
	CoverImage string `json:"cover_image" validate:"url"`
	CategoryID uint   `json:"category_id" validate:"required"`
}

type UpdateBookRequest struct {
	ID         int64  `json:"id"`
	Title      string `json:"title" validate:"required"`
	Author     string `json:"author" validate:"required"`
	ISBN       string `json:"isbn" validate:"required"`
	Year       int    `json:"year" validate:"required"`
	Summary    string `json:"summary" validate:"required"`
	CoverImage string `json:"cover_image" validate:"url"`
	CategoryID uint   `json:"category_id" validate:"required"`
}
