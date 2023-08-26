package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, Create())
		})
	}
}

func TestUser_Auth(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name     string
		args     args
		authFail bool
	}{
		{
			name:     "valid auth",
			args:     args{username: "asd", password: "asd"},
			authFail: false,
		},
		{
			name:     "invalid auth",
			args:     args{username: "asd", password: "a"},
			authFail: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := Create()
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Authorized Content"))
			})

			req, err := http.NewRequest(http.MethodPost, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.SetBasicAuth(tt.args.username, tt.args.password)

			w := httptest.NewRecorder()
			authMiddleware := mockUser.Auth(mockHandler)
			authMiddleware.ServeHTTP(w, req)

			assert.Equal(t, tt.authFail, w.Code == http.StatusUnauthorized)
		})
	}
}

func TestUser_verifyUserPass(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid",
			args: args{username: "asd", password: "asd"},
			want: true,
		},
		{
			name: "invalid",
			args: args{username: "as", password: "asd"},
			want: false,
		},
		{
			name: "invalid",
			args: args{username: "asd", password: "ad"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := Create()

			assert.Equal(t, tt.want, users.verifyUserPass(tt.args.username, tt.args.password))
		})
	}
}
