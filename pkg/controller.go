package pkg

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/health", cr.health)
}

func (cr Controller) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
