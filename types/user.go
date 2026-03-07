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

type UpdateUserParams struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

func (p UpdateUserParams) UpdateUserToBSON() bson.M {
	update := bson.M{}
	if p.FirstName != nil && len(*p.FirstName) > 0 {
		update["firstName"] = *p.FirstName
	}
	if p.LastName != nil && len(*p.LastName) > 0 {
		update["lastName"] = *p.LastName
	}
	return update
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("firstName must be at least %d characters", minFirstNameLength)
	}
	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("lastName must be at least %d characters", minLastNameLength)
	}
	if len(params.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters", minPasswordLength)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("invalid email")
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
