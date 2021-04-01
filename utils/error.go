package utils

import (
	"context"
	"net/http"
)

type CustomError struct {
	SendError
	Description string
	ctx         context.Context
}

type SendError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *CustomError) Error() string {
	return err.Description
}

func (err *CustomError) logWarn() {
	WarnLog(err.ctx, Fields{
		"code":        err.Code,
		"message":     err.Message,
		"description": err.Description,
	})
}

func (err *CustomError) logInfo() {
	InfoLog(err.ctx, Fields{
		"code":        err.Code,
		"message":     err.Message,
		"description": err.Description,
	})
}

func NewCustomError(ctx context.Context, code int, mess string) *CustomError {
	err := CustomError{
		SendError: SendError{
			Code:    code,
			Message: mess,
		},
		Description: mess,
		ctx:         ctx,
	}

	err.logWarn()
	return &err
}

func BadRequest(ctx context.Context, desc string) *CustomError {
	err := CustomError{
		SendError: SendError{
			Code:    http.StatusBadRequest,
			Message: "bad request",
		},
		Description: desc,
		ctx:         ctx,
	}

	err.logInfo()
	return &err
}

func FailServer(ctx context.Context, desc string) *CustomError {
	err := CustomError{
		SendError: SendError{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		},
		Description: desc,
		ctx:         ctx,
	}

	err.logWarn()
	return &err
}

func Authorization(ctx context.Context, desc string) *CustomError {
	err := CustomError{
		SendError: SendError{
			Code:    http.StatusUnauthorized,
			Message: "not authorized",
		},
		Description: desc,
		ctx:         ctx,
	}

	err.logInfo()
	return &err
}
