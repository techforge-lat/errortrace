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
	e.err = err
	return e
}

func (e *Error) Error() string {
	var stringBuilder strings.Builder
	var errStr string

	err := e.RootErr()
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

func (e *Error) RootErr() error {
	return e.err
}

// Err returns the causing error of the trace chain
func (e *Error) Err() error {
	if e.err == nil {
		return nil
	}

	return e
}

func (e *Error) SetStatusCode(t status.Code) *Error {
	e.statusCode = t
	return e
}

// StatusCode returns the last status in the trace chain
func (e *Error) StatusCode() string {
	return string(e.statusCode)
}

func (e *Error) SetPresentationMsg(msg string) *Error {
	e.presentationMsg = msg
	return e
}

// PresentationMsg returns the last PresentationMsg in the trace chain chain
func (e *Error) PresentationMsg() string {
	return e.presentationMsg
}

func (e *Error) AddMetadata(key string, value any) *Error {
	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}

	e.metadata[key] = value
	return e
}

// Metadata reeturns all metadata in the trace chain
func (e *Error) Metadata() map[string]any {
	return e.metadata
}

// Where returns the where trace chain
func (e *Error) Where() string {
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
