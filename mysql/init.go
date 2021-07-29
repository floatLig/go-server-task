package mysql

import (
	"database/sql"
	"fmt"
	"time"

	config "shopee.com/zeliang-entry-task/const"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitMySql() error {
	var err error

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", config.MysqlUser, config.MysqlPassword, config.MysqlHost, config.MysqlPort, config.MysqlDatabase)
	DB, err = sql.Open(config.MysqlDriver, dataSourceName)
	if err != nil {
		return err
	}
	DB.SetMaxIdleConns(config.MysqlMaxIdleConns)
	DB.SetMaxOpenConns(config.MysqlMaxOpenConns)
	DB.SetConnMaxLifetime(time.Second)

	if err = DB.Ping(); err != nil {
		return err
	}

	go SelectAllUserToChan()
	return nil
}
