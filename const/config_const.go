package _const

import "time"

const (
	TcpNetWork = "tcp"
	TcpHost    = "127.0.0.1"
	TcpPort    = "3030"
	TcpTimeout = time.Second * 30

	MysqlDriver       = "mysql"
	MysqlUser         = "root"
	MysqlPassword     = "root"
	MysqlHost         = "127.0.0.1"
	MysqlPort         = "3306"
	MysqlDatabase     = "entryTask"
	MysqlTable        = "user"
	MysqlMaxIdleConns = 100
	MysqlMaxOpenConns = 100

	LoginFile  = "HTTPServer/template/login.html"
	LoginToken = "token"
	UserFile   = "HTTPServer/template/user.html"

	TokenSuffix   = "_token"
	UserKeyPrefix = "key:"

	HTTPConnectionPoolSize = 100
	HttpHost               = "127.0.0.1"
	HttpPort               = "8888"

	ImgStorePath = "./img"

	RedisHost     = "localhost"
	RedisPort     = "6379"
	RedisPassword = ""
	RedisDb       = 0
	RedisPoolSize = 60
	RedisExpire   = time.Minute * 30
)
