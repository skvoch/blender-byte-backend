package model

// UserData ...
type UserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Validate ...
func (r *UserData) Validate() bool {
	if len(r.Email) < 6 {
		return false
	}

	if len(r.Login) < 6 {
		return false
	}

	if len(r.Password) < 6 {
		return false
	}

	return true
}

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
