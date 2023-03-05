package pg

import "tax-management/pkg"

type ServiceImp struct {
	Repository pkg.ClientRepository
	Client     pkg.ClientLoggerExtension
}
