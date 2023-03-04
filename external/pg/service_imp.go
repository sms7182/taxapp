package pg

import "tax-app/pkg"

type ServiceImp struct {
	Repository pkg.ClientRepository
	Client     pkg.ClientLoggerExtension
}
