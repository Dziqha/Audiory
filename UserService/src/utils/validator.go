package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func ValidateStruct(s interface{}) error {
	val := validator.New()

	if err := val.Struct(s); err != nil {
		return errors.Wrap(err, "validation failed")
	}
	return nil
}
