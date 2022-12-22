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

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts/:limit/:offset", server.listAccounts)

	router.POST("/users", server.createUser)
	router.POST("/login", server.LoginUser)

	router.POST("/transfers", server.createTransfer)
	/*
		router.GET("/auth", authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
			return
		})
	*/

	server.router = router
}

// Start new HTTP request and listen for requests !
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
