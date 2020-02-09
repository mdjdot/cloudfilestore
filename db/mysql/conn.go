package mysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql驱动
)

// FileDB 文件的数据库
var FileDB *sql.DB

func init() {
	fileDB, err := sql.Open("mysql", "root:dmtest@tcp(127.0.0.1:3301)/fileserver?charset=utf8")
	if err != nil {
		log.Fatalln("file to connect to mysql,err: ", err.Error())
	}
	fileDB.SetMaxOpenConns(30)
	FileDB = fileDB

	err = FileDB.Ping()
	if err != nil {
		log.Fatalln("file to connect to mysql,err: ", err.Error())
	}
}
