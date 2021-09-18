package qencode

import (
	"fmt"
	"io"
	"net/http"
)

// RequestError is returned when a request to qEncode
// fails.
type RequestError struct {
	Message      string
	ResponseBody io.Reader
	StatusCode   int
}

// Error implements error interface.
func (r RequestError) Error() string {
	return fmt.Sprintf("[%d %s]: %s", r.StatusCode, http.StatusText(r.StatusCode), r.Message)
}
