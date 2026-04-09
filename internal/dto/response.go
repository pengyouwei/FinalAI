package dto

type UserRegisterRes struct {
	Username string `json:"username"`
}

type UserLoginRes struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
