package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() (*gorm.DB) {
	var err error;
	cs:= os.Getenv("cs")
	db, err:= gorm.Open(mysql.Open(cs), &gorm.Config{})
	if err != nil {
		fmt.Println("error connecting database: error = %v", err)
		return nil
	}
	return db
}