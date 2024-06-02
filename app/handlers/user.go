package handlers

import (
	"user-profile-apis/controllers"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	controllers.UserRegister(c)
}

func UserLogin(c *gin.Context) {
	controllers.UserLogin(c)
}

func UpdateUser(c *gin.Context) {
	controllers.UpdateUser(c)
}

func DeleteUser(c *gin.Context) {
	controllers.DeleteUser(c)
}
