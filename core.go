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
	metadata        map[string]string
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
	strBuilder := strings.Builder{}

	if e.statusCode != "" {
		strBuilder.WriteString(fmt.Sprintf("\n\nstatus: %q\n", e.statusCode))
	}

	strBuilder.WriteString(fmt.Sprintf("where: %q\n", e.Where()))

	for _, k := range e.getSortedMetadataKeys() {
		v := e.metadata[k]
		strBuilder.WriteString(fmt.Sprintf("%s: %q\n", k, v))
	}

	if e.presentationMsg != "" {
		strBuilder.WriteString(fmt.Sprintf("presentationMsg: %q\n", e.presentationMsg))
	}

	if e.err != nil {
		strBuilder.WriteString(fmt.Sprintf("\t%s\n\n", e.err))
	}

	return strBuilder.String()
}

func (e *Error) ErrorWithContext() string {
	return fmt.Sprintf("[%s], [%s]", e.Where(), e.err)
}

func (e *Error) SetStatusCode(t StatusCode) *Error {
	e.statusCode = t
	return e
}

func (e *Error) StatusCode() string {
	return string(e.statusCode)
}

func (e *Error) SetPresentationMsg(msg string) *Error {
	e.presentationMsg = msg
	return e
}

func (e *Error) PresentationMsg() string {
	return e.presentationMsg
}

func (e *Error) AddMetadata(key, value string) *Error {
	e.metadata[key] = value
	return e
}

func (e *Error) Where() string {
	return e.where
}

func whereChain(e *Error) string {
	wrappedErr := &Error{}
	if errors.As(e.err, &wrappedErr) {
		return fmt.Sprintf("%s => %s", e.where, whereChain(wrappedErr))
	}

	return e.where
}

func (e *Error) getSortedMetadataKeys() []string {
	keys := []string{}
	for key := range e.metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
