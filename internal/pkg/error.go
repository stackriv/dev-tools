package pkg

import (
	"net/http"
	"strconv"
)

// ErrorMessage	 returns the good message by code
func ErrorMessage(code int) map[string]string {
	var msg string
	switch code {
	case http.StatusBadRequest:
		msg = "Bad Request"
	case http.StatusNotFound:
		msg = "Not Found"
	case http.StatusMethodNotAllowed:
		msg = "Method Not Allowed"
	case http.StatusInternalServerError:
		msg = "Internal Server Error"
	case http.StatusConflict:
		msg = "Conflict"
	case http.StatusForbidden:
		msg = "Forbidden"
	case http.StatusUnauthorized:
		msg = "Unauthorized"
	}

	return map[string]string{
		"code": strconv.Itoa(code),
		"msg":  msg,
	}
}
