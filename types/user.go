package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 8
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}
	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Sprintf("firstName must be at least %d characters", minFirstNameLength))
	}
	if len(params.LastName) < minLastNameLength {
		errors = append(errors, fmt.Sprintf("lastName must be at least %d characters", minLastNameLength))
	}
	if len(params.Password) < minPasswordLength {
		errors = append(errors, fmt.Sprintf("password must be at least %d characters", minPasswordLength))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, fmt.Sprintf("invalid email"))
	}
	return errors
}

func isEmailValid(email string) bool {
	// A very basic email validation regex. For production use, consider using a more robust solution.
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

type User struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string        `bson:"firstName" json:"firstName"`
	LastName          string        `bson:"lastName" json:"lastName"`
	Email             string        `bson:"email" json:"email"`
	EncryptedPassword string        `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
