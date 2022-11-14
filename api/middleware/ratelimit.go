package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ClientUpload struct {
	Size int64
	Time time.Time
}
type ClientUploads []*ClientUpload

type Client struct {
	Identifier  string
	Limiter     *rate.Limiter
	LastRequest time.Time
	Uploads     *ClientUploads
}

type Ratelimiter struct {
	Clients   *[](*Client)
	LastReset time.Time
}

func safeListAccess[T *Client | *ClientUpload, L *[]T](f L) L {
	if f == nil {
		return &[]T{}
	}
	return f
}

func NewRatelimiter() *Ratelimiter {
	return &Ratelimiter{
		Clients:   &[](*Client){},
		LastReset: time.Now(),
	}
}

func (limiterBase *Ratelimiter) addClient(client *Client) {
	newValue := append(*safeListAccess(limiterBase.Clients), client)
	limiterBase.Clients = &newValue
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
			Uploads:     &ClientUploads{},
			Limiter:     &rate.Limiter{},
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

// Remove clients that have not made any requests since their limits expired
// TODO: This can also take a bit of memory as a new array is created and appended.
// fix is possible via removing the clients straight from the Clients struct
func (limiterBase *Ratelimiter) Clean() {
	clients := safeListAccess(limiterBase.Clients)
	activeClients := &[](*Client){}
	for _, client := range *clients {
		if time.Since(client.LastRequest) > client.Limiter.Reserve().Delay() {
			*activeClients = append(*activeClients, client)
		}
	}

	limiterBase.Clients = activeClients
	limiterBase.LastReset = time.Now()
}

// Limit returns a middleware to create a new ratelimiter for each IP.
// This can take a lot of memory with higher use, though.
// TODO: Optimize for larger scale
func (limiterBase *Ratelimiter) Limit(count int, per time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawip := []byte(ctx.ClientIP())
		ip := hex.EncodeToString(sha256.New().Sum(rawip))

		found, client := limiterBase.getClientByIdentifierOrCreate(ip)
		client.LastRequest = time.Now()
		if !found {
			client.Limiter = rate.NewLimiter(rate.Every(per), count)
			return
		}

		if client.Limiter.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}

// Allow a certain amount of data in specific duration
// LimitDependsBodySize checks the history of a client and
// limits their access based on found data
func (limiterBase *Ratelimiter) LimitDependsBodySize(
	duration time.Duration,
	allowedData int64,
) gin.HandlerFunc {
	if limiterBase == nil {
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

		sum := int64(0)
		for _, upload := range *client.Uploads {
			if time.Since(upload.Time) > duration {
				break
			}
			sum += upload.Size
		}

		if sum > allowedData {
			ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"error": "you've exceeded your upload capacity",
			})
			return
		}

	}
}
