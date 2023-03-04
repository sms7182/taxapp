package taxorganization

import (
	"tax-app/external/gateway"
)

type ClientImpl struct {
	HttpClient gateway.ClientLoggerExtension
}
