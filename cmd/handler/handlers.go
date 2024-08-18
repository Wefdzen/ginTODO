package handler

import (
	"fmt"
	"net/http"
	"strconv"

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

func CreateNewPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title_of_post")
		text := c.PostForm("source_of_post")
		postgres.InsertNewPost(title, text)
		c.String(http.StatusOK, "Данные отправлены успешно!")
		c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/postes")
	}
}

func GetAllPostes() gin.HandlerFunc {
	return func(c *gin.Context) {
		postesFromDb := postgres.GetAllPost()
		//fmt.Println(postesFromDb)
		for _, v := range postesFromDb {
			fmt.Println("name of title:", v.Title, " text:", v.Post)
		}
		c.JSON(http.StatusOK, gin.H{"Lol": "postes"})
	}
}

// TODO реализовать функцию delete post in postgres
func DeletePostes() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idCheck, err := strconv.Atoi(id)
		if err != nil || idCheck < 0 {
			fmt.Println(err.Error(), " or id < 0")
		} else { // all good
			postgres.DeletePostByID(idCheck)
		}
	}
}

// TODO реализовать функцию watchpost in postgres
func WatchPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idCheck, err := strconv.Atoi(id)
		if err != nil || idCheck < 0 {
			fmt.Println(err.Error(), " or id < 0")
		} else { // all good
			postgres.WatchPostByID(idCheck)
		}
	}
}

// TODO реализовать функцию edit post in postgres
func EditingPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idCheck, err := strconv.Atoi(id)
		if err != nil || idCheck < 0 {
			fmt.Println(err.Error(), " or id < 0")
		} else { // all good
			postgres.EditPostByID(idCheck)
		}
	}
}
