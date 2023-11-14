package errortrace

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/techforge-lat/errortrace/status"
)

type Error struct {
	err             error
	statusCode      status.Code
	presentationMsg string
	where           string
	metadata        map[string]any
}

func (e *Error) HasErr() bool {
	return e.err != nil
}

func (e *Error) HasStatusCode() bool {
	return e.statusCode != ""
}

func (e *Error) HasPresentationMsg() bool {
	return e.presentationMsg != ""
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

	stringBuilder.WriteString(fmt.Sprintf("[where=%s] ", e.Where()))

	metadata := e.AggregateMetadata()
	if metadata == nil {
		metadata = make(map[string]any)
	}

	statusCode := e.StatusCode()
	if statusCode != "" {
		stringBuilder.WriteString(fmt.Sprintf("[status_code=%s] ", statusCode))
	}

	presentationMsg := e.PresentationMsg()
	if presentationMsg != "" {
		stringBuilder.WriteString(fmt.Sprintf("[presentation_msg=%s] ", presentationMsg))
	}

	if errStr != "" {
		stringBuilder.WriteString(fmt.Sprintf("[error=%s] ", errStr))
	}

	for _, key := range GetSortedMetadataKeys(metadata) {
		value := metadata[key]

		valueStr, ok := value.(string)
		if !ok {
			stringBuilder.WriteString(fmt.Sprintf("[%s=%v] ", key, value))
			continue
		}

		if valueStr == "" {
			continue
		}

		stringBuilder.WriteString(fmt.Sprintf("[%s=%s] ", key, valueStr))
	}

	return stringBuilder.String()[:stringBuilder.Len()-1]
}

// Err returns the first non *Error type
func (e *Error) Err() error {
	var customErr *Error
	if !errors.As(e.err, &customErr) {
		return e.err
	}

	return customErr.Err()
}

func (e *Error) SetStatusCode(t status.Code) *Error {
	e.statusCode = t
	return e
}

// StatusCode returns the first status code in the *Error chain
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

// PresentationMsg returns the first PresentationMsg in the *Error chain
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

func (e *Error) AddMetadata(key string, value any) *Error {
	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}

	e.metadata[key] = value
	return e
}

// AggregateMetadata reeturns every metadata in every *Error in the chain
func (e *Error) AggregateMetadata() map[string]any {
	var customErr *Error
	if !errors.As(e.err, &customErr) {
		return e.metadata
	}

	customErr.AggregateMetadata()
	for k, v := range customErr.metadata {
		e.AddMetadata(k, v)
	}

	return e.metadata
}

// Where returns the first `where` in the *Error chain
func (e *Error) Where() string {
	var wrappedErr *Error
	if errors.As(e.err, &wrappedErr) {
		return fmt.Sprintf("%s => %s", e.where, wrappedErr.Where())
	}

	return e.where
}

func GetSortedMetadataKeys(metadata map[string]any) []string {
	keys := []string{}
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
