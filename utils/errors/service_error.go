package errors

import (
	"context"
	"strings"

	"github.com/borscht/backend/utils/logger"
)

func CreateErrorWithService(ctx context.Context, err error) *CustomError {
	logger.UtilsLevel().InfoLog(ctx, logger.Fields{"ERROR": err.Error()})
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
			custErr := NewErrorWithMessage(nameError)
			logger.UtilsLevel().DebugLog(ctx, logger.Fields{"CUSTOM ERROR": custErr.Error()})
			return custErr
		}
	}

	custErr := FailServerError(err.Error())
	logger.UtilsLevel().DebugLog(ctx, logger.Fields{"CUSTOM ERROR": custErr.Error()})
	return custErr
}
