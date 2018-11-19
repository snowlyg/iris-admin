package models

import (
	"github.com/pborman/uuid"
)

const (
	DEFAULT_EXPIRATION = 3600 // 24*60*60s=86400
)

type Session struct {
	Id         string `json:"sessionId"`
	UserId     int64  `json:"userId"`
	Username   string `json:"username"`
}

func NewSession(userId int64, username string) *Session {
	return &Session{
		Id:         uuid.New(),
		UserId:     userId,
		Username:   username,
	}
}
