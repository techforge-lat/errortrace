package errortrace

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strings"
)

type StatusCode string

const (
	BadRequest    StatusCode = "bad_request"
	InternalError StatusCode = "internal_error"
)

var HTTPStatusByStatusCode = map[StatusCode]int{
	BadRequest:    http.StatusBadRequest,
	InternalError: http.StatusInternalServerError,
}

type Error struct {
	err             error
	statusCode      StatusCode
	presentationMsg string
	where           string
	metadata        map[string]any
}

func New(err error) *Error {
	fun, _, line, _ := runtime.Caller(1)

	e := &Error{
		err:   err,
		where: fmt.Sprintf("%s:%d", runtime.FuncForPC(fun).Name(), line),
	}

	return e
}

func (e *Error) Error() string {
	var stringBuilder strings.Builder
	var errStr string

	err := e.Err()
	if err != nil {
		errStr = err.Error()
	}

	stringBuilder.WriteString(fmt.Sprintf("[where=%q] ", e.AgregateWhere()))

	metadata := make(map[string]any, 0)
	e.AgregateMetadata(metadata)

	metadata["status_code"] = e.StatusCode()
	metadata["presentation_msg"] = e.PresentationMsg()
	metadata["error"] = errStr

	for _, key := range getSortedMetadataKeys(metadata) {
		value := metadata[key]

		valueStr, ok := value.(string)
		if !ok {
			stringBuilder.WriteString(fmt.Sprintf("[%s=%v] ", key, value))
			continue
		}

		if valueStr == "" {
			continue
		}

		stringBuilder.WriteString(fmt.Sprintf("[%s=%q] ", key, valueStr))
	}

	return stringBuilder.String()[:stringBuilder.Len()-1]
}

func (e *Error) SetStatusCode(t StatusCode) *Error {
	e.statusCode = t
	return e
}

// Err returns the first non *Error type
func (e *Error) Err() error {
	var customErr *Error
	if !errors.As(e.err, &customErr) {
		return e.err
	}

	return customErr.Err()
}

func (e *Error) StatusCode() string {
	if e.statusCode != "" {
		return string(e.statusCode)
	}

	var customErr *Error
	if errors.As(e.err, &customErr) {
		return customErr.StatusCode()
	}

	return ""
}

func (e *Error) SetPresentationMsg(msg string) *Error {
	e.presentationMsg = msg
	return e
}

func (e *Error) PresentationMsg() string {
	if e.presentationMsg != "" {
		return e.presentationMsg
	}

	var customErr *Error
	if errors.As(e.err, &customErr) {
		return customErr.PresentationMsg()
	}

	return ""
}

func (e *Error) AddMetadata(key, value string) *Error {
	e.metadata[key] = value
	return e
}

func (e *Error) AgregateMetadata(metadata map[string]any) {
	var customErr *Error
	if !errors.As(e.err, &customErr) {
		return
	}

	for k, v := range customErr.metadata {
		metadata[k] = v
	}

	customErr.AgregateMetadata(metadata)
}

func (e *Error) Where() string {
	return e.where
}

func (e *Error) AgregateWhere() string {
	var wrappedErr *Error
	if errors.As(e.err, &wrappedErr) {
		return fmt.Sprintf("%s => %s", e.where, wrappedErr.AgregateWhere())
	}

	return e.where
}

func getSortedMetadataKeys(metadata map[string]any) []string {
	keys := []string{}
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
