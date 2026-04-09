package dto

type UserRegisterReq struct {
	Username        string `json:"username" validate:"required,min=3,max=20"`
	Password        string `json:"password" validate:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserLoginReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
