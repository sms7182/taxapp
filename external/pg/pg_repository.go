package pg

import "gorm.io/gorm"

type RepositoryImpl struct {
	DB *gorm.DB
}

func (repository RepositoryImpl) LogReqRes(requestTraceId string, signature string, packetType string, url string, statusCode int, req string, res *string, errorMsg *string) error {
	return nil
}
