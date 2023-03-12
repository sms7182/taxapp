package pkg

import "net/http"

type ClientLoggerExtension interface {
	Do(taxRawId *uint, taxProcessId *uint, requestUniqueId string, request *http.Request, apiname string) (*http.Response, error)
}
