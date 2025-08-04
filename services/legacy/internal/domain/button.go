package domain

import (
	"bytes"
	"strconv"
	"time"
)

type Button struct {
	UserID int    `json:"-" gorm:"primaryKey;autoIncrement:false"`
	User   User   `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Year   int    `json:"year" gorm:"primaryKey;autoIncrement:false"`
	Month  int    `json:"month" gorm:"primaryKey;autoIncrement:false"`
	Day    int    `json:"day" gorm:"primaryKey;autoIncrement:false"`
	Count  int64  `json:"count"`
	Text   string `json:"text" gorm:"-"`
}

func (b *Button) UpdateText() {
	buff := bytes.Buffer{}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	bTime := time.Date(b.Year, time.Month(b.Month), b.Day, 0, 0, 0, 0, time.UTC)

	diff := now.Sub(bTime)

	switch {
	case diff <= 0:
		b.Text = "Сегодня"

	case diff <= time.Hour*24:
		b.Text = "Вчера"

	case diff <= time.Hour*24*6:
		switch bTime.Weekday() {
		case time.Monday:
			b.Text = "Понедельник"
		case time.Tuesday:
			b.Text = "Вторник"
		case time.Wednesday:
			b.Text = "Среда"
		case time.Thursday:
			b.Text = "Четверг"
		case time.Friday:
			b.Text = "Пятница"
		case time.Saturday:
			b.Text = "Суббота"
		case time.Sunday:
			b.Text = "Воскресенье"
		}

	default:
		buff.WriteString(strconv.Itoa(b.Year))
		buff.WriteString(".")

		if b.Month < 10 {
			buff.WriteString("0")
		}

		buff.WriteString(strconv.Itoa(b.Month))
		buff.WriteString(".")

		if b.Day < 10 {
			buff.WriteString("0")
		}

		buff.WriteString(strconv.Itoa(b.Day))

		b.Text = buff.String()
	}
}
