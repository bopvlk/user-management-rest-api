package datastore

import (
	"fmt"
	"log"

	"git.foxminded.com.ua/3_REST_API/interal/config"
	"git.foxminded.com.ua/3_REST_API/interal/interface/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(c *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(mysql:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.MysqlUser, c.MysqlPassword, c.MysqlDatabase)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&repository.DBUser{})

	return db
}
