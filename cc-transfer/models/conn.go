package models

import (
	"cc-transfer/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

var db *sqlx.DB

func init() {
	//用sql.Open的方式创建的对象是长连接的，协程安全的
	db, _ = sqlx.Open("mysql", config.Username+ ":" + config.Pwd + "@tcp(" + config.Host + ")/" + config.DBName + "?charset=utf8")
	err := db.Ping()
	if err != nil {
		fmt.Println("Fail to connect to MySQl_master, error: " + err.Error())
		os.Exit(1)
	} else {
		fmt.Println("MySQL_master is connected successfully!")
	}
}

//返回数据库对象
func DBConn() *sqlx.DB {
	return db
}