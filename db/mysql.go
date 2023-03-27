package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Address  string `yaml:"Address"`
	Port     int    `yaml:"Port"`
	DBName   string `yaml:"DBName"`
}

const baseDsn = "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"

var DB *gorm.DB

func Init(mysqlConf Config) error {
	db, err := newDBClient(mysqlConf)
	if err != nil {
		return err
	}
	DB = db
	fmt.Println("数据库连接成功")

	// 数据库迁移
	Migration()

	return nil
}

// 根据config结构体中的内容初始化数据库连接
func newDBClient(mysqlConf Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(baseDsn, mysqlConf.Username, mysqlConf.Password, mysqlConf.Address, mysqlConf.Port, mysqlConf.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info), // 执行数据库操作时，在日志打印对应的sql语句
		SkipDefaultTransaction: true,                                // 跳过事务，减少链接开销
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
