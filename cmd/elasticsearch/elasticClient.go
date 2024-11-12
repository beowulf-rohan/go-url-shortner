package elasticsearch

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

var EsClient *elasticsearch.Client

func GetElasticClient(url, index, username, password string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Username: username,
		Password: password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 360 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("error creating elastic client:", err)
		return &elasticsearch.Client{}, err
	}

	_, err = es.Info()
	if err != nil {
		log.Println("info error for elastic:", err)
		return &elasticsearch.Client{}, err
	}
	
	log.Println("elastic client created")
	return es, nil
}
