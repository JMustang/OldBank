package api

import (
	db "github.com/JMustang/OldBank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking servoce.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Add router
	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
