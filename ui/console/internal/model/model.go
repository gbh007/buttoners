package model

type Connection struct {
	Addr     string `json:"addr"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
