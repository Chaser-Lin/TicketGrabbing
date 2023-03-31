package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

const (
	limitKeyPrefix = "ticket_stock_"
	orderLimitKey  = "order_limit"
)

func GetStockKey(ticketID int) string {
	return fmt.Sprintf("%s%d", limitKeyPrefix, ticketID)
}

// Lua脚本，用于查询剩余车票并将库存-1，保证操作的原子性
var script = redis.NewScript(`
local num = redis.call("GET", KEYS[1])
if num == 0 then
    return false
end

local current = redis.call("INCRBY", KEYS[1], -1)
if current < 0 then
	redis.call("INCRBY", KEYS[1], 1)
    return false
else
    return true
end`)

// Limit 用于redis缓存的余票检查和预扣库存
func Limit(key string) bool {
	eval := script.Run(RedisClient, []string{key}, []string{})
	ok, err := eval.Bool()
	if err != nil {
		log.Println("limit error", err.Error())
	}
	return ok
}

// AddStock 为车票设置库存
func AddStock(key string, num uint32) error {
	return RedisClient.Set(key, num, 0).Err()
}

// StockAddOne 用于用户购买过程中出问题时补偿库存
func StockAddOne(key string) error {
	log.Println(key, "stock + 1")
	return RedisClient.Incr(key).Err()
}

// OrderLimit 用户在购买车票时，判断是否有限制
func OrderLimit(passengerID int, ticketID int) (bool, error) {
	userTicket := fmt.Sprintf("%d-%d", passengerID, ticketID)
	return RedisClient.SIsMember(orderLimitKey, userTicket).Result()
}

// AddOrderLimit 添加用户购票限制，在用户下单成功时调用
func AddOrderLimit(passengerID int, ticketID int) error {
	userTicket := fmt.Sprintf("%d-%d", passengerID, ticketID)
	return RedisClient.SAdd(orderLimitKey, userTicket).Err()
}

// DeleteOrderLimit 删除用户购票限制，在用户取消订单以及订单过期时调用
func DeleteOrderLimit(passengerID int, ticketID int) error {
	userTicket := fmt.Sprintf("%d-%d", passengerID, ticketID)
	return RedisClient.SRem(orderLimitKey, userTicket).Err()
}
