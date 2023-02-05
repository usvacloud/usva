package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// rH
type requestHandler struct {
	nextReset     time.Time
	tokens        int16
	maximumTokens int16
}

func newHandler(requestCount int16, saveDuration time.Duration) *requestHandler {
	return &requestHandler{
		nextReset:     time.Now().Add(saveDuration),
		tokens:        requestCount,
		maximumTokens: requestCount,
	}
}

func (h *requestHandler) useToken(count int16) bool {
	ok := h.tokens >= count
	if ok {
		h.tokens -= count
	}
	return ok
}

func (h *requestHandler) reset() {
	h.tokens = h.maximumTokens
}

// cU
type clientUpload struct {
	size      int64
	timestamp time.Time
}

// C
type Client struct {
	Identifier  string
	handler     *requestHandler
	lastRequest time.Time
	uploads     []clientUpload
}

// RateLimiter struct
type Ratelimiter struct {
	clients     []Client
	lastCleanup time.Time
}

func NewRatelimiter() *Ratelimiter {
	return &Ratelimiter{
		clients:     []Client{},
		lastCleanup: time.Now(),
	}
}

func (r *Ratelimiter) getExistingClient(identifier string) *Client {
	for _, client := range r.clients {
		if client.Identifier == identifier {
			return &client
		}
	}
	return nil
}

func (r *Ratelimiter) newClient(identifier string, handler *requestHandler) *Client {
	client := r.getExistingClient(identifier)
	if client == nil {
		client = &Client{
			Identifier:  identifier,
			handler:     handler,
			lastRequest: time.Now(),
		}
		r.clients = append(r.clients, *client)
	} else if client.handler == nil && handler != nil {
		client.handler = handler
	}
	return client
}

func (r *Ratelimiter) NewUpload(identifier string, upload clientUpload) {
	client := r.newClient(identifier, nil)
	client.uploads = append(client.uploads, upload)
}

// Remove clients that have full ratelimit capacity.
// TODO: This can also take quite a bit of memory as a new array is created and appended.
// fix is possible via removing the clients straight from the Clients struct
func (r *Ratelimiter) Clean() {
	nl := []Client{}
	for _, client := range r.clients {
		if client.handler.maximumTokens > client.handler.tokens {
			nl = append(nl, client)
		}
	}

	r.clients = nl
	r.lastCleanup = time.Now()
}

// RestrictRequests returns a middleware to create a new ratelimiter for each IP.
// This can take a lot of memory with higher use, though.
// TODO: Optimize for larger scale
func (r *Ratelimiter) RestrictRequests(count int16, per time.Duration) gin.HandlerFunc {
	if count == 0 {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}
	return func(ctx *gin.Context) {
		identifier := ctx.Request.Header.Get(Headers.Identifier)

		client := r.newClient(identifier, newHandler(count, per))

		if client.handler.nextReset.Before(time.Now()) {
			time.AfterFunc(per, client.handler.reset)
			client.handler.nextReset = time.Now().Add(per)
		}

		setResponseHeaders(ctx, count, client.handler.tokens, int16(per.Seconds()))
		if client.handler.useToken(1) {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}

		client.lastRequest = time.Now()
	}
}

// RestrictUploads checks the history of a client and
// limits their access based on found data.
// Allows a certain amount of data in specific duration.
func (r *Ratelimiter) RestrictUploads(
	duration time.Duration,
	allowedData uint64,
) gin.HandlerFunc {
	if allowedData == 0 {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}
	return func(ctx *gin.Context) {
		identifier := ctx.Request.Header.Get(Headers.Identifier)
		client := r.newClient(identifier, nil)
		client.lastRequest = time.Now()

		totalUploaded := uint64(ctx.Request.ContentLength)
		for _, upload := range client.uploads {
			if time.Since(upload.timestamp) > duration {
				continue
			}
			totalUploaded += uint64(upload.size)
		}

		if totalUploaded > allowedData {
			ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"error": "you've exceeded your upload capacity",
			})
			return
		}

		ctx.Header(Headers.AllowedBytes, fmt.Sprint(allowedData-totalUploaded))
	}
}
