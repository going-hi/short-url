package auth

import "time"


type RefreshToken struct {
	Id int
	Token string
	UserId int
	CreatedAt time.Time
}

