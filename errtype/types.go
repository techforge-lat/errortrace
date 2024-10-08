package errtype

import "net/http"

type Code string

const (
	BadRequest          Code = "bad_request"
	UnprocessableEntity Code = "unprocessable_entity"
	InternalError       Code = "internal_error"
	BindFailed          Code = "bind_failed"
	Unauthorized        Code = "unauthorized"
	Forbidden           Code = "forbidden"
	NotFound            Code = "not_found"
)

var HTTPStatusByStatusCode = map[Code]int{
	BadRequest:          http.StatusBadRequest,
	UnprocessableEntity: http.StatusUnprocessableEntity,
	InternalError:       http.StatusInternalServerError,
	BindFailed:          http.StatusBadRequest,
	Unauthorized:        http.StatusUnauthorized,
	Forbidden:           http.StatusForbidden,
	NotFound:            http.StatusNotFound,
}
