package model

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var MiscLogsDB *sql.DB

func InitConn(url string) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(300)
	// 打开数据库连接的最大数量
	db.SetMaxOpenConns(500)
	// 连接可复用的最大时间
	db.SetConnMaxLifetime(time.Second * 30)

	MiscLogsDB = db
}
