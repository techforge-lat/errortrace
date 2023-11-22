package errortrace

import (
	"errors"
	"testing"

	"github.com/techforge-lat/errortrace/status"
)

var storageErr = errors.New("psql: could not create user")

var storageErrWithTracing = New(storageErr)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
		want    string
	}{
		{
			name:    "with default error",
			err:     New(storageErr).SetStatusCode(status.InternalError),
			want:    `[where=github.com/techforge-lat/errortrace.TestError_Error:23] [status_code=internal_error] [error=psql: could not create user]`,
			wantErr: true,
		},
		{
			name:    "with errcontext.Error wrapped",
			err:     New(storageErrWithTracing).SetStatusCode(status.BadRequest),
			want:    `[where=github.com/techforge-lat/errortrace.TestError_Error:29 => github.com/techforge-lat/errortrace.init:12] [status_code=bad_request] [error=psql: could not create user]`,
			wantErr: true,
		},
		{
			name: "without wrapped error",
			err:  New(nil).SetStatusCode(status.BadRequest).SetPresentationMsg("Boolean validation failed").Err(),
			want: ``,
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
