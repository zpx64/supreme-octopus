package model

import (
	"time"
)

type UserNCred struct {
	User        User            `json:"user"`
	Credentials UserCredentials `json:"credentials"`
}

type User struct {
	UserId       int       `json:"user_id"`
	CreationDate time.Time `json:"creation_date"`
	Nickname     string    `json:"nickname"`
	Name         *string   `json:"name,omitempty"`
	Surname      *string   `json:"surname,omitempty"`
}

type UserCredentials struct {
	UserId   int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Pow      string `json:"pow"`
	// i really hate local pows
	// but i think we need it(
}

type UserToken struct {
	TokenId      int    `json:"token_id"`
	UserId       int    `json:"user_id"`
	DeviceId     string `json:"device_id"`
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	TokenDate    int64  `json:"token_date"`
}
