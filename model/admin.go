package model

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 管理员模型
type AdminModel struct {
	Model
}

// 处理登陆请求
func (this *AdminModel) HandleSignIn(token, name, password string) error {
	// 检验账号或密码是否为空
	if strings.Trim(name, " ") == "" {
		return errors.New("账号不能为空")
	}
	if strings.Trim(password, " ") == "" {
		return errors.New("密码不能为空")
	}
	// 执行数据库查询
	return this.dbOperate(func(db *sql.DB) error {
		// 根据账号获取正确密码
		str := "select id, password from admin where name = ?"
		var id int
		var corrPasswd string
		err := db.QueryRow(str, name).Scan(&id, &corrPasswd)
		// 判断账号密码是否正确
		switch {
		// 没找到数据库记录
		case err == sql.ErrNoRows:
			return errors.New("账号或者密码不正确")
		// 数据库查询发生错误
		case err != nil:
			return errors.New("数据库出问题了")
		// 查询成功，有记录
		default:
			// 判断传入密码和正确的密码是否一致
			if res := this.validatePassword(password, corrPasswd); !res {
				return errors.New("账号或者密码不正确")
			}
			// 检验通过
			this.Sess.Set(token, "AdminId", strconv.Itoa(id))
			this.Sess.Set(token, "AdminName", name)
			return nil
		}
	})
}

// 处理注销请求
func (this *AdminModel) HadleSignOut(token string) {
	this.Sess.Drop(token)
}

// 检查有没有登陆
func (this *AdminModel) HadSignIn(token string) (adminName string, had bool) {
	adminName = this.Sess.Get(token, "AdminName")
	had = adminName != ""
	return
}

// 比较传入密码和正确的密码，一致的话返回true
func (this *AdminModel) validatePassword(password, corrPasswd string) bool {
	// 加密传入密码：sha1
	cryptedPasswd := fmt.Sprintf("%x", sha1.Sum([]byte(password)))
	// 比较
	if cryptedPasswd != corrPasswd {
		return false
	}
	return true
}
