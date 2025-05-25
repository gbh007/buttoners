package logclient

import "time"

const (
	ContentType  = "application/json; charset=utf-8"
	ActivityPath = "/api/v1/activity"
)

// /api/v1/activity

type ActivityRequest struct {
	UserID int64 `json:"user_id"`
}

type ActivityResponse struct {
	RequestCount int64     `json:"request_count"`
	LastRequest  time.Time `json:"last_request"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Details string `json:"details"`
}
