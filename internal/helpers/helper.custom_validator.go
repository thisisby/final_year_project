package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"strings"
)

func BindAndValidate(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {

		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return fmt.Errorf("error parsing validation tags: %v", invalidValidationError)
		}

		var errorsOut validationErrors
		for _, err := range err.(validator.ValidationErrors) {
			var e error
			switch err.Tag() {
			case "required":
				e = fmt.Errorf("field '%s' cannot be blank", err.Field())
			case "email":
				e = fmt.Errorf("field '%s' must be a valid email address", err.Field())
			case "eth_addr":
				e = fmt.Errorf("field '%s' must  be a valid Ethereum address", err.Field())
			case "len":
				e = fmt.Errorf("field '%s' must be exactly %v characters long", err.Field(), err.Param())
			default:
				e = fmt.Errorf("field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
			}
			errorsOut = append(errorsOut, e)
		}

		return errorsOut
	}

	return nil
}

type validationErrors []error

func (v validationErrors) Error() string {

	buff := bytes.NewBufferString("")

	for i := 0; i < len(v); i++ {

		buff.WriteString(v[i].Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
