package rest

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/bradleyshawkins/berror"

	"github.com/google/uuid"
)

const (
	accountID = "accountID"
	userID    = "userID"
)

type RegisterUserRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	EmailAddress string `json:"emailAddress"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

func (r *RegisterUserRequest) validate() error {
	var invalidFields []berror.InvalidField
	if len(r.Username) == 0 {
		invalidFields = append(invalidFields, berror.InvalidField{
			Field:  "username",
			Reason: berror.ReasonMissing,
		})
	}
	if len(r.EmailAddress) == 0 {
		invalidFields = append(invalidFields, berror.InvalidField{
			Field:  "emailAddress",
			Reason: berror.ReasonMissing,
		})
	}
	if len(r.Password) == 0 {
		invalidFields = append(invalidFields, berror.InvalidField{
			Field:  "password",
			Reason: berror.ReasonMissing,
		})
	}
	if len(r.FirstName) == 0 {
		invalidFields = append(invalidFields, berror.InvalidField{
			Field:  "firstName",
			Reason: berror.ReasonMissing,
		})
	}
	if len(r.LastName) == 0 {
		invalidFields = append(invalidFields, berror.InvalidField{
			Field:  "lastName",
			Reason: berror.ReasonMissing,
		})
	}

	if len(invalidFields) > 0 {
		return berror.New("invalid fields provided", berror.WithInvalidFields(invalidFields...))
	}
	return nil
}

type RegisterUserResponse struct {
	AccountID uuid.UUID `json:"accountID"`
	UserID    uuid.UUID `json:"userID"`
}

func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	var rr RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&rr)
	if err != nil {
		return berror.Wrap(err, berror.WithInvalidPayload(), berror.WithMessage("unable to decode request"))
	}

	if err := rr.validate(); err != nil {
		return err
	}

	emailAddress, err := mail.ParseAddress(rr.EmailAddress)
	if err != nil {
		return berror.Wrap(err, berror.WithInvalidFields(berror.InvalidField{
			Field:  "emailAddress",
			Reason: berror.ReasonInvalid,
		}))
	}

	user, err := s.signUpManager.SignUp(rr.Username, rr.Password, emailAddress, rr.FirstName, rr.LastName)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	// TODO: Figure out response object
	err = json.NewEncoder(w).Encode(RegisterUserResponse{UserID: user.ID.AsUUID(), AccountID: uuid.New()})
	if err != nil {
		return berror.WrapInternal(err, "unable to serialize response")
	}

	return nil
}
