package custom_validator

import (
	"auth-service/internal/custom_error"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() CustomValidator {
	return CustomValidator{validator.New()}
}

func (cv CustomValidator) Validate(dto interface{}) *custom_error.AppError {
	fmt.Println("Validation dto: " + fmt.Sprintf("%v", dto))
	err := cv.validator.Struct(dto)
	if err != nil {
		fmt.Println("Validation error")
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if invalidValidationError, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			fmt.Println(invalidValidationError)
			return nil
		}

		appErr := &custom_error.AppError{Message: "Validation is failed: ", Err: err, HttpErrorCode: 400}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("Validation error field")
			appErr.Message += "Field name: `" + err.Field() + "` and field value: `" + fmt.Sprintf("%v", err.Value()) + "`."
		}

		return appErr
	}
	fmt.Println("Validation is successful")
	return nil
}
