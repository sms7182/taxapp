package pkg

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service Service
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/health", cr.health)
	e.GET("/tax/fire_inquiry", cr.inquiry)
	e.GET("/failedNotify", cr.failedNotify)
	e.GET("/retryInvoice/:taxRawId", cr.retryInvoice)
	e.GET("/autoRetryInvoice", cr.autoRetry)
}

func (cr Controller) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) failedNotify(c *gin.Context) {
	go cr.Service.NotifyFailedTax()
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) inquiry(c *gin.Context) {
	go cr.Service.TaxRequestInquiry()
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) autoRetry(c *gin.Context) {
	go cr.Service.AutoRetry(c)
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) retryInvoice(c *gin.Context) {
	taxRawIdStr := c.Param("taxRawId")
	taxRawId, err := strconv.Atoi(taxRawIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if e := cr.Service.RetryInvoice(c, uint(taxRawId)); e != nil {
		c.JSON(http.StatusInternalServerError, e.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
