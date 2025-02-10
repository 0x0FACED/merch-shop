package validator

import (
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Index int    `json:"index,omitempty"` // Индекс записи в массиве
	Field string `json:"field"`           // Поле, где произошла ошибка
	Tag   string `json:"tag"`             // Тэг валидации, который не прошел
	//Err   string `json:"error,omitempty"` // Нормальное описание ошибки (чтоб понятно было)
}

type ValidationErrorsResponse struct {
	Errors []*ValidationError `json:"errors"`
}

// Реализуем интерфейс ошибки для структуры ValidationErrorsResponse
func (v *ValidationErrorsResponse) Error() string {
	return "validation errors"
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
				//Err:   fmt.Sprintf("Validation failed on '%s' for field '%s'", fieldErr.Tag(), fieldErr.Field()),
			})
		}
	}

	return validationErrors
}

/*
func (v *APIValidator) ValidateArray(arr []*CreateRawEmailRequest) []*ValidationError {
	var errors []*ValidationError

	for i, item := range arr {
		errs := v.ValidateStruct(item)
		for _, e := range errs {
			e.Index = i
			errors = append(errors, e)
		}
	}

	return errors
}
*/
