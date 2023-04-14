package handler

import (
	"Project/MyProject/cache"
	R "Project/MyProject/response"
	"Project/MyProject/services"
	"Project/MyProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/*
	用户界面需要有的功能
1. 用户注册、登录、修改信息
2. 查询售票信息
3. 订单添加、查询、查询用户所有、从列表删除
4. 选择指定订单支付

*/

type UserHandler struct {
	UserService      services.UserServiceImplement
	PassengerService services.PassengerServiceImplement
	TicketService    services.TicketServiceImplement
	OrderService     services.OrderServiceImplement
	EmailService     services.EmailServiceImplement
	SpikeService     services.SpikeServiceImplement
}

// 用户交互界面所有服务接口
type UserInteractImplement interface {
	UserImplement
	PassengerImplement
	ListTicketsOnSale(ctx *gin.Context)
	OrderImplement
	SpikeImplement
	RenewAccessToken(ctx *gin.Context)
}

// 用户服务接口
type UserImplement interface {
	SendVerifyCode(ctx *gin.Context)
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	GetUserInfo(ctx *gin.Context)
	UpdateUserPassword(ctx *gin.Context)
	UpdateUserEmail(ctx *gin.Context)
}

// 乘车人服务接口
type PassengerImplement interface {
	AddPassenger(ctx *gin.Context)
	DeletePassenger(ctx *gin.Context)
	GetPassenger(ctx *gin.Context)
	ListUserPassengers(ctx *gin.Context)
}

// 订单服务接口
type OrderImplement interface {
	//BuyTicket(ctx *gin.Context)
	PayOrder(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
	ListOrders(ctx *gin.Context)
	GetOrder(ctx *gin.Context)
	DeleteOrder(ctx *gin.Context)
}

// 抢票服务接口
type SpikeImplement interface {
	BuyTicket(ctx *gin.Context)
}

func NewUserHandler(userService services.UserServiceImplement,
	passengerService services.PassengerServiceImplement,
	ticketService services.TicketServiceImplement,
	orderService services.OrderServiceImplement,
	emailService services.EmailServiceImplement,
	spikeService services.SpikeServiceImplement) UserInteractImplement {
	return &UserHandler{
		UserService:      userService,
		PassengerService: passengerService,
		TicketService:    ticketService,
		OrderService:     orderService,
		EmailService:     emailService,
		SpikeService:     spikeService,
	}
}

func (handler *UserHandler) RenewAccessToken(ctx *gin.Context) {
	var renewAccessTokenService services.RenewAccessTokenService
	if err := ctx.ShouldBind(&renewAccessTokenService); err == nil {
		if accessToken, err := renewAccessTokenService.RenewAccessToken(); err == nil {
			R.Ok(ctx, "更新token信息成功", gin.H{
				"access_token": accessToken,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) SendVerifyCode(ctx *gin.Context) {
	var sendVerifyCodeService services.SendVerifyCodeService
	if err := ctx.ShouldBind(&sendVerifyCodeService); err == nil {
		if exist, err := handler.UserService.CheckUserExist(sendVerifyCodeService.Email); err != nil || exist {
			R.Error(ctx, err.Error(), nil)
			return
		}
		if err = handler.EmailService.SendVerifyCode(&sendVerifyCodeService); err == nil {
			R.Ok(ctx, "发送成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	var loginService services.UserLoginService
	if err := ctx.ShouldBind(&loginService); err == nil {
		// TODO:登录成功需要根据是否为管理员跳转到不同页面
		if accessToken, refreshToken, isAdmin, err := handler.UserService.Login(&loginService); err == nil {
			R.Ok(ctx, "登录成功", gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"is_admin":      isAdmin,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) Register(ctx *gin.Context) {
	var registerService services.UserRegisterService
	if err := ctx.ShouldBind(&registerService); err == nil {
		// 判断邮箱是否被注册
		if exist, err := handler.UserService.CheckUserExist(registerService.Email); err != nil || exist {
			R.Error(ctx, err.Error(), nil)
			return
		}
		// 判断验证码是否正确
		if ok, err := handler.EmailService.CheckVerifyCode(registerService.Email, registerService.VerifyCode); err != nil || !ok {
			R.Error(ctx, err.Error(), nil)
			return
		}
		if err := handler.UserService.Register(&registerService); err == nil {
			R.Ok(ctx, "注册成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) UpdateUserEmail(ctx *gin.Context) {
	var updateUserInfoService services.UpdateEmailService
	if err := ctx.ShouldBind(&updateUserInfoService); err == nil {
		if userID, exists := ctx.Get("user_id"); exists {
			updateUserInfoService.UserID = userID.(int)
		} else {
			R.Error(ctx, err.Error(), nil)
			return
		}
		userInfo, err := handler.UserService.GetUserInfo(updateUserInfoService.UserID)
		if err != nil {
			R.Error(ctx, err.Error(), nil)
			return
		}
		if userInfo.Email == updateUserInfoService.Email {
			R.Error(ctx, "该邮箱与当前用户邮箱相同，不需要更改！", nil)
			return
		}
		// 判断邮箱是否被注册
		if exist, err := handler.UserService.CheckUserExist(updateUserInfoService.Email); err != nil || exist {
			R.Error(ctx, err.Error(), nil)
			return
		}
		// 判断验证码是否正确
		if ok, err := handler.EmailService.CheckVerifyCode(updateUserInfoService.Email, updateUserInfoService.VerifyCode); err != nil || !ok {
			R.Error(ctx, err.Error(), nil)
			return
		}
		if err := handler.UserService.UpdateEmail(&updateUserInfoService); err == nil {
			R.Ok(ctx, "邮箱修改成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) UpdateUserPassword(ctx *gin.Context) {
	var updateUserInfoService services.UpdatePasswordService
	if err := ctx.ShouldBind(&updateUserInfoService); err == nil {
		if userID, exists := ctx.Get("user_id"); exists {
			updateUserInfoService.UserID = userID.(int)
		} else {
			R.Error(ctx, err.Error(), nil)
			return
		}
		userInfo, err := handler.UserService.GetUserInfo(updateUserInfoService.UserID)
		if err != nil {
			R.Error(ctx, err.Error(), nil)
			return
		}
		if err = utils.CheckPassword(userInfo.HashedPassword, updateUserInfoService.OldPassword); err != nil {
			R.Error(ctx, "修改密码失败，原密码输入错误", nil)
			return
		}
		if updateUserInfoService.NewPassword != updateUserInfoService.NewPasswordAgain {
			R.Error(ctx, "两次输入的新密码不一致，请重新输入", nil)
			return
		}
		if err := handler.UserService.UpdatePassword(&updateUserInfoService); err == nil {
			R.Ok(ctx, "密码修改成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) UpdateUsername(ctx *gin.Context) {
	var updateUserInfoService services.UpdateUsernameService
	if err := ctx.ShouldBind(&updateUserInfoService); err == nil {
		if userID, exists := ctx.Get("user_id"); exists {
			updateUserInfoService.UserID = userID.(int)
		} else {
			R.Error(ctx, err.Error(), nil)
			return
		}
		if err := handler.UserService.UpdateUsername(&updateUserInfoService); err == nil {
			R.Ok(ctx, "用户名修改成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) GetUserInfo(ctx *gin.Context) {
	// 1. 通过token解析获取userID
	//accessToken := ctx.MustGet("token").(string)
	//userID := utils.ParseToken(accessToken)
	// 2. 通过url参数获取userID
	//userID, _ := strconv.Atoi(ctx.Param("user_id"))
	// 3. 在中间件直接将userID放入context中
	userID, exists := ctx.Get("user_id")
	if exists {
		if userInfo, err := handler.UserService.GetUserInfo(userID.(int)); err == nil {
			R.Ok(ctx, "获取用户信息成功", gin.H{
				"email":    userInfo.Email,
				"username": userInfo.Username,
				"user_id":  userInfo.UserID,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Error(ctx, "用户不存在", nil)
	}
}

func (handler *UserHandler) ListTicketsOnSale(ctx *gin.Context) {
	var listTicketService services.ListTicketsService
	if err := ctx.ShouldBind(&listTicketService); err == nil {
		if tickets, err := handler.TicketService.ListTickets(&listTicketService); err == nil {
			R.Ok(ctx, "查询当前在售车票信息成功", gin.H{
				"tickets": tickets,
			})
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) BuyTicket(ctx *gin.Context) {
	var spikeServiceReq services.SpikeServiceReq
	if err := ctx.ShouldBind(&spikeServiceReq); err == nil {
		userID, ok := ctx.Get("user_id")
		if ok {
			if ok, err = handler.PassengerService.CheckPassengerBelongToUser(spikeServiceReq.PassengerID, userID.(int)); !ok || err != nil {
				R.Error(ctx, err.Error(), nil)
				return
			}
			if err := handler.SpikeService.BuyTicket(userID.(int), spikeServiceReq.TicketID, spikeServiceReq.PassengerID); err == nil {
				R.Ok(ctx, "购票成功", nil)
			} else {
				R.Error(ctx, err.Error(), nil)
			}
		} else {
			R.Error(ctx, "系统错误，购票失败", nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) AddPassenger(ctx *gin.Context) {
	var addPassengerService services.AddPassengerService
	if err := ctx.ShouldBind(&addPassengerService); err == nil {
		userID, ok := ctx.Get("user_id")
		if ok {
			addPassengerService.UserID = userID.(int)
			if err := handler.PassengerService.AddPassenger(&addPassengerService); err == nil {
				R.Ok(ctx, "添加乘客信息成功", nil)
			} else {
				R.Error(ctx, err.Error(), nil)
			}
		} else {
			R.Error(ctx, "系统错误，添加乘客信息失败", nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) DeletePassenger(ctx *gin.Context) {
	var deletePassengerService services.DeletePassengerService
	var err error
	deletePassengerService.UserPassengerID, err = strconv.Atoi(ctx.Param("user_passenger_id"))
	if err != nil {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
		return
	}
	if err := handler.PassengerService.DeletePassenger(&deletePassengerService); err == nil {
		R.Ok(ctx, "删除乘客信息成功", nil)
	} else {
		R.Error(ctx, err.Error(), nil)
	}

}
func (handler *UserHandler) GetPassenger(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if exists {
		if userPassengerID, err := strconv.Atoi(ctx.Param("user_passenger_id")); err == nil {
			if passenger, err := handler.PassengerService.GetPassenger(userPassengerID); err == nil {
				if passenger.UserID != userID.(int) {
					R.Error(ctx, R.ErrPassengerNotExist.Error(), nil)
				} else {
					R.Ok(ctx, "查询乘客信息成功", gin.H{
						"passenger": passenger,
					})
				}
			} else {
				R.Error(ctx, err.Error(), nil)
			}
		} else {
			R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
		}
	} else {
		R.Error(ctx, "查询乘客信息失败", nil)
	}
}
func (handler *UserHandler) ListUserPassengers(ctx *gin.Context) {
	var listUserPassengersService services.ListUserPassengersService
	if err := ctx.ShouldBind(&listUserPassengersService); err == nil {
		userID, ok := ctx.Get("user_id")
		if ok {
			listUserPassengersService.UserID = userID.(int)
			if passengers, err := handler.PassengerService.ListPassengers(&listUserPassengersService); err == nil {
				R.Ok(ctx, "获取乘客信息列表成功", gin.H{
					"passengers": passengers,
				})
			} else {
				R.Error(ctx, err.Error(), nil)
			}
		} else {
			R.Error(ctx, "系统错误，获取乘客信息列表失败", nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) PayOrder(ctx *gin.Context) {
	var updateOrderStatusService services.UpdateOrderStatusService
	if err := ctx.ShouldBind(&updateOrderStatusService); err == nil {
		userID, exists := ctx.Get("user_id")
		if exists {
			order, err := handler.OrderService.GetOrder(updateOrderStatusService.OrderID)
			if err != nil {
				R.Error(ctx, err.Error(), nil)
				return
			}
			// token所属的用户id与订单的所有者不匹配，无法支付
			if order.UserID != userID.(int) {
				R.Error(ctx, fmt.Sprintf("支付失败，该订单的所有者为:(%d)，不属于用户:(%d)", order.UserID, userID.(int)), nil)
				return
			}
			if order.Status == 0 {
				updateOrderStatusService.Status = 1 // status = 1表示已支付订单
				if err := handler.OrderService.UpdateOrderStatus(&updateOrderStatusService); err == nil {
					R.Ok(ctx, "支付成功", nil)
				} else {
					R.Error(ctx, err.Error(), nil)
				}
			} else if order.Status == 1 {
				R.Error(ctx, "该订单已支付，请勿重复支付！", nil)
			} else if order.Status == 2 {
				R.Error(ctx, "支付失败，该订单已过期", nil)
			} else if order.Status == 3 {
				R.Error(ctx, "支付失败，该订单已被取消", nil)
			}
		} else {
			R.Error(ctx, "支付失败", nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) CancelOrder(ctx *gin.Context) {
	var updateOrderStatusService services.UpdateOrderStatusService
	if err := ctx.ShouldBind(&updateOrderStatusService); err == nil {
		userID, exists := ctx.Get("user_id")
		if exists {
			order, err := handler.OrderService.GetOrder(updateOrderStatusService.OrderID)
			if err != nil {
				R.Error(ctx, err.Error(), nil)
				return
			}
			// token所属的用户id与订单的所有者不匹配，无法取消订单
			if order.UserID != userID.(int) {
				R.Error(ctx, fmt.Sprintf("取消订单失败，该订单的所有者为:(%d)，不属于用户:(%d)", order.UserID, userID.(int)), nil)
				return
			}
			if err = cache.DeleteOrderLimit(order.PassengerID, order.TicketID); err != nil {
				R.Error(ctx, "订单取消失败", nil)
				return
			}
			if order.Status == 2 {
				R.Error(ctx, "取消订单失败，该订单已过期", nil)
				return
			} else if order.Status == 3 {
				R.Error(ctx, "该订单已被取消，请勿重复取消订单", nil)
				return
			}
			updateOrderStatusService.Status = 3 // status = 3表示已取消订单
			err = handler.OrderService.UpdateOrderStatus(&updateOrderStatusService)
			if err != nil {
				R.Error(ctx, err.Error(), nil)
				return
			}
			// 订单取消后需要将车票库存+1
			if err = handler.TicketService.AddNumberOne(order.TicketID); err == nil {
				msg := "订单取消成功"
				if order.Status == 1 {
					msg = "退票成功"
				}
				R.Ok(ctx, msg, nil)
			} else {
				R.Error(ctx, err.Error(), nil)
			}
		} else {
			R.Error(ctx, "订单取消失败", nil)
		}
	} else {
		R.Response(ctx, http.StatusBadRequest, http.StatusBadRequest, "参数错误", err.Error())
	}
}

func (handler *UserHandler) ListOrders(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if exists {
		if orders, err := handler.OrderService.ListOrders(userID.(int)); err == nil {
			R.Ok(ctx, "查询订单列表成功", gin.H{
				"orders": orders,
			})
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Error(ctx, "查询订单列表失败", nil)
	}
}

func (handler *UserHandler) GetOrder(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if exists {
		orderID := ctx.Param("order_id")
		order, err := handler.OrderService.GetOrderInfo(orderID)
		if err != nil {
			R.Error(ctx, err.Error(), nil)
			return
		}
		// token所属的用户id与订单的所有者不匹配，无法查看订单
		if order.UserID != userID.(int) {
			R.Error(ctx, fmt.Sprintf("查询订单失败，该订单的所有者为:(%d)，不属于用户:(%d)", order.UserID, userID.(int)), nil)
		} else {
			R.Ok(ctx, "查询订单成功", gin.H{
				"order": order,
			})
		}
	} else {
		R.Error(ctx, "查询订单失败", nil)
	}
}

func (handler *UserHandler) DeleteOrder(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if exists {
		orderID := ctx.Param("order_id")
		order, err := handler.OrderService.GetOrder(orderID)
		if err != nil {
			R.Error(ctx, err.Error(), nil)
			return
		}
		// token所属的用户id与订单的所有者不匹配，无法删除订单
		if order.UserID != userID.(int) {
			R.Error(ctx, fmt.Sprintf("删除订单失败，该订单的所有者为:(%d)，不属于用户:(%d)", order.UserID, userID.(int)), nil)
			return
		}
		if order.Status == 0 || order.Status == 1 {
			R.Error(ctx, "订单删除失败，该订单为有效订单，无法删除", nil)
			return
		}
		if err := handler.OrderService.DeleteOrder(orderID); err == nil {
			R.Ok(ctx, "订单删除成功", nil)
		} else {
			R.Error(ctx, err.Error(), nil)
		}
	} else {
		R.Error(ctx, "订单删除失败", nil)
	}
}
