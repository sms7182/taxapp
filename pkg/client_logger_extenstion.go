package pkg

import "net/http"

type ClientLoggerExtension interface {
	Do(taxRawId uint, taxProcessId uint, request *http.Request, apiname string) (*http.Response, error)
}
