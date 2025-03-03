package model

import "github.com/go-playground/validator/v10"

type Todos struct {
	ID          int    `gorm:"primary_key" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (ctr *CreateTodoRequest) Validate() error {
	return validate.Struct(ctr)
}
