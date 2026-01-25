package controller

const userSessionCookieName = "baton-session"

type loginRequest struct {
	Login string `json:"login" validate:"required,min=3"`
	Pass  string `json:"password" validate:"required,min=3"`
}

type errorModel struct {
	Message string `json:"message"`
}

type buttonPowerRequest struct{
	User int `query:"user" validate:"required,gt=0"`
}