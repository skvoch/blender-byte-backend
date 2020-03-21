package model

// FailedValidationError ...
type FailedValidationError struct {
}

func (f *FailedValidationError) Error() string {
	return "Failed data structure validation"
}

// UserAlreadyExistError ...
type UserAlreadyExistError struct {
}

func (f *UserAlreadyExistError) Error() string {
	return "User already exist error"
}

// CannotFindUserError ...
type CannotFindUserError struct {
}

func (f *CannotFindUserError) Error() string {
	return "Cannot find user by login"
}
