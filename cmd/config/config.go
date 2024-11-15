package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/joho/godotenv"
)

const (
	LOCAL_CONFIG_FILE_PATH = ".env"
)

var GlobalConfig model.Config

func GetConfigFilePath() (configFilePath string) {
	// add env based config file paths here
	configFilePath = LOCAL_CONFIG_FILE_PATH
	return
}

func LoadEnvVaraibles(envVariableList []string) error {
	config := model.Config{}
	configFilePath := GetConfigFilePath()
	err := godotenv.Load(configFilePath)
	if err != nil {
		return err
	}

	for _, variable := range envVariableList {
		switch variable {
		case "DB_ADDR":
			config.DbAddr = os.Getenv("DB_ADDR")
			if config.DbAddr == "" {
				return fmt.Errorf("no env value found for \"DB_ADDR\" in env file")
			}
		case "DB_PASS":
			config.DbPass = os.Getenv("DB_PASS")
			// if config.DbPass == "" {
			// 	return fmt.Errorf("no env value found for \"DB_PASS\" in env file")
			// }
		case "APP_PORT":
			config.AppPort = os.Getenv("APP_PORT")
			if config.AppPort == "" {
				return fmt.Errorf("no env value found for \"APP_PORT\" in env file")
			}
		case "DOMAIN":
			config.Domain = os.Getenv("DOMAIN")
			if config.Domain == "" {
				return fmt.Errorf("no env value found for \"DOMAIN\" in env file")
			}
		case "API_RATE_LIMIT":
			apiRateLimit := os.Getenv("API_RATE_LIMIT")
			if apiRateLimit == "" {
				return fmt.Errorf("no env value found for \"API_RATE_LIMIT\" in env file")
			} else if config.ApiRateLimit, err = strconv.Atoi(apiRateLimit); err != nil {
				return err
			}
		case "ES_ENDPOINT":
			config.ElasticEndpoint = os.Getenv("ES_ENDPOINT")
			if config.ElasticEndpoint == "" {
				return fmt.Errorf("no env value found for \"ES_ENDPOINT\" in env file")
			}
		case "ES_USERNAME":
			config.ElasticUsername = os.Getenv("ES_USERNAME")
			if config.ElasticUsername == "" {
				return fmt.Errorf("no env value found for \"ES_USERNAME\" in env file")
			}
		case "ES_PASSWORD":
			config.ElasticPassword = os.Getenv("ES_PASSWORD")
			if config.ElasticPassword == "" {
				return fmt.Errorf("no env value found for \"ES_PASSWORD\" in env file")
			}
		case "URL_METADATA_ES_INDEX":
			config.UrlMetadataIndex = os.Getenv("URL_METADATA_ES_INDEX")
			if config.UrlMetadataIndex == "" {
				return fmt.Errorf("no env value found for \"URL_METADATA_ES_INDEX\" in env file")
			}
		}
	}
	GlobalConfig = config

	return nil
}
