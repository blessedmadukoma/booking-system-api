package api

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// rateLimit - IP-based rate limiting
func (srv *Server) rateLimit() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// background goroutine to remove old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				// check if the client hasn't been seen for the past 3 minutes
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(ctx *gin.Context) {
		if srv.config.Limiter.ENABLED {
			ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)

			if err != nil {
				log.Fatal("error splitting network address:", err)
				return
			}

			// lock the mutex to prevent concurrent execution
			mu.Lock()

			// check if the IP exists in the map, if it doesn't, initialize a new rate limiter and add the IP address and limiter to the map
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(srv.config.Limiter.RPS), srv.config.Limiter.BURST),
				}
			}

			// update the client's last seen
			clients[ip].lastSeen = time.Now()

			// if the request is not allowed, unlock the mutex and send 429 error
			if !clients[ip].limiter.Allow() {
				// fmt.Println("IP:", ip, "\nLast seen:", clients[ip].lastSeen.String(), "\nTokens:", clients[ip].limiter.Tokens(), "\n...")
				mu.Unlock()

				srv.rateLimitExceededResponse(ctx)
				return
			}
			// Very Important: unlock the mutex before calling the next handler in the chain.
			mu.Unlock()
		}

		ctx.Next()
	}
}
