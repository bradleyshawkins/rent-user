package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bradleyshawkins/berror"
	"github.com/lib/pq"
)

const (
	foreignKeyFailed string = "23503"
	duplicateEntry   string = "23505"
)

func toRentError(err error) error {
	if err == sql.ErrNoRows {
		return berror.Wrap(err, berror.WithNotExists(), berror.WithMessage("entity does not exist"))
	}

	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch string(pgErr.Code) {
		case duplicateEntry:
			return berror.Wrap(err, berror.WithDuplicate(), berror.WithMessage(fmt.Sprintf("duplicate entry found. Details: %s", pgErr.Detail)))
		case foreignKeyFailed:
			return berror.Wrap(err, berror.WithRequiredEntityNotExist(), berror.WithMessage(fmt.Sprintf("required entity does not exist. Details: %s", pgErr.Detail)))
		}

	}
	return berror.Wrap(err, berror.WithInternal(), berror.WithMessage("unexpected error occurred"))
}
