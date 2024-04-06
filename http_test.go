package errs_test

import (
	"github.com/valensto/errs"
	"testing"
)

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "nil",
			err:  nil,
			want: 200,
		},
		{
			name: "not found",
			err:  errs.New(errs.SlugNotFound),
			want: 404,
		},
		{
			name: "invalid",
			err:  errs.New(errs.SlugInvalid),
			want: 400,
		},
		{
			name: "internal",
			err:  errs.New(errs.SlugInternal),
			want: 500,
		},
		{
			name: "unknown",
			err:  errs.New(errs.SlugUnknown),
			want: 500,
		},
		{
			name: "unauthorized",
			err:  errs.New(errs.SlugUnauthorized),
			want: 401,
		},
		{
			name: "forbidden",
			err:  errs.New(errs.SlugForbidden),
			want: 403,
		},
		{
			name: "duplicate",
			err:  errs.New(errs.SlugDuplicate),
			want: 409,
		},
		{
			name: "not implemented",
			err:  errs.New(errs.SlugNotImplemented),
			want: 501,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errs.HTTPStatus(tt.err); got != tt.want {
				t.Errorf("HTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
