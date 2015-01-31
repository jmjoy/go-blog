package controller

import (
	"../model"
	"net/http"
)

type AdminController struct {
	Controller
}

var (
	adminModel = new(model.AdminModel)
)

// 管理员控制器相关的路由规则
func RouteAdmin() {
	c := new(AdminController)
	http.HandleFunc("/admin", c.Index)
	http.HandleFunc("/admin/index", c.Index)
	http.HandleFunc("/admin/handle-sign-in", c.HandleSignIn)
	http.HandleFunc("/admin/list-article", c.ListArticle)
	http.HandleFunc("/admin/add-article", c.AddArticle)
}

// 登陆
func (this *AdminController) Index(w http.ResponseWriter, r *http.Request) {
	// 获取登陆失败的闪存数据
	errMsg := this.getFlashCookie(w, r, "errMsg")
	errName := this.getFlashCookie(w, r, "errName")
	data := map[string]string{
		"errMsg":  errMsg,
		"errName": errName,
	}
	// 渲染模板
	this.render(w, "admin/index", data)
}

// 处理登陆请求
func (this *AdminController) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	// 获取管理员账号和密码
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")
	// 检验账号和密码正确性
	err := adminModel.HandleSignIn(name, password)
	// 校验不通过
	if err != nil {
		// 用闪存cookie保存错误信息
		this.setFlashCookie(w, "errMsg", err.Error())
		this.setFlashCookie(w, "errName", name)
		// 跳回登陆页面
		http.Redirect(w, r, "/admin/index", 302)
		return
	}
	// 检验通过，跳到文章显示页面
	http.Redirect(w, r, "/admin/list-article", 302)
}

func (this *AdminController) ListArticle(w http.ResponseWriter, r *http.Request) {
	this.render(w, "admin/list-article", nil)
}

func (this *AdminController) AddArticle(w http.ResponseWriter, r *http.Request) {
	this.render(w, "admin/add-article", nil)
}
