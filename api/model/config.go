package model

type Config struct {
	DbAddr   string `json:"db_addr"`
	DbPass   string `json:"db_pass"`
	AppPort  string `json:"app_port"`
	Domain   string `json:"domain"`
	ApiQuota int    `json:"api_quota"`
	ApiPort  int    `json:"api_port"`
}
