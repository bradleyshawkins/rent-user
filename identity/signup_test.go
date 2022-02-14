package identity_test

import (
	"errors"
	"net/mail"
	"testing"

	"github.com/matryer/is"

	"github.com/bradleyshawkins/rent-user/identity"
)

type mockUserCreatorService struct {
	RegisterUserUserRegistration       *identity.User
	RegisterUserCredentials            *identity.Credentials
	RegisterUserError                  error
	RegisterAccountUserID              identity.UserID
	RegisterAccountAccountRegistration *identity.Account
	RegisterAccountError               error
	AddUserToAccountAccountID          identity.AccountID
	AddUserToAccountUserID             identity.UserID
	AddUserToAccountRole               identity.Role
	AddUserToAccountError              error
}

func (m *mockUserCreatorService) RegisterUser(u *identity.User, c *identity.Credentials) error {
	m.RegisterUserUserRegistration = u
	m.RegisterUserCredentials = c
	return m.RegisterUserError
}

func (m *mockUserCreatorService) RegisterAccount(a *identity.Account) error {
	m.RegisterAccountAccountRegistration = a
	return m.RegisterAccountError
}

func (m *mockUserCreatorService) AddUserToAccount(aID identity.AccountID, pID identity.UserID, role identity.Role) error {
	m.AddUserToAccountAccountID = aID
	m.AddUserToAccountUserID = pID
	m.AddUserToAccountRole = role
	return m.AddUserToAccountError
}

type mockUserCreator struct {
	mpcs *mockUserCreatorService
}

func (m *mockUserCreator) SignUp(suf *identity.SignUpForm) error {
	return suf.SignUp(m.mpcs)
}

func TestRegisterUser(t *testing.T) {
	i := is.New(t)
	mpcs := &mockUserCreatorService{}
	mpc := &mockUserCreator{mpcs}
	registrar := identity.NewSignUpManager(mpc)
	emailAddress, _ := mail.ParseAddress("email.address@test.com")
	firstName := "First"
	lastName := "Last"
	username := "username"
	password := "Password"

	user, err := registrar.SignUp(username, password, emailAddress, firstName, lastName)
	if err != nil {
		t.Fatal("Unexpected error occurred. Error:", err)
	}

	if user == nil {
		t.Fatal("signUpForm was nil.")
	}

	t.Log("user:", user)

	i.Equal(user.EmailAddress.Address, emailAddress.Address)
	i.Equal(user.FirstName, firstName)
	i.Equal(user.LastName, lastName)
	i.True(!user.ID.IsZero())

}

func TestRegisterUser_Fail(t *testing.T) {
	tests := []struct {
		name               string
		createUserError    error
		createAccountError error
	}{
		{
			name:               "Create user Fails",
			createUserError:    errors.New("unable to create user"),
			createAccountError: nil,
		},
		{
			name:               "Create account Fails",
			createUserError:    nil,
			createAccountError: errors.New("unable to create account"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mpcs := &mockUserCreatorService{
				RegisterUserError:    tt.createUserError,
				RegisterAccountError: tt.createAccountError,
			}
			mpc := &mockUserCreator{mpcs: mpcs}

			registrar := identity.NewSignUpManager(mpc)

			emailAddress, _ := mail.ParseAddress("email.address@test.com")

			user, err := registrar.SignUp("username", "password", emailAddress, "First", "Last")
			if err == nil {
				t.Error("Received a nil error. Should have received an error")
			}
			if user != nil {
				t.Error("user was non-nil. user:", user)
			}
			t.Logf("error: %v", err)
		})
	}
}
