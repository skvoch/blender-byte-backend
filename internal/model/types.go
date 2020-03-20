package model

// RegisterRequest ...
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Validate ...
func (r *RegisterRequest) Validate() bool {
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
