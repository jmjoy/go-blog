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
	http.HandleFunc("/admin/", c.Index)
	http.HandleFunc("/admin/index", c.Index)
	http.HandleFunc("/admin/handle-sign-in", c.HandleSignIn)
	http.HandleFunc("/admin/handle-sign-out", c.HandleSignOut)
	http.HandleFunc("/admin/manage", c.Manage)
	http.HandleFunc("/admin/add-article", c.AddArticle)
}

// 登陆
func (this *AdminController) Index(w http.ResponseWriter, r *http.Request) {
	// 如果已经登陆了跳转到后台管理界面
	if _, had := this.hadSignIn(w, r); had {
		http.Redirect(w, r, "/admin/manage", 302)
		return
	}
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
	// 生成一个SessionToken
	token := this.getSessionToken(w, r)
	// 检验账号和密码正确性
	err := adminModel.HandleSignIn(token, name, password)
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
	http.Redirect(w, r, "/admin/manage", 302)
}

// 处理注销请求
func (this *AdminController) HandleSignOut(w http.ResponseWriter, r *http.Request) {
	// 检测有没有登陆
	_, had := this.hadSignIn(w, r)
	if !had {
		this.notSignIn(w, r)
		return
	}
	// 登陆了，让他注销
	token := this.getSessionToken(w, r)
	adminModel.HadleSignOut(token)
	http.Redirect(w, r, "/admin/index", 302)
}

// 管理页面首页，列出文章
func (this *AdminController) Manage(w http.ResponseWriter, r *http.Request) {
	// 检测有没有登陆
	adminName, had := this.hadSignIn(w, r)
	if !had {
		this.notSignIn(w, r)
		return
	}
	data := map[string]string{
		"adminName": adminName,
	}
	this.render(w, "admin/manage", data)
}

func (this *AdminController) AddArticle(w http.ResponseWriter, r *http.Request) {
	this.render(w, "admin/add-article", nil)
}

// 检查有没有登陆
func (this *AdminController) hadSignIn(w http.ResponseWriter, r *http.Request) (adminName string, had bool) {
	token := this.getSessionToken(w, r)
	adminName, had = adminModel.HadSignIn(token)
	return
}

// 没有登陆或者登陆超时而进行管理
func (this *AdminController) notSignIn(w http.ResponseWriter, r *http.Request) {
	this.setFlashCookie(w, "errMsg", "登陆超时")
	http.Redirect(w, r, "/admin/index", 302)
}
