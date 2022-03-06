package database

import (
	"os"

	"github.com/sing3demons/todoapi/todo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	database, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(&todo.Todo{})
	db = database
}

func GetDB() *gorm.DB {
	return db
}
