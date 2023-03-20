package response

import "errors"

var (
	// token过期错误提示
	ErrExpiredToken = errors.New("token已过期")
	// token非法错误提示
	ErrInvalidToken = errors.New("token格式不合法")
	// 参数类型错误
	ErrInvalidParam = errors.New("参数类型错误")
	// 用户名或密码错误
	ErrWrongPassword = errors.New("用户名或者密码错误")
	// token生成失败
	ErrCreateToken = errors.New("token生成失败")
	// 密码加密失败
	ErrEncrypt = errors.New("密码加密失败")
	// 用户已存在
	ErrUserExist = errors.New("用户已存在")
	// 用户不存在
	ErrUserNotExist = errors.New("用户不存在")
	// 邮箱已被注册
	ErrEmailExist = errors.New("邮箱已被注册")
	// 用户名已被注册
	ErrUsernameExist = errors.New("用户名已被注册")

	// 验证码过期
	ErrVerifyCodeExpired = errors.New("验证码已失效")
	// 验证码不存在
	ErrVerifyCodeNotExist = errors.New("没有找到该邮箱对应的验证码信息，请先获取邮箱验证码")
	// 验证码错误
	ErrVerifyCode = errors.New("验证码错误")
	// 验证码错误次数超过上限
	ErrVerifyCodeMaxTry = errors.New("验证码错误，且已达到最大重试次数，验证码失效，请重新获取验证码")

	// 路线已存在
	ErrRouteExist = errors.New("添加失败，路线已存在")
	// 路线不存在
	ErrRouteNotExist = errors.New("路线不存在")
	// 路线列表为空
	EmptyRouteList = errors.New("当前路线列表为空")

	// 列车已存在
	ErrTrainExist = errors.New("添加失败，列车已存在")
	// 列车不存在
	ErrTrainNotExist = errors.New("列车不存在")
	// 列车列表为空
	EmptyTrainList = errors.New("当前列车列表为空")

	// 售票信息已存在
	ErrTicketExist = errors.New("添加失败，该车票已在售卖中")
	// 售票信息不存在
	ErrTicketNotExist = errors.New("查询不到对应的车票")
	// 发售车票列表为空
	EmptyTicketList = errors.New("当前发售的车票列表为空")
	// 在售车票列表为空
	EmptyOnSaleTicketList = errors.New("当前在售的车票列表为空")
	// 当前车次没有余票
	ErrNoRemainingTicket = errors.New("当前车次没有余票")

	// 订单生成失败
	ErrAddOrderFailed = errors.New("订单添加失败")
	// 用户对同一车次的车票重复下订单，返回错误
	ErrSameOrderExist = errors.New("该用户已经购买过该车次的车票")
	// 订单不存在
	ErrOrderNotExist = errors.New("订单不存在")
	// 订单列表为空
	EmptyOrderList = errors.New("订单列表为空")
	// 更新订单状态失败
	ErrUpdateOrderStatus = errors.New("更新订单状态失败")

	// 数据库操作错误
	ErrDbOperation = errors.New("数据库操作错误")

	// Redis缓存操作错误
	ErrRedisOperation = errors.New("Redis缓存操作错误")

	// 消息转换成Json格式失败
	ErrFailedChangeToJson = errors.New("数据编码失败")
	// 车票库存修改失败
	ErrFailedSubStock = errors.New("车票库存修改失败")
)
