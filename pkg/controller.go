package pkg

import (
	"encoding/json"
	"net/http"
	"tax-management/external/exkafka/messages"
	"tax-management/utility"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service      Service
	KafkaService KafkaService
	TaxClient    TaxClient
}

func (cr Controller) SetRoutes(e *gin.Engine) {
	e.GET("/getToken", cr.getToken)
	e.GET("/health", cr.health)
	e.GET("/encrypt", cr.encryption)
	e.GET("/decrypt", cr.decryption)
}
func (cr Controller) getToken(c *gin.Context) {
	cr.TaxClient.GetToken()
	c.JSON(http.StatusOK, gin.H{})
}

func (cr Controller) health(c *gin.Context) {

	cr.KafkaService.Publish(messages.RawTransaction{
		After: messages.AfterData{
			Trn:     "13981512",
			Taxid:   "A118GE",
			Indatim: 1677974399000000,
			Inty:    2,
			Inp:     1,
			Ins:     1,
			Tins:    "14004958663",
			Tinb:    "1",
			Tob:     2,
			Tprdis:  10560000,
			Tdis:    0,
			Tadis:   10560000,
			Tvam:    0,
			Todam:   0,
			Tbill:   10560000,
			Setm:    1,
			Sstid:   "1",
			Am:      2,
			Fee:     5280000,
			Prdis:   10560000,
			Dis:     0,
			Adis:    10560000,
			Vra:     0,
			Vam:     0,
			Tsstam:  10560000,
		},
		Source: messages.SourceData{
			Version:   "1.9.5.Final",
			Connector: "postgresql",
			Name:      "dbz.warehouse",
			TsMs:      1678112970619,
			Snapshot:  "true",
			DB:        "data",
			Sequence:  "[null,'53992413464408']",
			Schema:    "public",
			Table:     "tax_dom_hotel_retail",
			TxID:      539331996,
			Lsn:       53992413464408,
			Xmin:      nil,
		},
		Op:          "r",
		TsMs:        1678112970619,
		Transaction: nil,
	})
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
