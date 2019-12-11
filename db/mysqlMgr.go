package db

import (
	"fmt"
	"godb/config"
	// "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

)

// var mysqldb *sql.DB

var mysqldb *sqlx.DB
 
func InitMySqlDB() {
	var conf = config.GetConfig()
	if url, ok := conf.DB["mysql"]; ok {
		// sqldb, err := sql.Open("mysql", url)
		sqldb, err := sqlx.Open("mysql", url)
		if  err != nil {
			fmt.Printf("sql.Open(mysql) dsn(%s) error(%v)", url, err)
		}
		mysqldb = sqldb
	}
}

func MysqlDB() *sqlx.DB {
	return mysqldb
}

