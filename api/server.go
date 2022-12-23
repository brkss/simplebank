package api

import (
	"github.com/brkss/simplebank/token"
	"github.com/brkss/simplebank/utils"

	"fmt"

	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serve HTTP request for our banking service
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}

// NewServer creaet new HTTP server and setup routes
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.SetupRouter()
	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/login", server.LoginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/account/:id", server.getAccount)
	authRoutes.GET("/accounts/:limit/:offset", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start new HTTP request and listen for requests !
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
