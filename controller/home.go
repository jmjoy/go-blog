package controller

import (
	"../model"
	"fmt"
	"net/http"
)

type HomeController struct {
	Controller
}

var (
	homeModel = new(model.HomeModel)
)

func RouteHome() {
	c := new(HomeController)
	http.HandleFunc("/", c.Index)
	http.HandleFunc("/home", c.Index)
	http.HandleFunc("/home/", c.Index)
	http.HandleFunc("/home/index", c.Index)
}

// 首页
func (this *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := homeModel.ShowAllArticles()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	data := map[string]interface{}{
		"articles": articles,
	}
	this.render(w, "home/index", data)
}
