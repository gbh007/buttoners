package dto

// Не совсем корректное расположение данных, более чистое решение разместить в DTO auth
type UserInfo struct {
	ID int64 `json:"id"`
}
