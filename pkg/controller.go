package pkg

import (
	"net/http"
	"tax-management/utility"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service      Service
	KafkaService KafkaService
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/health", cr.health)
	e.GET("/encrypt", cr.encryption)
}

func (cr Controller) health(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) encryption(c *gin.Context) {
	result := utility.Encrypt()
	c.JSON(http.StatusOK, result)
}
