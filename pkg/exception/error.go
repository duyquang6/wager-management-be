package exception

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	_validator "github.com/duyquang6/wager-management-be/pkg/validator"
	"strings"
)

const (
	validationErrorFmt = "Field:%s Error:%s"
)

// AppError describes application error.
type AppError struct {
	Meta          ErrorMeta `json:"meta"`
	OriginalError error     `json:"-"`
}

// AppErrorResponse describes return error response format.
type AppErrorResponse struct {
	Error string `json:"error"`
	//Code  int    `json:"code"`
}

// ErrorMeta is the metadata of AppError.
type ErrorMeta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (appErr AppError) Error() string {
	if appErr.OriginalError != nil {
		return appErr.OriginalError.Error()
	}
	return appErr.Meta.Message
}

// New returns an AppError with args.
func New(errCode int, msg string) error {
	return AppError{
		Meta: ErrorMeta{
			Code:    errCode,
			Message: msg,
		},
		OriginalError: nil,
	}
}

// Newf returns an AppError with args and message.
func Newf(errCode int, template string, args ...interface{}) error {
	return AppError{
		Meta: ErrorMeta{
			Code:    errCode,
			Message: fmt.Sprintf(template, args...),
		},
		OriginalError: nil,
	}
}

// Wrap returns an AppError with err, args.
func Wrap(errCode int, err error, msg string) error {
	return AppError{
		Meta: ErrorMeta{
			Code:    errCode,
			Message: msg,
		},
		OriginalError: err,
	}
}

// Wrapf returns an AppError with err, args and message.
func Wrapf(errCode int, err error, template string, args ...interface{}) error {
	return AppError{
		Meta: ErrorMeta{
			Code:    errCode,
			Message: fmt.Sprintf(template, args...),
		},
		OriginalError: err,
	}
}

func (appErr AppError) GetHTTPStatusCode() int {
	if appErr.Meta.Code < 1000 {
		return appErr.Meta.Code
	}
	return appErr.Meta.Code / 10000
}

func (appErr AppError) AggregateMetaError() AppError {
	res := []string{}

	var err error
	err = appErr
	for err != nil {
		if _appErr, ok := err.(AppError); ok {
			res = append(res, _appErr.Meta.Message)
			err = _appErr.OriginalError
			continue
		}
		res = append(res, err.Error())
		break
	}
	appErr.Meta.Message = strings.Join(res, "\n")
	return appErr
}

func (appErr AppError) ToAppErrorResponse() AppErrorResponse {
	var ve validator.ValidationErrors
	if errors.As(appErr.OriginalError, &ve) {
		out := make([]string, len(ve))
		for i, fe := range ve {
			out[i] = fmt.Sprintf(validationErrorFmt, fe.Field(), _validator.MsgForTag(fe))
		}
		return AppErrorResponse{
			strings.Join(out, "\n"),
		}
	}
	return AppErrorResponse{
		//Code:  appErr.Meta.Code,
		Error: appErr.AggregateMetaError().Meta.Message,
	}
}
