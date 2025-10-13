package link

import "time"

type Link struct {
	Id string
	Code string
	Clicks int
	Url string
	UserId int
	CreateAt time.Time
}