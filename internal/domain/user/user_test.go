package user_test

import (
	"strings"
	"testing"

	"gravitum/internal/domain/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserNew(t *testing.T) {
	testCases := []struct {
		id        uuid.UUID
		login     string
		firstName string
		lastName  string
		err       error
	}{
		{
			id:        uuid.New(),
			login:     "login",
			firstName: "John",
			lastName:  "Doe",
			err:       nil,
		},

		{
			id:        uuid.New(),
			login:     "LoGiN",
			firstName: "John",
			lastName:  "Doe",
			err:       nil,
		},
		{
			id:        uuid.New(),
			login:     "",
			firstName: "John",
			lastName:  "Doe",
			err:       user.ErrLogin,
		},
		{
			id:        uuid.New(),
			login:     "login1",
			firstName: "",
			lastName:  "Doe",
			err:       user.ErrFirstName,
		},
		{
			id:        uuid.New(),
			login:     "login2",
			firstName: "John",
			lastName:  "",
			err:       user.ErrLastName,
		},
	}

	for _, tc := range testCases {
		u, err := user.New(tc.id, tc.login, tc.firstName, tc.lastName)
		assert.Equal(t, tc.err, err)
		if err != nil {
			assert.Equal(t, uuid.Nil, u.ID())
			assert.Equal(t, "", u.Login())
			assert.Equal(t, "", u.FirsName())
			assert.Equal(t, "", u.LastName())
		} else {
			assert.Equal(t, tc.id, u.ID())
			assert.Equal(t, strings.ToLower(strings.TrimSpace(tc.login)), u.Login())
			assert.Equal(t, tc.firstName, u.FirsName())
			assert.Equal(t, tc.lastName, u.LastName())
		}
	}
}
