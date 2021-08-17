package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	init_p()
	init_storage()

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
		panic("Failed to connect database")
	}

	db.AutoMigrate(&User{}, &Role{}, &File{})
	initDefaultRoles()
	initRootUser()
}
