package errortrace

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/techforge-lat/errortrace/errtype"
)

type Error struct {
	err      error
	code     errtype.Code
	title    string
	detail   string
	where    string
	metadata map[string]any
}

// Wrap wraps an error with tracing information
func Wrap(err error) *Error {
	fun, _, line, _ := runtime.Caller(1)
	where := fmt.Sprintf("%s:%d", runtime.FuncForPC(fun).Name(), line)

	customeErr := &Error{}
	if errors.As(err, &customeErr) {
		customeErr.where = fmt.Sprintf("%s\n%s", where, customeErr.where)
		return customeErr
	}

	e := &Error{
		err:   err,
		where: fmt.Sprintf("%s:%d", runtime.FuncForPC(fun).Name(), line),
	}

	return e
}

func (e *Error) HasErr() bool {
	return e.err != nil
}

func (e *Error) SetErr(err error) *Error {
	e.err = err
	return e
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

func (e *Error) HasStatusCode() bool {
	return e.code != ""
}

func (e *Error) SetErrCode(t errtype.Code) *Error {
	e.code = t
	return e
}

// ErrCode returns the last errtype in the trace chain
func (e *Error) ErrCode() string {
	return string(e.code)
}

func (e *Error) HasTitle() bool {
	return e.title != ""
}

func (e *Error) SetTitle(title string) *Error {
	e.title = title
	return e
}

func (e *Error) Title() string {
	return e.title
}

func (e *Error) HasDetail() bool {
	return e.detail != ""
}

func (e *Error) SetDetail(msg string) *Error {
	e.detail = msg
	return e
}

// Detail returns the last Detail in the trace chain chain
func (e *Error) Detail() string {
	return e.detail
}

// Where returns the where trace chain
func (e *Error) Where() string {
	return e.where
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

func (e *Error) Error() string {
	var stringBuilder strings.Builder
	var errStr string

	err := e.RootErr()
	if err != nil {
		errStr = err.Error()
	}

	metadata := e.Metadata()
	if metadata == nil {
		metadata = make(map[string]any)
	}

	if !isEmpty(errStr) {
		stringBuilder.WriteString(fmt.Sprintf("\n%s\n", errStr))
	}

	title := e.Title()
	if !isEmpty(title) {
		stringBuilder.WriteString(fmt.Sprintf("\ntitle:\t%s", title))
	}

	detail := e.Detail()
	if !isEmpty(detail) {
		stringBuilder.WriteString(fmt.Sprintf("\ndetail:\t%s", detail))
	}

	errtypeCode := e.ErrCode()
	if !isEmpty(errtypeCode) {
		stringBuilder.WriteString(fmt.Sprintf("\ncode:\t%s", errtypeCode))
	}

	for _, key := range GetSortedMetadataKeys(metadata) {
		value := metadata[key]

		valueStr, ok := value.(string)
		if !ok {
			stringBuilder.WriteString(fmt.Sprintf("\n[%s:\t%v ", key, value))
			continue
		}

		if isEmpty(valueStr) {
			continue
		}

		stringBuilder.WriteString(fmt.Sprintf("\n%s:\t%s ", key, valueStr))
	}

	stringBuilder.WriteString(fmt.Sprintf("\n\n%s", e.Where()))

	return stringBuilder.String()[:stringBuilder.Len()-1]
}

func GetSortedMetadataKeys(metadata map[string]any) []string {
	keys := []string{}
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

func isEmpty(s string) bool {
	return s == ""
}
