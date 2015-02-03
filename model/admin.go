package model

import (
	"../helper/mytime"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 可以遍历结果集的接口
type RowScanner interface {
	Scan(dest ...interface{}) error
}

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

// 统计文章数量
func (this *AdminModel) CountArticle() (count int, err error) {
	err = this.dbOperate(func(db *sql.DB) error {
		return db.QueryRow("select count(*) from article").Scan(&count)
	})
	return
}

// 列出文章
func (this *AdminModel) ListArticle(page, rowList string) (result []map[string]string, err error) {
	// 转换成数字
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return
	}
	if pageInt <= 0 {
		pageInt = 1
	}
	rowListInt, err := strconv.Atoi(rowList)
	if err != nil {
		return
	}
	// 执行数据库查询
	err = this.dbOperate(func(db *sql.DB) error {
		// 执行查询，分页
		rows, err := db.Query("select id, title, ctime, mtime from article limit ?, ?", (pageInt-1)*rowListInt, rowListInt)
		if err != nil {
			return err
		}
		defer rows.Close()
		// 将查询结果放进结果数组
		var i int
		for rows.Next() {
			row, err := this.pushSingleArticleWithoutContent(rows)
			if err != nil {
				return err
			}
			result = append(result, row)
			i++
		}
		return nil
	})
	return
}

// 展示文章的内容
func (this *AdminModel) ShowArticle(id string) (article map[string]string, err error) {
	// 将文章id转成数字
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	// 根据文章的id查找文章
	err = this.dbOperate(func(db *sql.DB) error {
		rows := db.QueryRow("select * from article where id = ? limit 1", idInt)
		article, err = this.pushSingleArticle(rows)
		return err
	})
	return
}

// 处理Upsert文章的请求
func (this *AdminModel) HandleUpsertArticle(id, title, content string) error {
	// 因为是管理员上传的，所以不需要校验了
	return this.dbOperate(func(db *sql.DB) error {
		// 当前时间
		now := time.Now().Unix()
		if id == "" {
			// Id为空，说明是插入操作
			str := "insert into article values(null, ?, ?, ?, ?)"
			_, err := db.Exec(str, title, content, now, now)
			return err
		}
		// Id不为空，说明是更新操作
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		str := "update article set title = ?, content = ?, mtime = ? where id = ?"
		_, err = db.Exec(str, title, content, now, idInt)
		return err
	})
}

// 根据id删除文章
func (this *AdminModel) DelArticle(id string) error {
	// 将id转成int型
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	// 调用数据库方法
	return this.dbOperate(func(db *sql.DB) error {
		// 执行删除操作
		_, err := db.Exec("delete from article where id = ?", idInt)
		return err
	})
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

// 从结果集获取一个文章的数组
func (this *AdminModel) pushSingleArticleWithoutContent(rows RowScanner) (article map[string]string, err error) {
	var id, ctime, mtime int
	var title string
	if err = rows.Scan(&id, &title, &ctime, &mtime); err != nil {
		return
	}
	// 将获取的数据全转成相应的字符串，放进结果map中
	article = map[string]string{
		"id":    strconv.Itoa(id),
		"title": title,
		"ctime": mytime.GetDateTime(int64(ctime)),
		"mtime": mytime.GetDateTime(int64(mtime)),
	}
	return
}

// 从结果集获取一个文章的数组
func (this *AdminModel) pushSingleArticle(rows RowScanner) (article map[string]string, err error) {
	var id, ctime, mtime int
	var title, content string
	if err = rows.Scan(&id, &title, &content, &ctime, &mtime); err != nil {
		return
	}
	// 将获取的数据全转成相应的字符串，放进结果map中
	article = map[string]string{
		"id":      strconv.Itoa(id),
		"title":   title,
		"content": content,
		"ctime":   mytime.GetDateTime(int64(ctime)),
		"mtime":   mytime.GetDateTime(int64(mtime)),
	}
	return
}
