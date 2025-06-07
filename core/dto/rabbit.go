package dto

type RabbitMQData struct {
	RequestID string `json:"request_id"`
	UserID    int64  `json:"user_id"`
	Chance    int64  `json:"chance"`
	Duration  int64  `json:"duration"`
}
