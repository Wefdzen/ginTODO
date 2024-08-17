package main

import (
	"log"

	"wefdzen/cmd/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	route.LoadHTMLGlob("static/html/*")

	route.GET("/login", handler.Login())
	route.POST("/login", handler.LoginPost())

	route.GET("/mainpage", handler.MainPage()) // после mainpage не надо / он будет отправлять на mainpage/ ...
	route.POST("/mainpage", handler.MainPagePost())

	log.Fatal(route.Run())

}
