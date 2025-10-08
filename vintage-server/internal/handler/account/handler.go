package account

import "vintage-server/internal/domain/account"

type handler struct {
	svc account.AccountService
}
func NewAccountHandler(svc account.AccountService) account.AccountHandler {
	return &handler {
		svc: svc,
	}
}
