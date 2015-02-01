package main

import (
	_ "./component/session" // 运行Session的gc
	"./controller"
	"net/http"
)

// 路由规则
func route() {
	controller.RouteHome()
	controller.RouteAdmin()
}

// 静态文件路由
func asset() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

// 主函数
func main() {
	asset()
	route()
	http.ListenAndServe(":8080", nil)
}
