package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DataBase struct {
	engine *sql.DB
}

var DB *sqlx.DB

func InitSqlx(driver, dataSourceName string) error {
	var err error
	DB, err = sqlx.Connect(driver, dataSourceName)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(200) //设置最大连接数
	return nil
}

func Execute(sqlCmd string, args ...interface{}) (*sql.Result, error) {
	result, err := DB.Exec(sqlCmd, args...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func QueryNormal(sqlCmd string, args ...interface{}) (*sql.Rows, error) {
	if args == nil {
		return DB.Query(sqlCmd)
	}
	return DB.Query(sqlCmd, args...)
}
