package db

import (
	"database/sql"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db *sql.DB
)

func NewDB() {
	// Example DSN: "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "host=localhost user=postgres password=root dbname=short_url_db port=5432 sslmode=disable"

	var err error
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	Db, err = DB.DB()
	Db.SetMaxIdleConns(10)
	Db.SetMaxOpenConns(100)

	if err != nil {
		zap.L().Fatal("failed connecting to db", zap.Error(err))
	}
}

func StopDB()  {
	err := Db.Close()
	if err != nil {
		zap.L().Fatal("failed close db", zap.Error(err))
	}
}
