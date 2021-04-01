package mysql

import (
	"fmt"
	"github.com/dipperin/go-ms-toolkit/qyenv"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type CustomerLog func(msg string, items ...interface{})

var logConfig = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
	SlowThreshold: 1 * time.Second,
	LogLevel:      logger.Warn,
	Colorful:      true,
})

type gorm2Mysql struct {
	dbConfig    *DbConfig
	replicas    []*DbConfig
	db          *gorm.DB
	utilDB      *gorm.DB
	customerLog CustomerLog
}

func (gm *gorm2Mysql) CreateDB() {
	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + gm.dbConfig.DbName + " DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;"

	err := gm.utilDB.Exec(createDbSQL).Error
	if err != nil {
		fmt.Println("创建失败：" + err.Error() + " sql:" + createDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库创建成功")
}

func (gm *gorm2Mysql) DropDB() {
	dropDbSQL := "DROP DATABASE IF EXISTS " + gm.dbConfig.DbName + ";"

	err := gm.utilDB.Exec(dropDbSQL).Error
	if err != nil {
		fmt.Println("删除失败：" + err.Error() + " sql:" + dropDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库删除成功")
}

func (gm *gorm2Mysql) GetDB() *gorm.DB {
	return gm.db
}

func (gm *gorm2Mysql) GetUtilDB() *gorm.DB {
	gm.customerLog("init db connection: ", "db_host", gm.dbConfig.Host, "db_name", gm.dbConfig.DbName, "user", gm.dbConfig.Username)

	mysqlDialer := mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName, gm.dbConfig.DbCharset))
	gormLogLevel := logger.Error
	// 如果不是生产数据库则打开详细日志
	if substr(gm.dbConfig.DbName, len(gm.dbConfig.DbName)-4, 4) != "prod" {
		gormLogLevel = logger.Info
	}
	openedDb, err := gorm.Open(mysqlDialer, &gorm.Config{
		Logger: logConfig.LogMode(gormLogLevel),
	})
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	dbPool, err := openedDb.DB()
	if err != nil {
		panic("获取数据库连接池出错：" + err.Error())
	}
	dbPool.SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	dbPool.SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	dbPool.SetConnMaxLifetime(time.Hour * 1)

	return openedDb
}

func (gm *gorm2Mysql) ClearAllData() {
	if qyenv.IsUnitTestEnv() && strings.Contains(gm.dbConfig.DbName, "test") {
		tmpDb := gm.db
		if tmpDb == nil {
			panic("尚未初始化数据库, 清空数据库失败")
		}
		if rs, err := tmpDb.Raw("show tables;").Rows(); err == nil {
			var tName string
			for rs.Next() {
				if err := rs.Scan(&tName); err != nil || tName == "" {
					fmt.Println("表名获取失败", err, tName)
					panic("表名获取失败")
				}
				if err := tmpDb.Exec(fmt.Sprintf("delete from %s", tName)).Error; err != nil {
					panic("清空表数据失败:" + err.Error())
				}
			}
		} else {
			panic("表名列表获取失败：" + err.Error())
		}
	} else {
		panic("非法操作！在非测试环境下调用了清空所有数据的方法")
	}
}

func newGormMysql(dbConfig *DbConfig, cusLog CustomerLog, forUtil bool, replicas ...*DbConfig) *gorm2Mysql {
	gm := &gorm2Mysql{dbConfig: dbConfig, replicas: replicas, customerLog: cusLog}

	if forUtil {
		gm.initCdDb()
		return gm
	}

	// init db
	gm.initGormDB()

	return gm
}

func (gm *gorm2Mysql) initGormDB() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	gm.customerLog("init db connection: ", "db_host", gm.dbConfig.Host, "db_name", gm.dbConfig.DbName, "user", gm.dbConfig.Username)

	mysqlDialer := mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName, gm.dbConfig.DbCharset))
	gormLogLevel := logger.Silent
	// 如果不是生产数据库则打开详细日志
	if substr(gm.dbConfig.DbName, len(gm.dbConfig.DbName)-4, 4) != "prod" {
		gormLogLevel = logger.Info
	}
	openedDb, err := gorm.Open(mysqlDialer, &gorm.Config{
		Logger: logConfig.LogMode(gormLogLevel),
	})
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	dbPool, err := openedDb.DB()
	if err != nil {
		panic("获取数据库连接池出错：" + err.Error())
	}
	dbPool.SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	dbPool.SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	dbPool.SetConnMaxLifetime(time.Hour * 1)

	gm.db = openedDb

	// 读分离 - 丛库
	gm.addReplicas()
}

func (gm *gorm2Mysql) addReplicas() {
	if len(gm.replicas) <= 0 {
		return
	}
	var replicaDialers = make([]gorm.Dialector, 0, len(gm.replicas))
	for _, config := range gm.replicas {
		if config == nil {
			return
		}
		gm.customerLog("init db connection: ", "db_host", gm.dbConfig.Host, "db_name", gm.dbConfig.DbName, "user", gm.dbConfig.Username)

		replicaDialers = append(replicaDialers, mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName, gm.dbConfig.DbCharset)))
	}
	err := gm.db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas: replicaDialers,
			Policy:   dbresolver.RandomPolicy{}}). // 随机负载均衡
			SetMaxIdleConns(gm.dbConfig.MaxIdleConns).
			SetMaxOpenConns(gm.dbConfig.MaxOpenConns).
			SetConnMaxLifetime(time.Hour * 1))
	if err != nil {
		panic("初始化丛库出错: " + err.Error())
	}
}

func (gm *gorm2Mysql) initCdDb() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	cStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, "information_schema", gm.dbConfig.DbCharset)
	openedDb, err := gorm.Open(mysql.Open(cStr), &gorm.Config{Logger: logConfig})
	if err != nil {
		fmt.Println(cStr)
		panic("连接数据库出错:" + err.Error())
	}

	gm.utilDB = openedDb
}

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
