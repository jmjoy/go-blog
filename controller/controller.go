package controller

import (
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

type Controller struct {
}

var (
	layoutPath = "./view/layout/main.html"
)

// 渲染模板，其中view的格式为“controller/action”
func (this *Controller) render(w http.ResponseWriter, view string, data interface{}) {
	strs := strings.Split(view, "/")
	baseName := strs[1] + ".html"
	filePath := "./view/" + view + ".html"
	t := template.Must(template.New(baseName).ParseFiles(filePath, layoutPath))
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
