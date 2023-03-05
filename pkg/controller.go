package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service      Service
	KafkaService KafkaService
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/health", cr.health)
}

func (cr Controller) health(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{})
}
