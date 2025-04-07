package domain

import "time"

type Message struct {
	ID        string
	FromUser  string
	ToUser    string
	Content   string
	Timestamp time.Time
}