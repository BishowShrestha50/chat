package controller

import (
	"chat/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type Controller struct {
	DB  *gorm.DB
	Gin *gin.Engine
}

func NewController() *Controller {
	ctl := &Controller{}
	ctl.DB = NewDB()
	ctl.Gin = gin.Default()
	ctl.Routes()
	return ctl
}

func (ctl *Controller) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logrus.Infof("server start on port %s ", port)
	return ctl.Gin.Run(":" + port)
}

func NewDB() *gorm.DB {
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	db.AutoMigrate(&model.User{}, &model.Chat{}, model.ChatMessage{})
	return db

}
