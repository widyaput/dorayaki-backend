package database

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// DB is reusable gorm sql connection.
	DB *gorm.DB
)

// ConnectDB connects this application to database instance.
func ConnectDB() error {
	h := os.Getenv("MYSQL_HOST")
	u := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	p := os.Getenv("MYSQL_PORT")
	d := os.Getenv("MYSQL_DATABASE")
	dsn := u + ":" + pwd + "@tcp(" + h + ":" + p + ")/" + d + "?charset=utf8mb4&parseTime=True&loc=Local"

	dbConnection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = dbConnection
	return nil
}
