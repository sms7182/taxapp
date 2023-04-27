package pkg

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"tax-management/external"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service Service
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/health", cr.health)
	e.GET("/tax/fire_inquiry", cr.inquiry)

	e.GET("/autoRetryInvoice", cr.autoRetry)
	e.GET("/taxprocess/:id", cr.getTaxProcess)
	e.POST("/sendInvoice", cr.sendInvoice)
}

func (cr Controller) sendInvoice(c *gin.Context) {
	var rawTransaction external.RawTransaction
	request := c.Request
	reqBody, _ := ioutil.ReadAll(request.Body)
	request.Body.Close()
	_ = json.Unmarshal(reqBody, &rawTransaction)
	err := cr.Service.StartSendingInvoice(rawTransaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
func (cr Controller) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) getTaxProcess(c *gin.Context) {
	idStr := c.Param("id")
	taxProcessId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	tp, err := cr.Service.GetTaxProcess(context.Background(), uint(taxProcessId))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, tp)

}

func (cr Controller) inquiry(c *gin.Context) {
	go cr.Service.TaxRequestInquiry()
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) autoRetry(c *gin.Context) {
	go cr.Service.AutoRetry(c)
	c.JSON(http.StatusOK, gin.H{})
}
