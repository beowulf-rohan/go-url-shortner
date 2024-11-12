package model

type Config struct {
	DbAddr           string `json:"db_addr"`
	DbPass           string `json:"db_pass"`
	AppPort          string `json:"app_port"`
	Domain           string `json:"domain"`
	ApiQuota         int    `json:"api_quota"`
	ElasticEndpoint  string `json:"elastic_endpoint"`
	ElasticUsername  string `json:"elastic_username"`
	ElasticPassword  string `json:"elastic_password"`
	UrlMetadataIndex string `json:"url_metadata_index"`
}