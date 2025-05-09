package authclient

// /api/v1/login

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// /api/v1/logout

type LogoutRequest struct {
	Token string `json:"token"`
}

// /api/v1/register

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// /api/v1/info

type InfoRequest struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Login  string `json:"login"`
	UserID int64  `json:"user_id"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}
