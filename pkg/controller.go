package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service Service
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/health", cr.health)
	e.GET("/tax/fire_inquiry", cr.inquiry)
}

func (cr Controller) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) inquiry(c *gin.Context) {
	//cr.Service.TaxRequestInquiry()
}
