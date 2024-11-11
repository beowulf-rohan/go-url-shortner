package main

import (
	"log"

	"github.com/beowulf-rohan/go-url-shortner/controller"
	"github.com/beowulf-rohan/go-url-shortner/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	envVariableList = []string{
		"DB_ADDR",
		"DB_PASS",
		"APP_PORT",
		"DOMAIN",
		"API_QUOTA",
		"API_PORT",
	}
)

func main() {

	log.Println("loading configs.....")
	config, err := utils.LoadEnvVaraibles(envVariableList)
	if err != nil {
		log.Fatal("Error loading configs from env file", err)
	}
	log.Println("configs loaded successfully....", config)

	router := gin.Default()
	
	InitializeRouters(router)

	router.Run(":8080")

}

func InitializeRouters(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	urlGroup := router.Group("url")
	{
		urlGroup.POST("/resolve", controller.Resolve)
		urlGroup.POST("/shorten", controller.Shorten)
	}

}
