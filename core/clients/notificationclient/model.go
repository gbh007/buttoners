package notificationclient

import "time"

const (
	ContentType = "application/json; charset=utf-8"
	NewPath     = "/api/v1/new"
	ListPath    = "/api/v1/list"
	ReadPath    = "/api/v1/read"
)

type NewRequest struct {
	UserID  int64     `json:"user_id"`
	Kind    string    `json:"kind"`
	Level   string    `json:"level"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

type NotificationData struct {
	ID      int64     `json:"id"`
	Kind    string    `json:"kind"`
	Level   string    `json:"level"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

type ListRequest struct {
	UserID int64 `json:"user_id"`
}

type ListResponse struct {
	Notifications []NotificationData `json:"notifications"`
}

type ReadRequest struct {
	UserID int64 `json:"user_id"`
	ID     int64 `json:"id,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}
