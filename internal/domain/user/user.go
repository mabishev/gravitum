package user

import (
	"strings"

	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID
	login     string
	firstName string
	lastName  string
}

func New(id uuid.UUID, login string, firstName string, lastName string) (User, error) {
	login = strings.ToLower(strings.TrimSpace(login))
	if login == "" {
		return User{}, ErrLogin
	}

	firstName = strings.TrimSpace(firstName)
	if firstName == "" {
		return User{}, ErrFirstName
	}

	lastName = strings.TrimSpace(lastName)
	if lastName == "" {
		return User{}, ErrLastName
	}

	return User{
		id:        id,
		login:     login,
		firstName: firstName,
		lastName:  lastName,
	}, nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Login() string {
	return u.login
}

func (u *User) FirsName() string {
	return u.firstName
}

func (u *User) LastName() string {
	return u.lastName
}
