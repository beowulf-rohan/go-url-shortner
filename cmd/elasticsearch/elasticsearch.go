package elasticsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	config *model.Config
)

func Init(configurations *model.Config) {
	config = configurations
}

type ElasticClient struct {
	Client *elasticsearch.Client
}

func GetElasticClient() (*ElasticClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.ElasticEndpoint,
		},
		// Username: config.ElasticUsername,
		// Password: config.ElasticPassword,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 360 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			// TLSClientConfig: &tls.Config{
			// 	MinVersion: tls.VersionTLS13,
			// },
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("error creating elastic client:", err)
		return nil, err
	}

	_, err = es.Info()
	if err != nil {
		log.Println("info error for elastic:", err)
		return nil, err
	}

	return &ElasticClient{Client: es}, nil
}

func (ec *ElasticClient) CreateIndex(indexName string) error {
	exists, err := ec.CheckIfIndexExists(indexName)
	if err != nil {
		log.Println("error checking if index exists:", err)
		return err
	}

	if exists {
		log.Println("Index already exists:", indexName)
		return nil
	}

	log.Printf("creating %s index.....", indexName)
	mapping := getIndexMapping(indexName)

	reqBody := strings.NewReader(mapping)
	res, err := ec.Client.Indices.Create(indexName, ec.Client.Indices.Create.WithBody(reqBody))
	if err != nil {
		log.Println("error creating elastic index:", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			log.Println("error decoding", err)
			return err
		}
		return fmt.Errorf("failed to create index. error code: %d, response body: %v", res.StatusCode, res.Body)
	}
	log.Printf("created %s index sucessfully...", indexName)
	return nil
}

func (ec *ElasticClient) CheckIfIndexExists(indexName string) (bool, error) {
	res, err := ec.Client.Indices.Exists([]string{indexName})
	if err != nil {
		return false, fmt.Errorf("error checking index existence: %v", err)
	}
	defer res.Body.Close()

	exists := res.StatusCode == 200
	return exists, nil
}
