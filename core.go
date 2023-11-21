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
	// err is our root causing err, if it doesn't exist
	// there is nothing to trace
	if err == nil {
		return nil
	}

	fun, _, line, _ := runtime.Caller(1)
	where := fmt.Sprintf("%s:%d", runtime.FuncForPC(fun).Name(), line)

	customeErr := &Error{}
	if errors.As(err, &customeErr) {
		customeErr.where = fmt.Sprintf("%s => %s", where, customeErr.where)
		return customeErr
	}

	e := &Error{
		err:   err,
		where: fmt.Sprintf("%s:%d", runtime.FuncForPC(fun).Name(), line),
	}

	return e
}

func (e *Error) SetErr(err error) *Error {
	if e == nil {
		e = &Error{}
	}

	e.err = err
	return e
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	var stringBuilder strings.Builder
	var errStr string

	err := e.Err()
	if err != nil {
		errStr = err.Error()
	}

	stringBuilder.WriteString(fmt.Sprintf("[where=%s] ", e.Where()))

	metadata := e.Metadata()
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

// Err returns the causing error of the trace chain
func (e *Error) Err() error {
	if e == nil {
		return nil
	}

	return e.err
}

func (e *Error) SetStatusCode(t status.Code) *Error {
	if e == nil {
		return nil
	}

	e.statusCode = t
	return e
}

// StatusCode returns the last status in the trace chain
func (e *Error) StatusCode() string {
	if e == nil {
		return ""
	}

	return string(e.statusCode)
}

func (e *Error) SetPresentationMsg(msg string) *Error {
	if e == nil {
		return nil
	}

	e.presentationMsg = msg
	return e
}

// PresentationMsg returns the last PresentationMsg in the trace chain chain
func (e *Error) PresentationMsg() string {
	if e == nil {
		return ""
	}

	return e.presentationMsg
}

func (e *Error) AddMetadata(key string, value any) *Error {
	if e == nil {
		return nil
	}

	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}

	e.metadata[key] = value
	return e
}

// Metadata reeturns all metadata in the trace chain
func (e *Error) Metadata() map[string]any {
	if e == nil {
		return nil
	}

	return e.metadata
}

// Where returns the where trace chain
func (e *Error) Where() string {
	if e == nil {
		return ""
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
