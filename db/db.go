package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	userName = "work"
	password = "mhxzkhl"
	ip       = "127.0.0.1"
	port     = "8021"
	dbName   = "exchange_db"
)

var db *sql.DB
var err error

func InitDB() bool {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println("path=", path)
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	db, err = sql.Open("mysql", path)
	if err != nil {
		fmt.Println("open mysql fail=", err.Error())

		return false
	}
	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil {
		fmt.Println("opon database fail")
		return false
	}
	fmt.Println("connnect success")
	return true
}

func ConnDb() error {
	if db == nil {
		flag := InitDB()
		if !flag {
			return errors.New("init db fail")
		}
	}
	return nil
}

func Insert(sql *string) (int64, error) {
	stmt, err := db.Prepare(*sql)
	if err == nil {
		res, err := stmt.Exec()
		if err == nil {
			return res.RowsAffected()
		} else {
			return 0, errors.New("exec insert fail")
		}
	} else {
		return 0, errors.New("prepare stmt error")
	}
}

func Select(sql *string) (*sql.Rows, error) {
	row, err := db.Query(*sql)
	if err == nil {
		return row, nil
	}
	return nil, errors.New("select error")
}

func Close() {
	if db != nil {
		db.Close()
	}
}
