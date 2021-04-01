package mysql

import "gorm.io/gorm"

type DBUtil interface {
	CreateDB()
	DropDB()
	GetUtilDB() *gorm.DB
}

type DB interface {
	GetDB() *gorm.DB
	ClearAllData()
}
