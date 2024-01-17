package api

import (
	db "booking-api/db/sqlc"
	"booking-api/token"
	"booking-api/util"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server struct serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      *db.Store
	router     *gin.Engine
}

func healthy(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Healthy using Webhooks...")
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store *db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	gin.SetMode(config.GinMode)

	router := gin.Default()

	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{"*"}
	// corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	// corsConfig.AllowHeaders = []string{"Origin"}
	// corsConfig.AllowAllOrigins = true
	// corsConfig.AllowFiles	= true
	// corsConfig.AddAllowHeaders("Authorization")
	// corsConfig.AddAllowHeaders("content-type")
	// corsConfig.AddAllowHeaders("Access-Control-Allow-Origin", "*")
	// config.AllowOrigins = []string{"http://google.com", "http://facebook.com"}
	// corsConfig.AllowAllOrigins = true

	// router.Use(cors.New(corsConfig))

	corsConfig := cors.Default()
	router.Use(corsConfig)
	// router.Use(favicon.New("./favicon.ico"))
	router.Use(server.rateLimit())

	// do not trust all proxies
	// router.SetTrustedProxies([]string{"192.168.1.2"})
	router.SetTrustedProxies(nil)
	router.TrustedPlatform = gin.PlatformCloudflare

	Routes(router, server)

	server.router = router

	return server, nil
}

// StartServer runs the HTTP server on a specific address
func (srv *Server) StartServer(address string) error {
	fmt.Printf("Server starting on address: %s\n", address)
	return srv.router.Run(fmt.Sprintf("0.0.0.0:%s", address))
}
