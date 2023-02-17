package ratelimit

import (
	"bytes"
	"crypto/rand"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	if code := m.Run(); code != 0 {
		os.Exit(code)
	}
}

func TestRatelimiter_RestrictUploads(t *testing.T) {
	t.Parallel()
	type args struct {
		allowedData          uint64
		uploadSize           uint64
		resetInterval        time.Duration
		durationBeforeUpload time.Duration
	}
	tests := []struct {
		name       string
		client     Client
		args       args
		want       int
		wantHeader bool
	}{
		{
			name: "valid user",
			client: Client{
				Identifier: "test-valid",
				uploads:    nil,
			},
			args: args{
				allowedData:          32,
				uploadSize:           32,
				resetInterval:        time.Second * 5,
				durationBeforeUpload: 0,
			},
			want:       http.StatusOK,
			wantHeader: true,
		},
		{
			name: "invalid user",
			client: Client{
				Identifier:  "test-invalid",
				lastRequest: time.Now().Add(-time.Second),
				uploads: []ClientUpload{
					{
						size:      32,
						timestamp: time.Now().Add(-time.Second),
					},
				},
			},
			args: args{
				allowedData:          64,
				uploadSize:           33,
				resetInterval:        time.Second * 5,
				durationBeforeUpload: 0,
			},
			want:       http.StatusNotAcceptable,
			wantHeader: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiterBase := Ratelimiter{
				clients: []Client{tt.client},
			}

			res := httptest.NewRecorder()
			_, r := gin.CreateTestContext(res)

			limiter := limiterBase.RestrictUploads(tt.args.resetInterval, tt.args.allowedData)
			r.POST("/", limiter, func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			}) // initialize handler

			buf := make([]byte, tt.args.uploadSize)
			_, err := io.ReadFull(rand.Reader, buf)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest("POST", "/", bytes.NewReader(buf))
			req.Header.Set(Headers.Identifier, tt.client.Identifier)

			time.Sleep(tt.args.durationBeforeUpload)
			r.ServeHTTP(res, req)

			if tt.want != res.Result().StatusCode {
				t.Fatalf("wanted status code %d got %d", tt.want, res.Result().StatusCode)
			}
			if tt.wantHeader && res.Header().Get(Headers.AllowedBytes) == "" {
				t.Fatalf("%s header is empty", Headers.AllowedBytes)
			}
		})
	}
}

func TestRequestHandler_UseToken(t *testing.T) {
	t.Parallel()
	type fields struct {
		nextReset     time.Time
		tokens        int16
		maximumTokens int16
	}
	type args struct {
		count int16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test",
			fields: fields{
				nextReset:     time.Now().Add(time.Hour),
				tokens:        3,
				maximumTokens: 3,
			},
			args: args{
				count: 1,
			},
			want: true,
		},
		{
			name: "test 1",
			fields: fields{
				nextReset:     time.Now().Add(time.Hour),
				tokens:        0,
				maximumTokens: 3,
			},
			args: args{
				count: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := &tokenStorage{
				nextReset:     tt.fields.nextReset,
				tokens:        tt.fields.tokens,
				maximumTokens: tt.fields.maximumTokens,
			}
			if got := hand.useToken(tt.args.count); got != tt.want {
				t.Errorf("RequestHandler.UseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRatelimiter_RestrictRequests(t *testing.T) {
	t.Parallel()

	type args struct {
		lastCleanup      time.Time
		allowedRequests  int16
		testRequestCount uint
		beforeRequest    time.Duration
		cleanInterval    time.Duration
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "fail",
			args: args{
				lastCleanup:      time.Now().Add(-time.Minute),
				allowedRequests:  1,
				testRequestCount: 2,
				beforeRequest:    0,
				cleanInterval:    time.Second,
			},
			want: http.StatusTooManyRequests,
		},
		{
			name: "success",
			args: args{
				lastCleanup:      time.Now().Add(-time.Minute),
				allowedRequests:  2,
				testRequestCount: 2,
				beforeRequest:    0,
				cleanInterval:    time.Second,
			},
			want: http.StatusOK,
		},
		{
			name: "success-1",
			args: args{
				lastCleanup:      time.Now().Add(-time.Minute),
				allowedRequests:  0,
				testRequestCount: 10,
				beforeRequest:    0,
				cleanInterval:    time.Second,
			},
			want: http.StatusOK,
		},
		{
			name: "success-2",
			args: args{
				lastCleanup:      time.Now().Add(-time.Minute),
				allowedRequests:  3,
				testRequestCount: 4,
				beforeRequest:    time.Second / 2,
				cleanInterval:    time.Second,
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiterBase := NewRatelimiter()
			limiter := limiterBase.RestrictRequests(tt.args.allowedRequests, tt.args.cleanInterval)

			responseWriter := httptest.NewRecorder()
			_, r := gin.CreateTestContext(responseWriter)
			r.GET("/", limiter)

			var res *httptest.ResponseRecorder
			for iter := uint(0); iter < tt.args.testRequestCount; iter++ {
				req := httptest.NewRequest("GET", "/", strings.NewReader("hello"))
				req.Header.Set(Headers.Identifier, "api-identifier")

				time.Sleep(tt.args.beforeRequest)
				res = httptest.NewRecorder()
				r.Handler().ServeHTTP(res, req)
			}

			if tt.want != res.Code {
				t.Fatalf("want %d got %d", tt.want, res.Code)
			}
		})
	}
}
