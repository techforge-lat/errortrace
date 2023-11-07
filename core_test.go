package errortrace

import (
	"errors"
	"testing"
)

var storageErr = errors.New("psql: could not create user")

var storageErrWithTracing = New(errors.New("psql: cound not create user"))

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "with default error",
			err:  New(storageErr).SetStatusCode(InternalError),
			want: `error = [where: "github.com/techforge-lat/errortrace.TestError_Error:16"], [statusCode: "internal_error"], [wrappedErr: psql: could not create user]`,
		},
		{
			name: "with errcontext.Error wrapped",
			err:  New(storageErrWithTracing).SetStatusCode(BadRequest),
			want: `error = [where: "github.com/techforge-lat/errortrace.TestError_Error:21"], [statusCode: "bad_request"], [wrappedErr: error = [where: "github.com/techforge-lat/errortrace.TestError_Error:21"], [wrappedErr: psql: cound not create user]]`,
		},
		{
			name: "without wrapped error",
			err:  New(nil).SetStatusCode(BadRequest).SetPresentationMsg("Boolean validation failed"),
			want: `error = [where: "github.com/techforge-lat/errortrace.TestError_Error:26"], [statusCode: "bad_request"], [presentationMsg: "Boolean validation failed"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}
		})
	}
}
