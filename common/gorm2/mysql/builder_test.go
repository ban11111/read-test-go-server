package mysql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeDBUtil(t *testing.T) {
	conf := GetDefaultDbConfig("hahaha_test")

	assert.NotNil(t, MakeDBUtil(conf, cusLog))
}

func TestMakeDB(t *testing.T) {
	conf := GetDefaultDbConfig("hahaha_test")

	utilDB := MakeDBUtil(conf, cusLog)
	assert.NotNil(t, utilDB)

	utilDB.CreateDB()

	db := MakeDB(conf, cusLog)
	assert.NotNil(t, db)

	utilDB.DropDB()
}
