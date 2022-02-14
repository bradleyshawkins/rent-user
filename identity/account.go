package identity

import (
	"fmt"

	"github.com/bradleyshawkins/berror"

	"github.com/google/uuid"
)

type AccountID uuid.UUID

func NewAccountID() AccountID {
	return AsAccountID(uuid.New())
}

func AsAccountID(id uuid.UUID) AccountID {
	return AccountID(id)
}

func (a AccountID) AsUUID() uuid.UUID {
	return uuid.UUID(a)
}

func (a AccountID) IsZero() bool {
	return a.AsUUID() == uuid.Nil
}

func (a AccountID) String() string {
	return a.AsUUID().String()
}

type AccountStatus string

const (
	AccountDisabled AccountStatus = "Disabled"
	AccountActive   AccountStatus = "Active"
	AccountCanceled AccountStatus = "Canceled"
)

type Role string

const (
	RoleOwner  Role = "Owner"
	RoleWriter Role = "Writer"
	RoleReader Role = "Reader"
)

var roleMap = map[string]Role{
	"Owner":  RoleOwner,
	"Writer": RoleWriter,
	"Reader": RoleReader,
}

func ToRole(role string) (Role, error) {
	r, ok := roleMap[role]
	if !ok {
		return "", berror.New("invalid role", berror.WithMessage(fmt.Sprintf("role %s is not a valid role", role)))
	}
	return r, nil
}

type Account struct {
	ID      AccountID
	Status  AccountStatus
	Members map[UserID]Role
}

type AccountMemberships struct {
	AccountToRole map[AccountID]Role
}
