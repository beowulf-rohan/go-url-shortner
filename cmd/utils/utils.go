package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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

func GenerateShortUrl(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}

func GetShortenQuery(url string) string {
	return fmt.Sprintf(`{
		"query": {
			"term": {
				"url.keyword": "%s"
			}
		}
	}`, url)
}

func GetResolveQuery(shortUrl string) string {
	return fmt.Sprintf(`{
		"query": {
			"term": {
				"short_url.keyword": "%s"
			}
		}
	}`, shortUrl)
}