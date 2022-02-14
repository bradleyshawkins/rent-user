package identity

import (
	"net/mail"

	"github.com/bradleyshawkins/berror"

	"golang.org/x/crypto/bcrypt"
)

// signUpSteps contains all methods used to sign up
type signUpSteps interface {
	RegisterUser(u *User, c *Credentials) error
	RegisterAccount(a *Account) error
	AddUserToAccount(accountID AccountID, u UserID, role Role) error
}

// SignUpForm contains all necessary data needed to sign up
type SignUpForm struct {
	credentials *Credentials
	user        *User
	account     *Account
}

// SignUp applies steps necessary to register with the service
func (s *SignUpForm) SignUp(sus signUpSteps) error {
	err := sus.RegisterUser(s.user, s.credentials)
	if err != nil {
		return err
	}

	err = sus.RegisterAccount(s.account)
	if err != nil {
		return err
	}

	err = sus.AddUserToAccount(s.account.ID, s.user.ID, RoleOwner)
	if err != nil {
		return err
	}
	return nil
}

// signUpper is the interface that begins the signup process.
type signUpper interface {
	SignUp(s *SignUpForm) error
}

// SignUpManager handles initiating the signup process
type SignUpManager struct {
	su signUpper
}

// NewSignUpManager is a constructor for SignUpManager
func NewSignUpManager(uc signUpper) *SignUpManager {
	return &SignUpManager{su: uc}
}

// SignUp creates the types needed for signing off and kicks off the signing up steps
func (u *SignUpManager) SignUp(username string, password string, emailAddress *mail.Address, firstName, lastName string) (*User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return nil, berror.WrapInternal(err, "unable to encrypt password")
	}

	suf := &SignUpForm{
		credentials: &Credentials{
			Username: username,
			Password: string(pw),
		},
		user: &User{
			ID:           NewUserID(),
			EmailAddress: emailAddress,
			FirstName:    firstName,
			LastName:     lastName,
			Status:       UserActive,
		},
		account: &Account{
			ID:     NewAccountID(),
			Status: AccountActive,
		},
	}

	err = u.su.SignUp(suf)
	if err != nil {
		return nil, err
	}

	return suf.user, nil
}
