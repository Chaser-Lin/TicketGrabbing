package server

import (
	"Project/MyProject/cache"
	"Project/MyProject/config"
	"Project/MyProject/dao"
	"Project/MyProject/event"
	h "Project/MyProject/handler"
	"Project/MyProject/middleware"
	"Project/MyProject/services"
	"Project/MyProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

type Server struct {
	config *config.Config
	//store      db.Store
	tokenMaker    utils.PasetoMaker
	kafkaProducer *event.Producer
	kafkaConsumer *event.Consumer
	router        *gin.Engine
}

func NewServer(config *config.Config) (*Server, error) {
	tokenMaker := utils.NewTokenMaker(config.TokenSymmetricKey)
	kafkaProducer, err := event.NewProducer()
	if err != nil {
		return nil, err
	}
	kafkaConsumer, err := event.NewConsumer()
	if err != nil {
		return nil, err
	}
	server := &Server{
		config: config,
		//store:  store,
		tokenMaker:    tokenMaker,
		kafkaProducer: kafkaProducer,
		kafkaConsumer: kafkaConsumer,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	userService := services.NewUserServices(dao.NewUserDao())
	PassengerService := services.NewPassengerServices(dao.NewPassengerDao())
	routeService := services.NewRouteServices(dao.NewRouteDao())
	trainService := services.NewTrainServices(dao.NewTrainDao())
	ticketService := services.NewTicketServices(dao.NewTicketDao())
	orderService := services.NewOrderServices(dao.NewOrderDao(), dao.NewTicketDao(), dao.NewPassengerDao())
	emailService := services.NewEmailServices(services.NewVerifyCodeManager())
	spikeService := services.NewSpikeService(server.kafkaProducer)

	var userHandler = h.NewUserHandler(userService, PassengerService, ticketService, orderService, emailService, spikeService)
	var adminHandler = h.NewAdminHandler(routeService, trainService, ticketService)

	router := gin.Default()

	// 用户相关路由组
	userGroup := router.Group("/users")
	{
		userGroup.POST("/verify_code", userHandler.SendVerifyCode)                            // 发送验证码
		userGroup.POST("/", userHandler.Register)                                             // 用户注册
		userGroup.POST("/login", userHandler.Login)                                           // 用户登录
		userGroup.GET("/", middleware.Auth(), userHandler.GetUserInfo)                        // 根据token查询用户信息
		userGroup.POST("/update_password", middleware.Auth(), userHandler.UpdateUserPassword) // 修改密码
		userGroup.POST("/update_email", middleware.Auth(), userHandler.UpdateUserEmail)       // 修改邮箱
	}
	router.GET("/tickets/search", middleware.Auth(), userHandler.ListTicketsOnSale) // 用户通过起点和终点查询在售车票信息

	// 抢票接口使用限流中间件，每秒限制5000个请求，使用令牌桶算法，每秒填充5000个令牌
	router.POST("/spike", middleware.Auth(), middleware.RateLimit(time.Second, 5000, 5000), userHandler.BuyTicket) // 抢票接口
	//router.POST("/spike", middleware.RateLimit(time.Second, 5000, 5000), userHandler.BuyTicket) // 抢票接口

	router.POST("/tokens/renew_access", userHandler.RenewAccessToken) // 通过refreshToken刷新accessToken

	// 乘客相关路由组
	passengerGroup := router.Group("/passengers").Use(middleware.Auth())
	{
		passengerGroup.POST("/", userHandler.AddPassenger)                        // 添加用户的乘客信息
		passengerGroup.GET("/", userHandler.ListUserPassengers)                   // 获取用户添加的乘客列表
		passengerGroup.GET("/:user_passenger_id", userHandler.GetPassenger)       // 获取用户的乘客信息
		passengerGroup.DELETE("/:user_passenger_id", userHandler.DeletePassenger) // 删除用户的乘客信息
	}

	// 订单相关路由组
	orderGroup := router.Group("/orders").Use(middleware.Auth())
	{
		orderGroup.POST("/pay", userHandler.PayOrder)            // 用户支付车票订单
		orderGroup.POST("/cancel", userHandler.CancelOrder)      // 用户取消车票订单
		orderGroup.GET("/", userHandler.ListOrders)              // 用户查询所有订单
		orderGroup.GET("/:order_id", userHandler.GetOrder)       // 用户通过订单id查询具体订单
		orderGroup.DELETE("/:order_id", userHandler.DeleteOrder) // 用户删除指定id的订单
	}

	// 路线相关路由组
	routeGroup := router.Group("/routes").Use(middleware.AdminAuth())
	{
		routeGroup.POST("/", adminHandler.AddRoute)  // 管理员添加路线
		routeGroup.GET("/", adminHandler.ListRoutes) // 管理员查询所有路线
		//routeGroup.GET("/:start/:end", adminHandler.GetRoute)   // 管理员通过起点和终点查询具体路线
		routeGroup.GET("/:route_id", adminHandler.GetRouteByID)   // 管理员通过路线id查询具体路线
		routeGroup.DELETE("/:route_id", adminHandler.DeleteRoute) // 管理员通过路线id删除路线
	}

	// 列车相关路由组
	trainGroup := router.Group("/trains").Use(middleware.AdminAuth())
	{
		trainGroup.POST("/", adminHandler.AddTrain)               // 管理员添加列车
		trainGroup.GET("/", adminHandler.ListTrains)              // 管理员查询所有列车
		trainGroup.GET("/:train_id", adminHandler.GetTrain)       // 管理员通过列车id查询具体列车信息
		trainGroup.DELETE("/:train_id", adminHandler.DeleteTrain) // 管理员通过列车id删除列车信息
	}

	// 车票相关路由组
	ticketGroup := router.Group("/tickets").Use(middleware.AdminAuth())
	{
		ticketGroup.POST("/", adminHandler.AddTicket)                    // 管理员添加售票信息
		ticketGroup.GET("/", adminHandler.ListTickets)                   // 管理员通查询某一路线的售票信息
		ticketGroup.GET("/onsale", adminHandler.ListTicketsOnSale)       // 管理员查询某一路线的在售车票信息
		ticketGroup.GET("/all", adminHandler.ListAllTickets)             // 管理员通查询所有路线的售票信息
		ticketGroup.GET("/allonsale", adminHandler.ListAllTicketsOnSale) // 管理员查询所有路线的在售车票信息
		ticketGroup.POST("/stopsell", adminHandler.StopSell)             // 管理员停止售票
	}

	server.router = router
}

func (s *Server) StartKafkaConsumer() {
	kafkaService := services.NewKafkaMQService(s.kafkaConsumer)
	ticketService := services.NewTicketServices(dao.NewTicketDao())
	orderService := services.NewOrderServices(dao.NewOrderDao(), dao.NewTicketDao(), dao.NewPassengerDao())

	kafkaService.StartConsumer(orderService, ticketService)
}

func (s *Server) Start(addr string) error {
	go s.StartKafkaConsumer()    // 启动kafka消费者
	go s.AutoDeleteExpireOrder() // 启动定时任务，每秒更新一次过期订单信息
	err := s.LoadTicketStocks()  // 缓存预热
	if err != nil {
		return err
	}
	return s.router.Run(addr)
}

// LoadTicketStocks 缓存预热，在系统开始运行时先读取数据库中的余票信息
func (s *Server) LoadTicketStocks() error {
	tickerService := services.NewTicketServices(dao.NewTicketDao())

	tickets, err := tickerService.GetAllTickets()
	if err != nil {
		log.Println("LoadTicketStocks ListAllTickets err: ", err)
		return err
	}

	for _, ticket := range tickets {
		fmt.Printf("ticketID:(%d), stock:(%d)\n", ticket.TicketID, ticket.Stock)
		err := cache.AddStock(cache.GetStockKey(ticket.TicketID), ticket.Stock)
		if err != nil {
			log.Println("LoadTicketStocks cache.AddStock err: ", err)
			return err
		}
	}
	return nil
}

func (s *Server) AutoDeleteExpireOrder() {
	ticker := time.NewTicker(time.Second) // 每秒删除一次过期订单记录
	defer ticker.Stop()
	orderService := services.NewOrderServices(dao.NewOrderDao(), dao.NewTicketDao(), dao.NewPassengerDao())
	ticketService := services.NewTicketServices(dao.NewTicketDao())
	for range ticker.C {
		now := time.Now().Unix()
		orderIDs, err := cache.GetExpiredOrder("0", strconv.Itoa(int(now)))
		if err != nil {
			log.Println("AutoDeleteExpireOrder cache.GetExpiredOrder err: ", err)
			continue
		}
		// 不要手动删数据库，否则会出现redis和mysql数据不一致的情况，程序无法正常运行
		for _, orderID := range orderIDs {
			fmt.Println("expired orderID: ", orderID)
			passengerID, ticketID, err := orderService.GetOrderPassengerAndTicketID(orderID)
			if err != nil {
				log.Printf("AutoDeleteExpireOrder orderService.GetOrder error, orderID:(%v), err:(%v)", orderID, err)
				continue
			}
			// 把有效订单记录删除，避免用户在没有有效订单的情况下也购买不了车票
			err = cache.DeleteOrderLimit(passengerID, ticketID)
			if err != nil {
				log.Printf("AutoDeleteExpireOrder cache.DeleteOrderLimit error, ticketID:(%v), err:(%v)", ticketID, err)
			}

			updateOrderStatusService := services.UpdateOrderStatusService{
				OrderID: orderID,
				Status:  2, // 订单状态为 2 表示已过期
			}
			err = orderService.UpdateOrderStatus(&updateOrderStatusService)
			if err != nil {
				log.Printf("AutoDeleteExpireOrder orderService.UpdateOrderStatus error, orderID:(%v), err:(%v)", orderID, err)
				continue
			}
			// 订单过期后需要将车票库存+1
			err = ticketService.AddNumberOne(ticketID)
			if err != nil {
				log.Printf("AutoDeleteExpireOrder ticketService.AddNumberOne error, ticketID:(%v), err:(%v)", ticketID, err)
			}
		}
		//cache.RemoveExpiredOrder("0", strconv.Itoa(int(now)))
	}
}

//func errResponse(err error) gin.H {
//	return gin.H{"error": err.Error()}
//}
