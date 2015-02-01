package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

type Controller struct {
}

const (
	// 主模板路径
	LAYOUT_PATH = "./view/layout/main.html"
	// SessionCookie的键名
	SESSION_TOKEN_KEY = "GOBLOGSESS"
)

// 渲染模板，其中view的格式为“controller/action”
func (this *Controller) render(w http.ResponseWriter, view string, data interface{}) {
	strs := strings.Split(view, "/")
	baseName := strs[1] + ".html"
	filePath := "./view/" + view + ".html"
	t := template.Must(template.New(baseName).ParseFiles(filePath, LAYOUT_PATH))
	t.Execute(w, data)
}

// 获取闪存Cookie
func (this *Controller) getFlashCookie(w http.ResponseWriter, r *http.Request, name string) (value string) {
	// 获取cookie
	cookie, err := r.Cookie(name)
	if err != nil {
		value = ""
	} else {
		value, err = url.QueryUnescape(cookie.Value)
		if err != nil {
			value = ""
		}
	}
	// 删除闪存cookie
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		MaxAge: -1,
	})
	return
}

// 设置闪存Cookie
func (this *Controller) setFlashCookie(w http.ResponseWriter, name, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: url.QueryEscape(value),
	}
	http.SetCookie(w, cookie)
}

// 获取SessionCookie的值，如果不存在就创建！
func (this *Controller) getSessionToken(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(SESSION_TOKEN_KEY)
	// 不存在
	if err != nil {
		token := fmt.Sprintf("%x", time.Now().UnixNano())
		http.SetCookie(w, &http.Cookie{
			Name:  SESSION_TOKEN_KEY,
			Value: token,
		})
		return token
	}
	// 存在就直接返回
	return cookie.Value
}
