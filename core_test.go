package errortrace

import (
	"errors"
	"testing"

	"github.com/techforge-lat/errortrace/errtype"
)

var storageErr = errors.New("psql: could not create user")

var storageErrWithTracing = Wrap(storageErr)
var storageErrWithTracing2 = Wrap(storageErr)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
		want    string
	}{
		{
			name:    "with default error",
			err:     Wrap(storageErr).SetErrCode(errtype.InternalError),
			want:    "\npsql: could not create user\n\ncode:\tinternal_error\n\ngithub.com/techforge-lat/errortrace.TestError_Error:2",
			wantErr: true,
		},
		{
			name:    "with errcontext.Error wrapped",
			err:     Wrap(storageErrWithTracing).SetErrCode(errtype.BadRequest),
			want:    "\npsql: could not create user\n\ncode:\tbad_request\n\ngithub.com/techforge-lat/errortrace.TestError_Error:30\ngithub.com/techforge-lat/errortrace.init:1",
			wantErr: true,
		},
		{
			name: "without wrapped error",
			err:  Wrap(nil).SetErrCode(errtype.BadRequest).SetDetail("Boolean validation failed").Err(),
			want: ``,
		},
		{
			name:    "with title and detail",
			err:     Wrap(storageErrWithTracing2).SetErrCode(errtype.BadRequest).SetTitle("Validation Error").SetDetail("Boolean validation failed"),
			want:    "\npsql: could not create user\n\ntitle:\tValidation Error\ndetail:\tBoolean validation failed\ncode:\tbad_request\n\ngithub.com/techforge-lat/errortrace.TestError_Error:41\ngithub.com/techforge-lat/errortrace.init:1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr && tt.err != nil {
				t.Errorf("got = %q, want = nil", tt.want)
			}

			if !tt.wantErr {
				return
			}

			got := tt.err.Error()
			if tt.wantErr && got != tt.want {
				t.Errorf("got = %v, want = %v", got, tt.want)
			}

		})
	}
}
