package runscope

import (
	"fmt"
	"net/http"
)

type Error struct {
	Response *http.Response
	E        struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e Error) Status() int {
	if e.E.Status != 0 {
		return e.E.Status
	}
	return e.Response.StatusCode
}

func (e Error) Error() string {
	message := e.E.Message
	if e.E.Message == "" {
		message = e.Response.Status
	}

	return fmt.Sprintf("%d %s", e.Status(), message)
}
