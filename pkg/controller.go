package pkg

import (
	"encoding/json"
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
	e.GET("/decrypt", cr.decryption)
}

func (cr Controller) health(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) decryption(c *gin.Context) {
	key := c.Query("key")

	encrypted := c.Query("enc")
	result := utility.Decrypt(encrypted, key)
	var obj utility.DataToEncrypt
	err := json.Unmarshal(result, &obj)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, obj)
}

func (cr Controller) encryption(c *gin.Context) {
	result := utility.EncryptWithIV()
	c.JSON(http.StatusOK, result)
}
