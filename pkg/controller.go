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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "fa")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		//	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (cr Controller) SetRoutes(e *gin.Engine) {

	e.Use(CORSMiddleware())
	e.GET("/health", cr.health)
	e.GET("/tax/fire_inquiry/:userName", cr.inquiry)

	e.GET("/autoRetryInvoice", cr.autoRetry)
	e.GET("/taxprocess/:id", cr.getTaxProcess)
	e.POST("/init_customer", cr.initCustomer)
	e.POST("/sendInvoice", cr.sendInvoice)
	e.POST("/sendInvoices", cr.sendInvoices)

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

func (cr Controller) sendInvoices(c *gin.Context) {
	var rawTransaction external.RawTransaction

	request := c.Request
	reqBody, _ := ioutil.ReadAll(request.Body)
	request.Body.Close()
	_ = json.Unmarshal(reqBody, &rawTransaction)
	id, err := cr.Service.StartSendingInvoice(context.Background(), rawTransaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, id)
}

func (cr Controller) sendInvoice(c *gin.Context) {
	var rawTransaction external.RawTransaction
	request := c.Request
	reqBody, _ := ioutil.ReadAll(request.Body)
	request.Body.Close()
	_ = json.Unmarshal(reqBody, &rawTransaction)
	id, err := cr.Service.StartSendingInvoice(context.Background(), rawTransaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, id)
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
