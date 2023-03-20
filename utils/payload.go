package utils

import (
	"Project/MyProject/response"
	"github.com/google/uuid"
	"time"
)

// Payload包含了用于验证用户信息的用户Id和用户名信息，以及token的一个过期时间
type Payload struct {
	UUID        uuid.UUID
	UserID      int
	Email       string
	ExpiredTime time.Time
}

// 根据传入的用户名、用户id、有效时间，生成一个Payload信息
func NewPayload(userID int, email string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		UUID:        tokenID,
		UserID:      userID,
		Email:       email,
		ExpiredTime: time.Now().Add(duration),
	}
	return payload, nil
}

// 验证Payload是否以及过期
func (p *Payload) Verify() error {
	if time.Now().After(p.ExpiredTime) {
		return response.ErrExpiredToken
	}
	return nil
}
