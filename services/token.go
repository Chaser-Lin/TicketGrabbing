package services

import (
	"Project/MyProject/cache"
	"Project/MyProject/response"
	"Project/MyProject/utils"
	"time"
)

type RenewAccessTokenService struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

func (s *RenewAccessTokenService) RenewAccessToken() (string, error) {
	refreshPayload, err := utils.TokenMaker.VerifyToken(s.RefreshToken)
	if err != nil {
		return "", err
	}

	if time.Now().After(refreshPayload.ExpiredTime) {
		return "", response.ErrExpiredToken
	}

	// 判断传进来的refreshToken和是否还保存在redis缓存中，以及和redis缓存中保存的session是否相同
	session, err := cache.GetSession(refreshPayload.UserID)
	if err != nil {
		return "", response.ErrExpiredToken
	}

	if s.RefreshToken != session {
		return "", response.ErrInvalidRefreshToken
	}

	accessToken, _, err := utils.TokenMaker.CreateToken(refreshPayload.UserID, refreshPayload.Email, utils.AccessTokenDurationTime)
	if err != nil {
		return "", response.ErrCreateToken
	}

	return accessToken, nil
}
