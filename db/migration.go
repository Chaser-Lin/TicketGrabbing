package db

import (
	"Project/MyProject/dal/models"
	"fmt"
	"log"
)

// Migration 执行数据迁移
func Migration() {
	//自动迁移模式
	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&models.User{},
			&models.Train{},
			&models.Route{},
			&models.Order{},
			&models.Ticket{},
			&models.UserPassenger{},
			&models.Passenger{},
		)
	if err != nil {
		log.Fatal("创建数据库表失败")
	}
	fmt.Println("创建数据库表成功")
}
