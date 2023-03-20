package cache

import (
	"github.com/go-redis/redis"
)

const orderExpireKey = "order_expire_time_list"

// AddOrderExpireTime 将订单添加到redis中，score作为过期时间，后续使用定时任务将过期订单删除
func AddOrderExpireTime(score float64, member string) error {
	return RedisClient.ZAdd(orderExpireKey, redis.Z{Score: score, Member: member}).Err()
}

// RemoveFinishOrder 将已经成功支付或者已经取消或者处理完过期状态的订单从过期列表删除
func RemoveFinishOrder(member string) error {
	return RedisClient.ZRem(orderExpireKey, member).Err()
}

// GetExpiredOrder 获取当前已经过期的订单id，max为当前时间戳，过期时间小于当前时间戳的都是已过期的
func GetExpiredOrder(min, max string) ([]string, error) {
	result := make([]string, 0)
	opt := redis.ZRangeBy{
		Min: min, // 最小分数
		Max: max, // 最大分数
	}
	expireOrderList, err := RedisClient.ZRangeByScore(orderExpireKey, opt).Result()

	if err != nil {
		return nil, err
	}

	for _, order := range expireOrderList {
		result = append(result, order)
	}
	return result, nil
}

// RemoveExpiredOrder 删除当前已经过期的订单id，max为当前时间戳，过期时间小于当前时间戳的都是已过期的
//func RemoveExpiredOrder(min, max string) error {
//	return RedisClient.ZRemRangeByScore(orderExpireKey, min, max).Err()
//}
