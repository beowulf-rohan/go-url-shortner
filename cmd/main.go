package main

import (
	"log"

	"github.com/beowulf-rohan/go-url-shortner/api/controller"
	"github.com/beowulf-rohan/go-url-shortner/api/services"
	"github.com/beowulf-rohan/go-url-shortner/elasticsearch"
	"github.com/beowulf-rohan/go-url-shortner/model"
	database "github.com/beowulf-rohan/go-url-shortner/redis"
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
		"ES_ENDPOINT",
		"ES_USERNAME",
		"ES_PASSWORD",
		"URL_METADATA_ES_INDEX",
	}
)

func main() {

	log.Println("loading configs.....")
	config, err := utils.LoadEnvVaraibles(envVariableList)
	if err != nil {
		log.Fatal("Error loading configs from env file", err)
	}
	log.Println("configs loaded successfully....")

	utils.Init(config)
	database.Init(config)
	services.Init(config)
	elasticsearch.Init(config)

	router := gin.Default()
	InitializeRouters(router)
	InitializeElasticSerch(config)

	log.Fatal(router.Run(":" + config.AppPort))

}

func InitializeRouters(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/:url", controller.Resolve)
	router.POST("/api/v1", controller.Shorten)

}

func InitializeElasticSerch(config *model.Config) {
	ElasticClient, err := elasticsearch.GetElasticClient()
	if err != nil {
		log.Fatal(err)
	}
	ElasticClient.CreateIndex(config.UrlMetadataIndex)
	// TODO: add more indices here.

}
