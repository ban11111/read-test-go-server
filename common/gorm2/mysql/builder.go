package mysql

func MakeDBUtil(dbConfig *DbConfig, log CustomerLog) DBUtil {
	return newGormMysql(dbConfig, log, true)
}

// replicas 读写分离的丛库配置
func MakeDB(dbConfig *DbConfig, log CustomerLog, replicas ...*DbConfig) DB {
	return newGormMysql(dbConfig, log, false, replicas...)
}
