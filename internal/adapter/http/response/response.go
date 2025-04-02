package response

import (
	"errors"
	"net/http"
	"time"

	"go.bankyaya.org/app/backend/internal/pkg/datetime"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

type Response struct {
	Status     string `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	ServerTime string `json:"serverTime,omitempty"`
}

// Success returns status code 200 and success response with data.
func Success(data any) (int, *Response) {
	return http.StatusOK, &Response{
		Status:     "OK",
		Data:       data,
		ServerTime: currentServerTime(),
	}
}

// SuccessWithoutData returns status code 200 and success response without data.
func SuccessWithoutData() (int, *Response) {
	return http.StatusOK, &Response{
		Status:     "OK",
		ServerTime: currentServerTime(),
	}
}

// Error returns error status code and error.
func Error(err error) (int, *Response) {
	var s *status.Status
	if errors.As(err, &s) {
		return responseCode[s.Code], &Response{
			Status:     responseStatus[s.Code],
			Message:    s.Error(),
			ServerTime: currentServerTime(),
		}
	}
	return InternalServerError(err)
}

// BadRequest returns status code 400 and error response.
func BadRequest(err error) (int, *Response) {
	return http.StatusBadRequest, &Response{
		Status:     "BAD_REQUEST",
		Message:    err.Error(),
		ServerTime: currentServerTime(),
	}
}

// Unauthorized returns status code 401 and error response.
func Unauthorized(err error) (int, *Response) {
	return http.StatusUnauthorized, &Response{
		Status:     "UNAUTHORIZED",
		Message:    err.Error(),
		ServerTime: currentServerTime(),
	}
}

// InternalServerError returns status code 500 and error response.
func InternalServerError(err error) (int, *Response) {
	return http.StatusInternalServerError, &Response{
		Status:     "INTERNAL_SERVER_ERROR",
		Message:    err.Error(),
		ServerTime: currentServerTime(),
	}
}

// currentServerTime provides the current server time using the default time layout.
func currentServerTime() string {
	return time.Now().Format(datetime.DefaultTimeLayout)
}

var responseCode = []int{
	http.StatusOK,
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusForbidden,
	http.StatusNotFound,
	http.StatusConflict,
	http.StatusInternalServerError,
}

var responseStatus = []string{
	"OK",
	"BAD_REQUEST",
	"UNAUTHORIZED",
	"FORBIDDEN",
	"NOT_FOUND",
	"CONFLICT",
	"INTERNAL_SERVER_ERROR",
}
