package model

// LoginResponse ...
type LoginResponse struct {
	PrivateUUID string `json:"private_uuid"`
}

// LoginRequest ...
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// WhoamiResponse ...
type WhoamiResponse struct {
	Login string `json:"login"`
}
