package postgres

import (
	"github.com/bradleyshawkins/rent-user/identity"
)

func (t *transaction) RegisterAccount(a *identity.Account) error {
	_, err := t.tx.Exec(`INSERT INTO account(id, status) VALUES ($1, $2)`, a.ID.AsUUID(), a.Status)
	if err != nil {
		return toRentError(err)
	}

	return nil
}

func (t *transaction) AddUserToAccount(aID identity.AccountID, uID identity.UserID, role identity.Role) error {
	_, err := t.tx.Exec(`INSERT INTO membership(account_id, app_user_id, role) VALUES ($1, $2, $3)`, aID.AsUUID(), uID.AsUUID(), role)
	if err != nil {
		return toRentError(err)
	}
	return nil
}
