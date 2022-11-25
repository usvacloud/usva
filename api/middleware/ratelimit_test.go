package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	fakerInstance = faker.New()
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	os.Exit(m.Run())
}

func TestRatelimiter_RestrictUploads(t *testing.T) {
	t.Parallel()
	type args struct {
		allowedData uint64
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
				Identifier: "test-identifier-valid",
				Uploads:    nil,
			},
			args: args{
				allowedData: 32,
			},
			want:       http.StatusOK,
			wantHeader: true,
		},
		{
			name: "invalid user",
			client: Client{
				Identifier:  "test-identifier-invalid",
				handler:     &RequestHandler{},
				LastRequest: time.Now().Add(-time.Minute),
				Uploads: &[](*ClientUpload){
					{
						Size: 32,
						Time: time.Now().Add(-time.Minute),
					},
				},
			},
			args: args{
				allowedData: 32,
			},
			want:       http.StatusNotAcceptable,
			wantHeader: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiterBase := Ratelimiter{
				Clients: &[](*Client){&tt.client},
			}

			res := httptest.NewRecorder()
			_, r := gin.CreateTestContext(res)

			limiter := limiterBase.RestrictUploads(time.Hour, tt.args.allowedData)
			r.POST("/", limiter, func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			req := httptest.NewRequest("POST", "/", strings.NewReader(fakerInstance.App().Name()))
			req.Header.Set("Api-Identifier", tt.client.Identifier)
			r.ServeHTTP(res, req)

			assert.Equal(t, tt.want, res.Result().StatusCode, tt.name)
			if tt.wantHeader {
				assert.NotEmpty(t, res.Header().Get("Usva-AllowedUploadBytes"), tt.name)
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
			name: "test case 1",
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
			name: "test case 1",
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
			hand := &RequestHandler{
				nextReset:     tt.fields.nextReset,
				tokens:        tt.fields.tokens,
				maximumTokens: tt.fields.maximumTokens,
			}
			if got := hand.UseToken(tt.args.count); got != tt.want {
				t.Errorf("RequestHandler.UseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRatelimiter_RestrictRequests(t *testing.T) {
	t.Parallel()

	type fields struct {
		lastCleanup      time.Time
		allowedRequests  int16
		testRequestCount uint
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "fail-test",
			fields: fields{
				lastCleanup:      time.Now().Add(-time.Minute),
				allowedRequests:  1,
				testRequestCount: 2,
			},
			want: http.StatusTooManyRequests,
		},
		{
			name: "success-test",
			fields: fields{
				lastCleanup:      time.Now().Add(-time.Minute),
				allowedRequests:  2,
				testRequestCount: 2,
			},
			want: http.StatusOK,
		},
		{
			name: "success-test-nolimits",
			fields: fields{
				lastCleanup:      time.Now().Add(-time.Minute),
				allowedRequests:  0,
				testRequestCount: 10,
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiterBase := NewRatelimiter()

			limiter := limiterBase.RestrictRequests(tt.fields.allowedRequests, time.Hour)

			responseWriter := httptest.NewRecorder()
			_, r := gin.CreateTestContext(responseWriter)
			r.GET("/", limiter)

			var res *httptest.ResponseRecorder
			for iter := int16(0); uint(iter) < tt.fields.testRequestCount; iter++ {
				req := httptest.NewRequest("GET", "/", strings.NewReader("hello"))
				req.Header.Set("Api-Identifier", "api-identifier")

				res = httptest.NewRecorder()
				r.Handler().ServeHTTP(res, req)
			}
			statusgot := res.Code
			assert.Equal(t, tt.want, statusgot)

		})
	}
}
