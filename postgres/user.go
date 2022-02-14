package postgres

import (
	"github.com/bradleyshawkins/rent-user/identity"
	"github.com/google/uuid"
)

// SignUp provides a transaction around the sign up process
func (d *Database) SignUp(suf *identity.SignUpForm) error {
	tx, err := d.begin()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.rollback()
	}()

	err = suf.SignUp(tx)
	if err != nil {
		return err
	}

	err = tx.commit()
	if err != nil {
		return err
	}
	return nil
}

func (t *transaction) RegisterUser(user *identity.User, c *identity.Credentials) error {
	detailsID := uuid.New()
	_, err := t.tx.Exec(`INSERT INTO app_user_details(id, first_name, last_name, email_address) VALUES ($1, $2, $3, $4)`, detailsID, user.FirstName, user.LastName, user.EmailAddress.Address)
	if err != nil {
		return toRentError(err)
	}

	credentialsID := uuid.New()
	_, err = t.tx.Exec(`INSERT INTO app_user_credentials(id, username, password) VALUES ($1, $2, $3)`, credentialsID, c.Username, c.Password)
	if err != nil {
		return toRentError(err)
	}

	_, err = t.tx.Exec("INSERT INTO app_user(id, status, app_user_credentials_id, app_user_details_id) VALUES ($1, $2, $3, $4)", user.ID.AsUUID(), user.Status, credentialsID, detailsID)
	if err != nil {
		return toRentError(err)
	}

	return nil
}
