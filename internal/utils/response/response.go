package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error, statusCode int) Response {
	return Response{

		Status:  statusCode,
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}
}

func GeneralResponse(msg string, statusCode int, data interface{}) Response {
	return Response{
		Status:  statusCode,
		Success: true,
		Message: msg,
		Data:    data,
	}
}

func ValidationError(errs validator.ValidationErrors, statusCode int) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}

	}

	return Response{
		Status:  statusCode,
		Success: false,
		Message: strings.Join(errMsgs, ", "),
		Data:    nil,
	}
}
