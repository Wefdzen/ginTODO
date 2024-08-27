package handler

import (
	"wefdzen/cmd/users"

	"github.com/gin-gonic/gin"
)

// TODO form from html don't work (
func JsonOrFormForLoginOrReg(c *gin.Context) (login, email, password string) {
	var jsonInput users.User

	// Парсинг JSON из тела запроса
	if err := c.BindJSON(&jsonInput); err == nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})

		login = jsonInput.Login
		email = jsonInput.Email
		password = jsonInput.Password
		return
	} else { // если json не соответствует json тогда user дал form
		login = c.PostForm("login")
		email = c.PostForm("email")
		password = c.PostForm("password")
		return
	}

}
