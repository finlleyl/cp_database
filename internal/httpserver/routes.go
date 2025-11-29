package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func pingRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}