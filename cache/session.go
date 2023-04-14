package cache

import (
	"Project/MyProject/utils"
	"fmt"
)

func SetSession(userID int, refreshToken string) error {
	userSessionKey := fmt.Sprintf("user_%d_refresh_token", userID)
	return RedisClient.Set(userSessionKey, refreshToken, utils.RefreshTokenDurationTime).Err()
}

func GetSession(userID int) (string, error) {
	userSessionKey := fmt.Sprintf("user_%d_refresh_token", userID)
	return RedisClient.Get(userSessionKey).Result()
}
