package pg

import "gorm.io/gorm"

type RepositoryImpl struct {
	DB *gorm.DB
}
