package model

import (
	"../helper/mytime"
	"database/sql"
)

// 管理员模型
type HomeModel struct {
	Model
}

// 展示文章
func (this *HomeModel) ShowAllArticles() (articles []map[string]string, err error) {
	// 执行数据库查询
	err = this.dbOperate(func(db *sql.DB) error {
		rows, err := db.Query("select title, content, ctime, mtime from article order by mtime desc")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var title, content string
			var ctime, mtime int
			rows.Scan(&title, &content, &ctime, &mtime)
			articles = append(articles, map[string]string{
				"title":   title,
				"content": content,
				"ctime":   mytime.GetDateTime(int64(ctime)),
				"mtime":   mytime.GetDateTime(int64(mtime)),
			})
		}
		return nil
	})
	return
}
