package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DB_USER     string = "root"
	DB_PASSWORD string = "root"
	DB_HOST     string = "127.0.0.1"
	DB_PORT     string = "7531"
	DB_NAME     string = "cmsdiy"
	DB_CHARSET  string = "utf8mb4"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=" + DB_CHARSET + "&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
}

func DB() *gorm.DB {
	return db
}
