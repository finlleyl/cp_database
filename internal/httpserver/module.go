package httpserver

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger) *gin.Engine {
	r := gin.New()

	r.Use(
		gin.Recovery(),
		gin.Logger(),
		GinLogger(logger),
		RequestID(),
		CORSMiddleware(),
	)

	return r
}
