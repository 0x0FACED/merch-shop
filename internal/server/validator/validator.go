package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Index int    `json:"index,omitempty"` // Индекс записи в массиве
	Field string `json:"field"`           // Поле, где произошла ошибка
	Tag   string `json:"tag"`             // Тэг валидации, который не прошел
}

type ValidationErrorsResponse struct {
	Errors []*ValidationError `json:"errors"`
}

// Реализуем интерфейс ошибки для структуры ValidationErrorsResponse
func (v *ValidationErrorsResponse) Error() string {
	if len(v.Errors) == 0 {
		return "validation errors: no errors"
	}

	var result string
	for _, err := range v.Errors {
		result += fmt.Sprintf("Index: %d, Field: %s, Tag: %s; ", err.Index, err.Field, err.Tag)
	}

	return fmt.Sprintf("validation errors: %s", result)
}

type APIValidator struct {
	validator *validator.Validate
}

func NewAPIValidator() *APIValidator {
	return &APIValidator{
		validator: validator.New(),
	}
}

func (v *APIValidator) Validate(i any) error {
	if errors := v.ValidateStruct(i); len(errors) > 0 {
		return &ValidationErrorsResponse{Errors: errors}
	}
	return nil
}

func (v *APIValidator) ValidateStruct(i any) []*ValidationError {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	validationErrors := []*ValidationError{}
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			validationErrors = append(validationErrors, &ValidationError{
				Field: fieldErr.StructNamespace(),
				Tag:   fieldErr.Tag(),
			})
		}
	}

	return validationErrors
}
