package mysql

type DbConfig struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DbName       string
	MaxIdleConns int
	MaxOpenConns int
	DbCharset    string
}

func GetDefaultDbConfig(dbName string) *DbConfig {
	return &DbConfig{
		Username:     "zebreay",
		Password:     "516504610",
		Host:         "localhost",
		Port:         "3306",
		DbName:       dbName,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		DbCharset:    "utf8mb4",
	}
}