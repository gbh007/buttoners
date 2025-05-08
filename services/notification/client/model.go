package client

import "time"

const (
	ButtonKind   = "button"
	SuccessLevel = "success"
	ErrorLevel   = "error"
)

type Notification struct {
	ID      int64
	Kind    string
	Level   string
	Title   string
	Body    string
	Created time.Time
}
