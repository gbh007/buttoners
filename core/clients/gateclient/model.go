package gateclient

import "time"

type NotificationData struct {
	ID      int64
	Kind    string
	Level   string
	Title   string
	Body    string
	Created time.Time
}
