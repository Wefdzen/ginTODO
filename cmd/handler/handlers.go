package handler

import (
	"fmt"
	"net/http"

	"wefdzen/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil) //parse html file
	}
}

func LoginPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("login")                     // Получение значения поля "login"
		email := c.PostForm("password")                 // Получение значения поля "password"
		fmt.Printf("Имя: %s, Email: %s\n", name, email) // Обработка данных
		c.String(http.StatusOK, "Данные отправлены успешно!")
	}
}

func MainPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "createPost.html", nil) //parse html file
	}
}

func MainPagePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title_of_post")
		text := c.PostForm("source_of_post")
		postgres.InsertNewPost(title, text)
		c.String(http.StatusOK, "Данные отправлены успешно!")
	}
}
