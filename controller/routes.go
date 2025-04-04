package controller

import (
	"chat/service"
)

func (ctl *Controller) Routes() {
	authService := service.AuthService{DB: ctl.DB}
	chatService := service.ChatService{DB: ctl.DB}
	authController := AuthController{Service: authService}
	chatController := ChatController{Service: chatService, AuthService: authService}
	ctl.Gin.Use(CORSMiddleware())
	ctl.Gin.POST("/login", authController.Login)
	ctl.Gin.POST("/register", authController.Register)
	ctl.Gin.GET("/chat", chatController.ServeWebSocket)
	chatGroup := ctl.Gin.Group("/chat").Use(AuthMiddleware(ctl.DB), CORSMiddleware())
	{
		chatGroup.GET("/:receiverID", chatController.GetChatHistory)
		chatGroup.GET("/online/:userID", chatController.IsUserOnline)
	}

}
