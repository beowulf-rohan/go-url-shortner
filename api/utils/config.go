package utils

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

func GetConfigFilePath() (configFilePath string) {
	// add env based config file paths here
	configFilePath = LOCAL_CONFIG_FILE_PATH
	return
}

func LoadEnvVaraibles(envVariableList []string) (*model.Config, error) {
	config := model.Config{}
	configFilePath := GetConfigFilePath()
	err := godotenv.Load(configFilePath)
	if err != nil {
		return &model.Config{}, err
	}

	for _, variable := range envVariableList {
		switch variable {
		case "DB_ADDR":
			config.DbAddr = os.Getenv("DB_ADDR")
			if config.DbAddr == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"DB_ADDR\" in env file")
			}
		case "DB_PASS":
			config.DbPass = os.Getenv("DB_PASS")
			if config.DbPass == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"DB_PASS\" in env file")
			}
		case "APP_PORT":
			config.AppPort = os.Getenv("APP_PORT")
			if config.AppPort == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"APP_PORT\" in env file")
			}
		case "DOMAIN":
			config.Domain = os.Getenv("DOMAIN")
			if config.Domain == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"DOMAIN\" in env file")
			}
		case "API_QUOTA":
			apiQuota := os.Getenv("API_QUOTA")
			if apiQuota == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"API_QUOTA\" in env file")
			} else if config.ApiQuota, err = strconv.Atoi(apiQuota); err != nil {
				return &model.Config{}, err
			}
		case "ES_ENDPOINT":
			config.ElasticEndpoint = os.Getenv("ES_ENDPOINT")
			if config.ElasticEndpoint == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"ES_ENDPOINT\" in env file")
			}
		case "ES_PASSWORD":
			config.ElasticPassword = os.Getenv("ES_PASSWORD")
			if config.ElasticPassword == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"ES_PASSWORD\" in env file")
			}
		case "URL_METADATA_ES_INDEX":
			config.UrlMetadataIndex = os.Getenv("URL_METADATA_ES_INDEX")
			if config.UrlMetadataIndex == "" {
				return &model.Config{}, fmt.Errorf("no env value found for \"URL_METADATA_ES_INDEX\" in env file")
			}
		}
	}

	return &config, nil
}
