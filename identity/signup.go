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

//
//// register is a closure that registers the user and an account
//func (u *SignUpManager) register(user *SignUpForm, account *AccountRegistration) RegistrationFunc {
//	return func(us signUpSteps) error {
//		err := us.RegisterUser(user)
//		if err != nil {
//			return err
//		}
//
//		err = us.RegisterAccount(user.ID, account)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
//}
//
//func (u *SignUpManager) RegisterUserToAccount(accountID AccountID, role string, emailAddress string, firstName string, lastName string, password string) (*SignUpForm, error) {
//	addr, err := mail.ParseAddress(emailAddress)
//	if err != nil {
//		return nil, err
//	}
//
//	r, ok := roleMap[role]
//	if !ok {
//		return nil, rent.NewError(fmt.Errorf("invalid role provided. Role %v", role), rent.WithInvalidPayload())
//	}
//
//	user := &SignUpForm{
//		ID:           NewUserID(),
//		EmailAddress: addr,
//		Password:     password,
//		FirstName:    firstName,
//		LastName:     lastName,
//		Status:       UserActive,
//	}
//
//	err = u.uc.SignUpManager(u.registerUserToAccount(accountID, user, r))
//	if err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}
//
//func (u *SignUpManager) registerUserToAccount(accountID AccountID, user *SignUpForm, role Role) RegistrationFunc {
//	return func(us signUpSteps) error {
//		err := us.RegisterUser(user)
//		if err != nil {
//			return err
//		}
//
//		err = us.AddUserToAccount(accountID, user.ID, role)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//}
