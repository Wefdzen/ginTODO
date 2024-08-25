package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wefdzen/cmd/users"
	"wefdzen/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil) //parse html file
	}
}
func LoginPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonInput users.User
		//TODO //add postform
		// Парсинг JSON из тела запроса
		if err := c.BindJSON(&jsonInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		login := jsonInput.Login
		email := jsonInput.Email
		password := jsonInput.Password

		//check login with hash bcrypt compare password
		if postgres.CheckDataForLogin(login, email, password) {
			// потом сделать доступ к поста токо с jwt tokens
			//TODO generate jwt token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"admin": "false",
				"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(), //30 day
			})
			tokenString, err := token.SignedString([]byte("secret-key")) // secret key os.getenv
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Failed to create token"})
			}
			//add to set-cookie
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
			c.JSON(http.StatusOK, gin.H{})
			//redirect to main page
			//c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/postes")
		} else {
			c.String(http.StatusUnauthorized, "loh!")
		}

	}
}

// type MyClaims struct {
// 	jwt.RegisteredClaims
// 	Admin bool  `json:"admin"`
// 	Exp   int64 `json:"exp"`
// }

func CheckJWToken(tokenString string) bool {
	//token, err
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret-key"), nil
	})
	if err != nil {
		return false
	}

	//parse claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) >= claims["exp"].(float64) { // проверка exp токена
			return false
		} else {
			return true
		}
	}
	return false
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
		token, err := c.Cookie("Authorization") //take cookie
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't take cookie auth",
			})
			return
		}
		if proverka := CheckJWToken(token); proverka {
			c.JSON(http.StatusOK, gin.H{
				"status": "good",
			})
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

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
