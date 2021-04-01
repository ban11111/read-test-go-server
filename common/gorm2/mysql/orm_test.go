package mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var cusLog = CustomerLog(func(msg string, items ...interface{}) {
	fmt.Print("\n"+msg)
	fmt.Print(items...)
	fmt.Print("\n")
})

func Test_newGormMysql(t *testing.T) {
	conf := GetDefaultDbConfig("hahaha_test")

	gm := newGormMysql(conf, cusLog, true)
	assert.NotNil(t, gm)
	gm.CreateDB()
	gm.GetUtilDB()
	defer gm.DropDB()

	gm2 := newGormMysql(conf,  cusLog, false)
	assert.NotNil(t, gm2)
	assert.NotNil(t, gm2.GetDB())
	gm2.ClearAllData()
}

func Test_gormMysql_Create(t *testing.T) {
	conf := GetDefaultDbConfig("hahaha_test")

	gm := newGormMysql(conf,  cusLog, true)
	assert.NotNil(t, gm)
	gm.CreateDB()
	gm.GetUtilDB()
	defer gm.DropDB()

	type Test struct {
		Data string `gorm:"data"`
	}
	gm2 := newGormMysql(conf,  cusLog, false)
	gm2.GetDB().AutoMigrate(&Test{})
	assert.NotNil(t, gm2)
	assert.NoError(t, gm2.GetDB().Create(&Test{}).Error)
	assert.Error(t, gm2.GetDB().Create([]*Test{}).Error)
	assert.NoError(t, gm2.GetDB().Create([]*Test{{Data: "1"}, {Data: "2"}, {Data: "3"}}).Error)
	assert.NoError(t, gm2.GetDB().CreateInBatches([]*Test{{Data: "1"}, {Data: "2"}, {Data: "3"},{Data: "1"}, {Data: "2"}, {Data: "3"}}, 2).Error)
	gm2.ClearAllData()
}
