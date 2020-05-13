package db

import (
	"EX_okexquant/config"
	"EX_okexquant/mylog"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	mysqlConn *sqlx.DB
)

func InitMysqlCli() {
	var err error

	user := config.Config.Mysql.User
	password := config.Config.Mysql.Password
	address := config.Config.Mysql.Address
	database := config.Config.Mysql.DataBase
	maxOpenConn := config.Config.Mysql.MaxOpenConn
	maxIdleConn := config.Config.Mysql.MaxIdleConn

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%v?charset=utf8&multiStatements=true",
		user, password, address, database)
	mysqlConn, err = sqlx.Open("mysql", dataSourceName)
	if err != nil {
		mylog.Logger.Fatal().Msgf("[InitMysqlCli] open mysql connection failed, err=%v, dataSource=%v", err, dataSourceName)
	}
	mysqlConn.SetMaxIdleConns(maxIdleConn)
	mysqlConn.SetMaxOpenConns(maxOpenConn)

	fmt.Println("[InitMysql] mysql succeed.")
}

func CloseMysqlCli() {
	mysqlConn.Close()
}
