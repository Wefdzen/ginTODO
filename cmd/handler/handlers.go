package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"wefdzen/cmd/users"
	"wefdzen/pkg/postgres"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil) //parse html file
	}
}
func LoginPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		login := c.PostForm("login")       // Получение значения поля "login"
		password := c.PostForm("password") // Получение значения поля "password"
		//check login with hash bcrypt compare password
		if postgres.CheckDataForLogin(login, password) {
			// потомо сделать доступ к поста токо с jwt tokens
			//TODO generate jwt token
			//add to set-cookie
			c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/postes")
		} else {
			c.String(http.StatusUnauthorized, "loh!")
		}

	}
}

func Registration() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "registration.html", nil) //parse html file
	}
}
func RegistrationPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		login := c.PostForm("login")
		email := c.PostForm("email")
		password := c.PostForm("password")
		//create hash
		HashPass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
		var newUser users.User = users.User{
			Login:    login,
			Email:    email,
			Password: string(HashPass),
		}

		//send to database with hash bcrypt
		if success := postgres.RegistrationUser(newUser); success {
			c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/login")
		} else {
			c.String(401, "login занят")
		}

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

func WatchPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		idCheck, err := strconv.Atoi(id)
		if err != nil || idCheck < 0 {
			fmt.Println(err.Error(), " or id < 0")
		} else { // all good
			sourceOfPost := postgres.WatchPostByID(idCheck)
			c.JSON(http.StatusOK, gin.H{sourceOfPost.Title: sourceOfPost.Post})
		}
	}
}

func EditingPost() gin.HandlerFunc {
	newTitle := "for test"
	newText := "i love anime yopta"
	return func(c *gin.Context) {
		id := c.Param("id")
		idCheck, err := strconv.Atoi(id)
		if err != nil || idCheck < 0 {
			fmt.Println(err.Error(), " or id < 0")
		} else { // all good
			postgres.FullEditPostByID(idCheck, newTitle, newText)
		}
	}
}
