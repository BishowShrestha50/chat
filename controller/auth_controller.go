package controller

import (
	"chat/utils"
	"github.com/gin-gonic/gin"
	"net/http"

	"chat/model"
	"chat/service"
)

type AuthController struct {
	Service service.AuthService
}

func (u *AuthController) Register(c *gin.Context) {
	var registerData model.Credentials
	if err := c.ShouldBindJSON(&registerData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	err := registerData.Validate()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = u.Service.Register(registerData)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (u *AuthController) Login(c *gin.Context) {
	var loginData model.Credentials
	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	err := loginData.Validate()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := u.Service.Login(loginData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "invalid credentials")
		return
	}
	c.JSON(http.StatusOK, model.TokenResponse{Token: token})
}
