package pkg

import (
	"context"
	"encoding/json"
	"fmt"
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
	e.POST("/init_customer", cr.initCustomer)
	e.POST("/sendInvoice", cr.sendInvoice)
}

func (cr Controller) initCustomer(c *gin.Context) {
	var cdto external.CustomerDto
	request := c.Request
	reqBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	_ = json.Unmarshal(reqBody, &cdto)

	if cdto.PrivateKey == "" || cdto.PublicKey == "" || cdto.UserName == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("privateKey or public or userName is nil"))
		return
	}
	customerId, err := cr.Service.InitialCustomer(&cdto)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, customerId)

}
func (cr Controller) sendInvoice(c *gin.Context) {
	var rawTransaction external.RawTransaction
	request := c.Request
	reqBody, _ := ioutil.ReadAll(request.Body)
	request.Body.Close()
	_ = json.Unmarshal(reqBody, &rawTransaction)
	err := cr.Service.StartSendingInvoice(context.Background(), rawTransaction)
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
	userName := c.Param("userName")

	go cr.Service.TaxRequestInquiry(userName)
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) autoRetry(c *gin.Context) {
	go cr.Service.AutoRetry(c)
	c.JSON(http.StatusOK, gin.H{})
}
