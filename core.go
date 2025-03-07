package errortrace

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/techforge-lat/errortrace/v2/errtype"
)

type Error struct {
	Title   string  `json:"title,omitempty"`
	Message string  `json:"message"`
	Code    string  `json:"code"`
	Cause   error   `json:"cause,omitempty"`
	Stack   []Frame `json:"stack,omitempty"`
}

type Frame struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// OnError starts a new error builder chain if the cause error is not nil.
// Returns nil if causeError is nil.
func OnError(causeError error) *Error {
	if causeError == nil {
		return nil
	}

	var errTrace *Error
	if errors.As(causeError, &errTrace) {
		errTrace.Stack = append(errTrace.Stack, captureStack())

		return errTrace
	}

	return &Error{
		Cause: causeError,
		Stack: []Frame{captureStack()},
	}
}

// Code sets the error code
func (b *Error) WithCode(code errtype.Code) *Error {
	b.Code = string(code)

	return b
}

// Message sets the error message
func (b *Error) WithMessage(msg string) *Error {
	b.Message = msg

	return b
}

// From adds a cause to the error
func (b *Error) From(cause error) *Error {
	b.Cause = cause

	return b
}

func (b *Error) HasTitle() bool {
	return b.Title != ""
}

func (b *Error) WithTitle(title string) *Error {
	b.Title = title

	return b
}

// Error implements the error interface with a logging-friendly format
func (b *Error) Error() string {
	var parts []string

	if len(b.Stack) > 0 {
		var stackPaths []string
		// Stack is already in correct order, just format it
		for _, frame := range b.Stack {
			file := filepath.Base(frame.File)
			dir := filepath.Base(filepath.Dir(frame.File))
			location := fmt.Sprintf("%s/%s:%d", dir, file, frame.Line)
			stackPaths = append(stackPaths, location)
		}
		parts = append(parts, fmt.Sprintf("[stack=%s]", strings.Join(stackPaths, " => ")))
	}

	if b.Code != "" {
		parts = append(parts, fmt.Sprintf("[code=%s]", strings.ToLower(b.Code)))
	}

	msg := b.Message
	if b.Cause != nil {
		msg = fmt.Sprintf("%s: %v", b.Message, b.Cause)
	}
	parts = append(parts, fmt.Sprintf("[error=%s]", msg))

	return strings.Join(parts, " ")
}

// captureStack now stores frames in reverse order
func captureStack() Frame {
	fn, file, line, _ := runtime.Caller(2)

	fullFuncName := runtime.FuncForPC(fn).Name()

	// Extract just the function name from the full path
	if lastDot := strings.LastIndex(fullFuncName, "."); lastDot != -1 {
		fullFuncName = fullFuncName[lastDot+1:]
	}

	return Frame{
		File:     file,
		Line:     line,
		Function: fullFuncName,
	}
}
