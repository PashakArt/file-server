package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateStruct(s any) error {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	var invalidErr *validator.InvalidValidationError
	if errors.As(err, &invalidErr) {
		return fmt.Errorf("internal validator error: %w", err)
	}

	var errMsgs []string
	for _, err := range err.(validator.ValidationErrors) {
		errMsgs = append(errMsgs, fmt.Sprintf("field '%s' failed on '%s'", err.Field(), err.Tag()))
	}

	return errors.New(strings.Join(errMsgs, "; "))
}
