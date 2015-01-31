package model

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

// 公共模型类
type Model struct {
}

type dbOperateFunc func(db *sql.DB) error

// 数据库操作通用方法
func (this *Model) dbOperate(op dbOperateFunc) error {
	db, err := sql.Open("sqlite3", "./db/blog.sq3")
	if err != nil {
		return errors.New("数据库出问题了")
	}
	defer db.Close()
	return op(db)
}
