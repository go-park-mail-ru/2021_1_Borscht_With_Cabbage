package errors

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
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

func (err *CustomError) logError() {
	logrus.WithFields(logrus.Fields{
		"code":        err.Code,
		"message":     err.Message,
		"description": err.Description,
	}).Warn("request_id: ", err.ctx.Value("request_id"))
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

	err.logError()
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

	err.logError()
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

	err.logError()
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

	err.logError()
	return &err
}
