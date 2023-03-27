package main

import (
	"Project/MyProject/cache"
	"Project/MyProject/config"
	"Project/MyProject/dal"
	"Project/MyProject/dal/models"
	"Project/MyProject/db"
	"Project/MyProject/event"
	"Project/MyProject/server"
	"Project/MyProject/service"
	"Project/MyProject/utils"
	"fmt"
	"log"
	"os"
)

func main() {
	f, _ := os.OpenFile("./fmt.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
	os.Stdout = f
	os.Stderr = f

	// 加载配置文件
	err := config.Load()
	if err != nil {
		log.Fatal("配置文件加载错误:", err)
	}

	// 初始化Mysql数据库连接
	err = db.Init(config.Conf.Mysql)
	if err != nil {
		log.Fatal("Mysql连接错误:", err)
	}
	service.CreateAdmin()

	// 初始化Redis缓存连接
	err = cache.Init(config.Conf.Redis)
	if err != nil {
		log.Fatal("Redis连接错误:", err)
	}

	// 初始化kafka配置信息
	event.Init(config.Conf.Kafka)

	//for i := 0; i < 7000; i++ {
	//    TestCreateUser()
	//}

	// 创建一个 server 结构体
	server, err := server.NewServer(config.Conf)
	if err != nil {
		log.Fatal("创建server失败:", err)
	}

	// 启动 server
	err = server.Start(config.Conf.ServerPort)
	if err != nil {
		log.Fatal("启动server失败:", err)
	}

	fmt.Println("start successfully")
}

// 测试数据库插入和查询，没问题
func TestCreateUser() {
	userDal := dal.NewUserDal()
	password := utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)

	user1 := &models.User{
		Username:       utils.RandomString(1),
		HashedPassword: hashedPassword,
		Email:          utils.RandomEmail(),
	}
	err = userDal.AddUser(user1)
	//require.NoError(t, err)
	if err != nil {
		fmt.Println("数据插入失败")
	}

	//user2, err := userDal.GetUserByEmail(user1.Email)
	//if err != nil {
	//    fmt.Println("数据查询失败")
	//}
	//fmt.Println(user1.UserID, user2.UserID)
	//fmt.Println(user1.Username, user2.Username)
	//fmt.Println(user1.HashedPassword, user2.HashedPassword)
	//fmt.Println(user1.Email, user2.Email)
}
