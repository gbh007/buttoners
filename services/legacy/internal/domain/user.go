package domain

type User struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Token string `json:"-" gorm:"uniqueIndex"`
}
