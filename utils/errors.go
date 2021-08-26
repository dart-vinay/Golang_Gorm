package utils

import "errors"

var (
	ErrBadRequest = errors.New("Bad Request")
	ErrInternal = errors.New("Internal Server Error")
	ErrInvalidStudentId  = errors.New("Invalid student id")
	ErrInvalidTeacherId = errors.New("Invalid teacher id")
	ErrUnknownUserType = errors.New("Unknown user type")
	ErrUnauthorized = errors.New("User Unauthorized")
)
