package main

import (
	"log"

	"wefdzen/cmd/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("static/html/*")

	v1 := router.Group("/")
	{
		v1.GET("/login", handler.Login())
		v1.POST("/login", handler.LoginPost())
		v1.GET("/registration", handler.Registration())
		v1.POST("/registration", handler.RegistrationPost())
	}

	v2 := router.Group("/")
	{
		v2.GET("/postes", handler.GetAllPostes())
		v2.GET("/postes/:id", handler.WatchPost())
		v2.DELETE("/postes/:id", handler.DeletePostes())
		v2.PUT("/postes/:id", handler.EditingPost())
		v2.GET("/createpost", handler.MainPage()) // после mainpage не надо / он будет отправлять на mainpage/ ...
		v2.POST("/createpost", handler.CreateNewPost())

	}

	log.Fatal(router.Run())

}
