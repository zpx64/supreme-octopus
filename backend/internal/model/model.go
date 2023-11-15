package model

import (
  "time"
)

type UserNCred struct {
  User        User            `json:"user"`
  Credentials UserCredentials `json:"credentials"`
}

type User struct {
  Id           int       `json:"id"`
  Nickname     string    `json:"nickname"`
  CreationDate time.Time `json:"creation_date"`
  Name         *string   `json:"name,omitempty"`
  Surname      *string   `json:"surname,omitempty"`
}

type UserCredentials struct {
  Id       int    `json:"-"`
  Email    string `json:"email"`
  Password string `json:"password"`
  Pow      string `json:"pow"`
  // i really hate local pows
  // but i think we need it(
}

type UserToken struct {
  Id           int       `json:"-"`
  DeviceId     string    `json:"device_id"`
  RefreshToken string    `json:"refresh_token"`
  TokenDate    time.Time `json:"token_date"`
}
