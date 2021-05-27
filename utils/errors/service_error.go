package errors

import "strings"

func CreateErrorWithService(err error) *CustomError {
	custNameError := [9]string{
		"user not found",
		"restaurant not found",
		"User with this email already exists",
		"User with this number already exists",
		"user not authorization",
		"not authorization",
		"Restaurant with this email already exists",
		"Restaurant with this number already exists",
		"Restaurant with this name already exists",
	}

	for _, nameError := range custNameError {
		if strings.Contains(err.Error(), nameError) {
			return NewErrorWithMessage(nameError)
		}
	}

	return FailServerError(err.Error())
}
