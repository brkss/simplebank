package api

import (
	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serve HTTP request for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creaet new HTTP server and setup routes
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts/:limit/:offset", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start new HTTP request and listen for requests !
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
