package common

import (
	"github.com/micro/go-micro/errors"
)

// BadRequest generates a 400 error.
func BadRequest(id string, fun interface{}, err error, format string, a ...interface{}) error {
	ErrorLog(id, fun, err, format)
	return errors.BadRequest(id, format, a)
}

// Unauthorized generates a 401 error.
func Unauthorized(id string, fun interface{}, err error, format string, a ...interface{}) error {
	ErrorLog(id, fun, err, format)
	return errors.Unauthorized(id, format, a)
}

// Forbidden generates a 403 error.
func Forbidden(id string, fun interface{}, err error, format string, a ...interface{}) error {
	ErrorLog(id, fun, err, format)
	return errors.Forbidden(id, format, a)
}

// NotFound generates a 404 error.
func NotFound(id string, fun interface{}, err error, format string, a ...interface{}) error {
	ErrorLog(id, fun, err, format)
	return errors.NotFound(id, format, a)
}

// InternalServerError generates a 500 error.
func InternalServerError(id string, fun interface{}, err error, format string, a ...interface{}) error {
	ErrorLog(id, fun, err, format)
	return errors.InternalServerError(id, format, a)
}
