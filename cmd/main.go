package main

import (
	"log"

	"github.com/beowulf-rohan/go-url-shortner/api/controller"
	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/elasticsearch"
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
		// "ES_USERNAME",
		// "ES_PASSWORD",
		"URL_METADATA_ES_INDEX",
	}
)

func main() {

	log.Println("loading configs.....")
	err := config.LoadEnvVaraibles(envVariableList)
	if err != nil {
		log.Fatal("Error loading configs from env file", err)
	}
	log.Println("configs loaded successfully....")

	router := gin.Default()
	InitializeRouters(router)
	InitializeElasticSerch()

	log.Fatal(router.Run(":" + config.GlobalConfig.AppPort))

}

func InitializeRouters(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	router.Use(gin.Recovery())

	router.GET("/:shortURL", controller.Resolve)
	router.POST("/shorten", controller.Shorten)

}

func InitializeElasticSerch() {
	config := config.GlobalConfig
	ElasticClient, err := elasticsearch.GetElasticClient(config.UrlMetadataIndex)
	if err != nil {
		log.Fatal(err)
	}
	ElasticClient.CreateIndex()
	// TODO: add more indices here.

}
