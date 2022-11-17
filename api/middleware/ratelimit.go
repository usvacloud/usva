package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	nextReset     time.Time
	tokens        int16
	maximumTokens int16
}

func NewHandler(requestCount int16, saveDuration time.Duration) *RequestHandler {
	return &RequestHandler{
		nextReset:     time.Now().Add(saveDuration),
		tokens:        requestCount,
		maximumTokens: requestCount,
	}
}

func (hand *RequestHandler) UseToken(count int16) bool {
	ok := hand.tokens >= count
	if ok {
		hand.tokens -= count
	}
	return ok
}
func (hand *RequestHandler) ResetTokens() {
	hand.tokens = hand.maximumTokens
}

type ClientUpload struct {
	Size int64
	Time time.Time
}
type Client struct {
	Identifier  string
	handler     *RequestHandler
	LastRequest time.Time
	Uploads     *[]*ClientUpload
}

type Ratelimiter struct {
	Clients     *[](*Client)
	LastCleanup time.Time
}

func safeListAccess[T *Client | *ClientUpload, L *[]T](f L) L {
	if f == nil {
		return &[]T{}
	}
	return f
}

func NewRatelimiter() *Ratelimiter {
	return &Ratelimiter{
		Clients:     &[](*Client){},
		LastCleanup: time.Now(),
	}
}

func (limiterBase *Ratelimiter) addClient(client *Client) {
	newValue := append(*safeListAccess(limiterBase.Clients), client)
	limiterBase.Clients = &newValue
}

func setResponseHeaders(ctx *gin.Context, limit, remaining, toreset int16) {
	ctx.Header("RateLimit-Limit", fmt.Sprint(limit))
	ctx.Header("RateLimit-Remaining", fmt.Sprint(remaining))
	ctx.Header("RateLimit-Reset", fmt.Sprint(toreset))
}

func (limiterBase *Ratelimiter) getClientByIdentifier(identifier string) (safe bool, client *Client) {
	for _, client := range *safeListAccess(limiterBase.Clients) {
		if client.Identifier == identifier {
			return true, client
		}
	}
	return false, nil
}

func (limiterBase *Ratelimiter) getClientByIdentifierOrCreate(identifier string) (found bool, client *Client) {
	found, client = limiterBase.getClientByIdentifier(identifier)
	if !found {
		client = &Client{
			Identifier:  identifier,
			Uploads:     &[]*ClientUpload{},
			handler:     &RequestHandler{},
			LastRequest: time.Now(),
		}
		limiterBase.addClient(client)
	}
	return found, client
}

func (limiterBase *Ratelimiter) AppendClientUploads(identifier string, upload ClientUpload) {
	_, client := limiterBase.getClientByIdentifierOrCreate(identifier)
	*client.Uploads = append(*client.Uploads, &upload)
}

// Remove clients that have full ratelimit capacity.
// TODO: This can also take quite a bit of memory as a new array is created and appended.
// fix is possible via removing the clients straight from the Clients struct
func (limiterBase *Ratelimiter) Clean() {
	clients := safeListAccess(limiterBase.Clients)
	nl := []*Client{}
	for _, client := range *clients {
		if client.handler.maximumTokens > client.handler.tokens {
			nl = append(nl, client)
		}
	}

	*clients = nl
	limiterBase.LastCleanup = time.Now()
}

// RestrictRequests returns a middleware to create a new ratelimiter for each IP.
// This can take a lot of memory with higher use, though.
// TODO: Optimize for larger scale
func (limiterBase *Ratelimiter) RestrictRequests(count int16, per time.Duration) gin.HandlerFunc {
	if count == 0 {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}
	return func(ctx *gin.Context) {
		rawip := []byte(ctx.ClientIP())
		ip := hex.EncodeToString(sha256.New().Sum(rawip))

		found, client := limiterBase.getClientByIdentifierOrCreate(ip)
		if !found {
			client.handler = NewHandler(count, per)
		}

		if client.handler.nextReset.Before(time.Now()) {
			time.AfterFunc(per, client.handler.ResetTokens)
			client.handler.nextReset = time.Now().Add(per)
		}

		setResponseHeaders(ctx, count, client.handler.tokens, int16(per.Seconds()))
		if client.handler.UseToken(1) {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
		client.LastRequest = time.Now()
	}
}

// RestrictUploads checks the history of a client and
// limits their access based on found data.
// Allows a certain amount of data in specific duration.
func (limiterBase *Ratelimiter) RestrictUploads(
	duration time.Duration,
	allowedData uint64,
) gin.HandlerFunc {
	if allowedData == 0 {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	}
	return func(ctx *gin.Context) {
		if allowedData == 0 {
			ctx.Next()
			return
		}

		rawip := []byte(ctx.ClientIP())
		ip := hex.EncodeToString(sha256.New().Sum(rawip))

		_, client := limiterBase.getClientByIdentifierOrCreate(ip)
		client.LastRequest = time.Now()

		if client.Uploads == nil {
			ctx.Next()
			return
		}

		sum := uint64(0)
		for _, upload := range *client.Uploads {
			if time.Since(upload.Time) > duration {
				break
			}
			sum += uint64(upload.Size)
		}

		if sum > allowedData {
			ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"error": "you've exceeded your upload capacity",
			})
			return
		}

	}
}
