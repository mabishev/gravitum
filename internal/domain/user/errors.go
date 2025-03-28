package user

import "errors"

var (
	ErrID        = errors.New("invalid user id")
	ErrLogin     = errors.New("invalid user login")
	ErrFirstName = errors.New("invalid user first name")
	ErrLastName  = errors.New("invalid user last name")

	ErrNotFound = errors.New("user not found")
)
