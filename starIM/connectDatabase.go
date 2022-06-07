package starIM

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func Conn() error { //连接
	dsn := "root:root@tcp(127.0.0.1:3306)/star?parseTime=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	return err
}
