package pkg

import "net/http"

type ClientLoggerExtension interface {
	Do(requestTranceId string, signature string, packetType string, request *http.Request, gateway string) (*http.Response, error)
}
