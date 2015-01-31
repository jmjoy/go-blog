package controller

import (
	"net/http"
)

type HomeController struct {
	Controller
}

func RouteHome() {
	c := new(HomeController)
	http.HandleFunc("/", c.Index)
	http.HandleFunc("/home", c.Index)
	http.HandleFunc("/home/index", c.Index)
}

func (this *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	this.render(w, "home/index", nil)
}
