package utils

import (
	"strings"

	"github.com/beowulf-rohan/go-url-shortner/config"
)

func CheckDomainError(url string) bool {
	config := config.GlobalConfig
	if url == config.Domain {
		return false
	}

	url = strings.Replace(url, "http://", "", 1)
	url = strings.Replace(url, "https://", "", 1)
	url = strings.Replace(url, "www.", "", 1)
	domain := strings.Split(url, "/")[0]

	return domain != config.Domain
}

func EnforceHttp(url string) string {
	if url[0:4] != "http" {
		return "https://" + url
	}
	return url
}
