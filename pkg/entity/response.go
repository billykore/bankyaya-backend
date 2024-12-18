package entity

import (
	"errors"
	"net/http"
	"time"

	"go.bankyaya.org/app/backend/pkg/datetime"
	"go.bankyaya.org/app/backend/pkg/status"
)

type Response struct {
	Status     string `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	ServerTime string `json:"serverTime,omitempty"`
}

// ResponseSuccess returns status code 200 and success response with data.
func ResponseSuccess(data any) (int, *Response) {
	return http.StatusOK, &Response{
		Status:     "OK",
		Data:       data,
		ServerTime: currentServerTime(),
	}
}

// ResponseSuccessNilData returns status code 200 and success response without data.
func ResponseSuccessNilData() (int, *Response) {
	return http.StatusOK, &Response{
		Status:     "OK",
		ServerTime: currentServerTime(),
	}
}

// ResponseError returns error status code and error.
func ResponseError(err error) (int, *Response) {
	var s *status.Status
	if errors.As(err, &s) {
		return responseCode[s.Code], &Response{
			Status:     responseStatus[s.Code],
			Message:    s.Message,
			ServerTime: currentServerTime(),
		}
	}
	return ResponseInternalServerError(err)
}

// ResponseBadRequest returns status code 400 and error response.
func ResponseBadRequest(err error) (int, *Response) {
	return http.StatusBadRequest, &Response{
		Status:     "BAD_REQUEST",
		Message:    err.Error(),
		ServerTime: currentServerTime(),
	}
}

// ResponseUnauthorized returns status code 401 and error response.
func ResponseUnauthorized(err error) (int, *Response) {
	return http.StatusUnauthorized, &Response{
		Status:     "UNAUTHORIZED",
		Message:    err.Error(),
		ServerTime: currentServerTime(),
	}
}

// ResponseInternalServerError returns status code 500 and error response.
func ResponseInternalServerError(err error) (int, *Response) {
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
