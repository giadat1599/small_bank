package api

import (
	db "github.com/giadat1599/small_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// StartServer runs the HTTP server on a specific address
func (server *Server) StartServer(addr string) error {
	return server.router.Run(addr)
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
