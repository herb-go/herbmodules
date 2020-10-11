package auth

import (
	"net/http"
)

//Session session interface
type Session interface {
	// Set set session by field name with given value.
	Set(r *http.Request, fieldname string, v interface{}) error
	//Get get session by field name with given value.
	Get(r *http.Request, fieldname string, v interface{}) error
	// Del del session value by field name .
	Del(r *http.Request, fieldname string) error
	// IsNotFoundError check if given error is session not found error.
	IsNotFoundError(err error) bool
}
