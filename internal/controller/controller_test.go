package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erupshis/kode.git/internal/config"
	"github.com/erupshis/kode.git/internal/logger"
	"github.com/erupshis/kode.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_controller_postHandler(t *testing.T) {
	cfg := config.Config{
		Host: "localhost:8080",
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	//defer log.Sync()

	storageRam, err := storage.CreateRamStorage()
	require.NoError(t, err)
	ts := httptest.NewServer(Create(log, storageRam).Route())
	defer ts.Close()

	type reqJSON struct {
		method   string
		username string
		password string
		url      string
		body     string
	}
	type wantJSON struct {
		code int
		body string
	}

	tests := []struct {
		name string
		req  reqJSON
		want wantJSON
	}{
		{
			"add new text for auth user",
			reqJSON{
				method:   http.MethodPost,
				username: "asd",
				password: "asd",
				url:      "/",
				body:     `{"data":"some_text"}`,
			},
			wantJSON{
				code: http.StatusOK,
			},
		},
		{
			"add second text for auth user",
			reqJSON{
				method:   http.MethodPost,
				username: "asd",
				password: "asd",
				url:      "/",
				body:     `{"data":"some_text"}`,
			},
			wantJSON{
				code: http.StatusOK,
			},
		},
		{
			"add third text for auth user",
			reqJSON{
				method:   http.MethodPost,
				username: "asd",
				password: "asd",
				url:      "/",
				body:     `{"data":"some_text"}`,
			},
			wantJSON{
				code: http.StatusOK,
			},
		},
		{
			"add third text for unauth user",
			reqJSON{
				method:   http.MethodPost,
				username: "asd",
				password: "",
				url:      "/",
				body:     `{"data":"some_text"}`,
			},
			wantJSON{
				code: http.StatusUnauthorized,
				body: "Unauthorized\n",
			},
		},
		{
			"add text for unauth user with wrong url",
			reqJSON{
				method:   http.MethodPost,
				username: "asd",
				password: "",
				url:      "/xfdfg",
				body:     `{"data":"some_text"}`,
			},
			wantJSON{
				code: http.StatusUnauthorized,
				body: "Unauthorized\n",
			},
		},
		{
			"add text for auth user with wrong url",
			reqJSON{
				method:   http.MethodPost,
				username: "asd",
				password: "asd",
				url:      "/xfdfg",
				body:     `{"data":"some_text"}`,
			},
			wantJSON{
				code: http.StatusNotFound,
				body: "404 page not found\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(tt.req.body)
			req, errReq := http.NewRequest(tt.req.method, ts.URL+tt.req.url, body)
			require.NoError(t, errReq)

			req.Header.Add("Content-Type", "application/json")
			req.SetBasicAuth(tt.req.username, tt.req.password)
			resp, errResp := ts.Client().Do(req)
			assert.NoError(t, errResp)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			//assert.Equal(t, tt.req.method, )
			assert.Equal(t, tt.want.body, string(respBody))
			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}

func Test_controller_getHandler(t *testing.T) {
	cfg := config.Config{
		Host: "localhost:8080",
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	//defer log.Sync()

	storageRam, err := storage.CreateRamStorage()
	require.NoError(t, err)
	ts := httptest.NewServer(Create(log, storageRam).Route())
	defer ts.Close()

	type reqJSON struct {
		method   string
		username string
		password string
		url      string
		body     string
	}
	type want struct {
		code int
		body string
	}

	type test struct {
		name string
		init []reqJSON
		get  reqJSON
		want want
	}
	tests := []test{
		{
			name: "common case",
			init: []reqJSON{
				{
					method:   http.MethodPost,
					username: "asd",
					password: "asd",
					url:      "/",
					body:     `{"data":"some_text"}`,
				},
				{
					method:   http.MethodPost,
					username: "asd",
					password: "asd",
					url:      "/",
					body:     `{"data":"some_text2"}`,
				},
				{
					method:   http.MethodPost,
					username: "asd",
					password: "asd",
					url:      "/",
					body:     `{"data":"some_text3"}`,
				},
			},
			get: reqJSON{
				method:   http.MethodGet,
				username: "asd",
				password: "asd",
				url:      "/",
			},
			want: want{
				code: http.StatusOK,
				body: "{\"dataArray\":[\"some_text\",\"some_text2\",\"some_text3\"]}",
			},
		},
		{
			name: "case another value",
			init: []reqJSON{
				{
					method:   http.MethodPost,
					username: "qwe",
					password: "qwe",
					url:      "/",
					body:     `{"data":"some_text"}`,
				},
				{
					method:   http.MethodPost,
					username: "qwe",
					password: "qwe",
					url:      "/",
					body:     `{"data":"some_text2"}`,
				},
			},
			get: reqJSON{
				method:   http.MethodGet,
				username: "qwe",
				password: "qwe",
				url:      "/",
			},
			want: want{
				code: http.StatusOK,
				body: "{\"dataArray\":[\"some_text\",\"some_text2\"]}",
			},
		},
		{
			name: "case error in url",
			get: reqJSON{
				method:   http.MethodGet,
				username: "qwe",
				password: "qwe",
				url:      "/csdv",
			},
			want: want{
				code: http.StatusNotFound,
				body: "404 page not found\n",
			},
		},
		{
			name: "case unauth",
			get: reqJSON{
				method:   http.MethodGet,
				username: "qwe",
				password: "qe",
				url:      "/",
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "Unauthorized\n",
			},
		},
	}

	postReq := func(data *test) {
		for _, reqData := range data.init {
			body := bytes.NewBufferString(reqData.body)
			req, errReq := http.NewRequest(reqData.method, ts.URL+reqData.url, body)
			require.NoError(t, errReq)

			req.Header.Add("Content-Type", "application/json")
			req.SetBasicAuth(reqData.username, reqData.password)
			resp, errResp := ts.Client().Do(req)
			assert.NoError(t, errResp)
			defer resp.Body.Close()

			_, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postReq(&tt)

			body := bytes.NewBufferString(tt.get.body)
			req, errReq := http.NewRequest(tt.get.method, ts.URL+tt.get.url, body)
			require.NoError(t, errReq)

			req.Header.Add("Content-Type", "application/json")
			req.SetBasicAuth(tt.get.username, tt.get.password)
			resp, errResp := ts.Client().Do(req)
			assert.NoError(t, errResp)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			//assert.Equal(t, tt.req.method, )
			assert.Equal(t, tt.want.body, string(respBody))
			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}
