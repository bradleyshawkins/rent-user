package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bradleyshawkins/berror"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Errors interface {
	Message() string
	ErrorCode() string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Unable to execute handler. Code: %d, Message: %s", e.Code, e.Message)
}

func ErrorHandler(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log.Printf("Unexpected Error: %v", err)
			statusCode := http.StatusInternalServerError
			msg := "An unknown error occurred"
			code := berror.CodeUnknown
			re, ok := err.(*berror.Error)
			if ok {
				statusCode = re.HttpStatusCode()
				msg = re.UserMessage()
				code = re.Code()
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_ = json.NewEncoder(w).Encode(Error{
				Message: msg,
				Code:    int(code),
			})
		}
	}
}
