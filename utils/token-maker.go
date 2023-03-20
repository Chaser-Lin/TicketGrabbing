package utils

// token包负责为登录用户生成token

import (
	"Project/MyProject/response"
	"github.com/o1egl/paseto"
	"time"
)

const (
	// 32位对称密钥
	TokenSymmetricKey = "12345678901234567980123456789012"
	// token有效期
	TokenDurationTime = 1 * time.Hour
)

var TokenMaker = NewTokenMaker(TokenSymmetricKey)

// PasetoMaker用于管理token
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// 生成一个tokenMaker管理token
func NewTokenMaker(symmetricKey string) PasetoMaker {
	return PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
}

// 根据传入的用户名、用户id、token有效时间，创建一个token
func (maker *PasetoMaker) CreateToken(userID int, email string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, email, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

// 验证token有效性
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, response.ErrInvalidToken
	}

	err = payload.Verify()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// 解析token，获取userID
func ParseToken(token string) (userID int) {
	payload, _ := TokenMaker.VerifyToken(token)

	return payload.UserID
}
