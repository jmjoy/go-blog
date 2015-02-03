package controller

import (
	"../model"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
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
	http.HandleFunc("/admin/upsert-article", c.UpsertArticle)
	http.HandleFunc("/admin/handle-upsert-article", c.HandleUpsertArticle)
	http.HandleFunc("/admin/show-article", c.ShowArticle)
	http.HandleFunc("/admin/del-article", c.DelArticle)
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
	if _, had := this.hadSignIn(w, r); !had {
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
	// 获取分页参数
	querys, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	page := "1"
	if len(querys["page"]) != 0 {
		page = querys["page"][0]
	}
	rowList := 5
	// 获取文章列表
	res, err := adminModel.ListArticle(page, strconv.Itoa(rowList))
	// 出现罕见数据库查询错误
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	// 获取文章总数
	count, err := adminModel.CountArticle()
	// 出现罕见数据库查询错误
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	// 获取页数、上一页、下一页
	pageCount := int(math.Ceil(float64(count) / float64(rowList)))
	pageInt, _ := strconv.Atoi(page)
	if pageInt <= 0 {
		pageInt = 1
	}
	prePage := pageInt - 1
	nextPage := pageInt + 1
	if nextPage >= pageCount {
		nextPage = pageCount
	}
	preArr := make([]int, 0, 5)
	for i := pageInt - 5; i < pageInt; i++ {
		if i < 1 {
			continue
		}
		preArr = append(preArr, i)
	}
	nextArr := make([]int, 0, 5)
	for i := pageInt + 1; i < pageInt+5; i++ {
		if i > pageCount {
			continue
		}
		nextArr = append(nextArr, i)
	}
	// 获取数据和渲染模板
	data := map[string]interface{}{
		"adminName": adminName,
		"resArr":    res,
		"page":      pageInt,
		"prePage":   prePage,
		"nextPage":  nextPage,
		"pageCount": pageCount,
		"preArr":    preArr,
		"nextArr":   nextArr,
	}
	this.render(w, "admin/manage", data)
}

// 增加或者修改文章页面
func (this *AdminController) UpsertArticle(w http.ResponseWriter, r *http.Request) {
	// 检测有没有登陆
	if _, had := this.hadSignIn(w, r); !had {
		this.notSignIn(w, r)
		return
	}
	// 获取id的值
	querys, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	// 要传递到模板的数据
	var data map[string]string
	if len(querys["id"]) <= 0 {
		// 如果id的值不存在，说明是增加操作
		data = map[string]string{
			"id":      "",
			"title":   "",
			"content": "",
		}
	} else {
		// id的值存在，说明是修改操作
		id := querys["id"][0]
		article, err := adminModel.ShowArticle(id)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		data = article
	}
	this.render(w, "admin/upsert-article", data)
}

// 处理文章增加或者修改
func (this *AdminController) HandleUpsertArticle(w http.ResponseWriter, r *http.Request) {
	// 检测有没有登陆
	if _, had := this.hadSignIn(w, r); !had {
		this.notSignIn(w, r)
		return
	}
	// 获取Form表单的数据
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")
	// 插入或更新文章
	err := adminModel.HandleUpsertArticle(id, title, content)
	// 出现罕见的错误
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	// 成功
	http.Redirect(w, r, "/admin/manage", 302)
}

// 展示文章的内容
func (this *AdminController) ShowArticle(w http.ResponseWriter, r *http.Request) {
	// 检测有没有登陆
	if _, had := this.hadSignIn(w, r); !had {
		this.notSignIn(w, r)
		return
	}
	// 获取GET参数中id的值
	querys, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(querys["id"]) <= 0 {
		fmt.Fprint(w, err.Error())
		return
	}
	id := querys["id"][0]
	// 数据库查询文章
	article, err := adminModel.ShowArticle(id)
	// 罕见错误
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	this.render(w, "admin/show-article", article)
}

// 根据id删除文章
func (this *AdminController) DelArticle(w http.ResponseWriter, r *http.Request) {
	// 检测有没有登陆
	if _, had := this.hadSignIn(w, r); !had {
		this.notSignIn(w, r)
		return
	}
	// 获取输入的id参数
	id := r.FormValue("id")
	// 数据库删除
	err := adminModel.DelArticle(id)
	if err != nil {
		this.simpleJsonReturn(w, 400, err.Error())
		return
	}
	this.simpleJsonReturn(w, 200, "")
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
