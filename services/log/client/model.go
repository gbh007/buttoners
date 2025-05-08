package client

import "time"

type LogData struct {
	RequestCount int64
	LastRequest  time.Time
}
