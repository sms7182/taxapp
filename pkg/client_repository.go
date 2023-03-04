package pkg

type ClientRepository interface {
	LogReqRes(requestTraceId string, signature string, packetType string, url string, statusCode int, req string, res *string, errorMsg *string) error
}
