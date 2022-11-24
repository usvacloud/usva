package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestRatelimiter_RestrictUploads(t *testing.T) {
	fakerInstance := faker.New()

	type args struct {
		statusWants int
		allowedData uint64
		wantsHeader bool
	}
	tests := []struct {
		name   string
		client Client
		args   args
		want   gin.HandlerFunc
	}{
		{
			name: "valid user",
			client: Client{
				Identifier: "test-identifier-valid",
				Uploads:    nil,
			},
			args: args{
				allowedData: 32,
				statusWants: 200,
				wantsHeader: true,
			},
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
				statusWants: 406,
				wantsHeader: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiterBase := Ratelimiter{
				Clients: &[](*Client){&tt.client},
			}

			gin.SetMode(gin.ReleaseMode)
			res := httptest.NewRecorder()
			_, r := gin.CreateTestContext(res)

			limiter := limiterBase.RestrictUploads(time.Hour, tt.args.allowedData)
			r.POST("/", limiter, func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			req := httptest.NewRequest("POST", "/", strings.NewReader(fakerInstance.App().Name()))
			req.Header.Set("Api-Identifier", tt.client.Identifier)
			r.ServeHTTP(res, req)

			assert.Equal(t, tt.args.statusWants, res.Result().StatusCode, tt.name)
			if tt.args.wantsHeader {
				assert.NotEmpty(t, res.Header().Get("Usva-AllowedUploadBytes"), tt.name)
			}
		})
	}
}
