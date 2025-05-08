package server

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func hashString(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) }

func randomSHA256String() string {
	// TODO: не годится для нагрузок
	return hashString(time.Now().String())
}

func saltPassword(password, salt string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+salt)))
}
