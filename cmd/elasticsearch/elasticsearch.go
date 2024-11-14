package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticClient struct {
	Client *elasticsearch.Client
	Index  string
	Ctx    context.Context
}

func GetElasticClient(index string) (*ElasticClient, error) {
	config := config.GlobalConfig

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

	return &ElasticClient{
		Client: es,
		Index:  index,
		Ctx:    context.Background(),
	}, nil
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

func (ec *ElasticClient) CreateIndex() error {
	exists, err := ec.CheckIfIndexExists(ec.Index)
	if err != nil {
		log.Println("error checking if index exists:", err)
		return err
	}

	if exists {
		log.Println("Index already exists:", ec.Index)
		return nil
	}

	log.Printf("creating %s index.....", ec.Index)
	mapping := getIndexMapping(ec.Index)

	reqBody := strings.NewReader(mapping)
	res, err := ec.Client.Indices.Create(ec.Index, ec.Client.Indices.Create.WithBody(reqBody))
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
	log.Printf("created %s index sucessfully...", ec.Index)
	return nil
}

func (ec *ElasticClient) DeleteIndex() error {
	log.Printf("deleting %s index.....", ec.Index)
	res, err := ec.Client.Indices.Delete([]string{ec.Index}, ec.Client.Indices.Delete.WithContext(context.Background()))
	if err != nil {
		log.Println("error deleting elastic index:", err)
		return err
	}
	defer res.Body.Close()
	log.Printf("deleted %s index sucessfully...", ec.Index)
	return nil
}

func (ec *ElasticClient) PushToElastic(doc interface{}, docID string) error {

	var bulkBody bytes.Buffer
	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshalling document to JSON: %w", err)
	}
	indexLine := fmt.Sprintf(`{ "index" : { "_index" : "%s", "_id": "%s"} } %s`, ec.Index, docID, "\n")
	bulkBody.WriteString(indexLine)
	bulkBody.Write(jsonDoc)
	bulkBody.WriteString("\n")

	req := esapi.BulkRequest{
		Body:    strings.NewReader(bulkBody.String()),
		Refresh: "true",
	}
	res, err := req.Do(ec.Ctx, ec.Client)
	if err != nil {
		return fmt.Errorf("error indexing document to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var errMsg map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			return fmt.Errorf("error parsing Elasticsearch response: %w", err)
		}
		return fmt.Errorf("elasticsearch indexing error: %v", errMsg)
	}

	return nil
}

func (ec *ElasticClient) GetFromElastic(docID string) (*model.Response, error) {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"url.keyword": "%s"
			}
		}
	}`, docID)

	res, err := ec.Client.Search(
		ec.Client.Search.WithIndex(ec.Index),
		ec.Client.Search.WithBody(strings.NewReader(query)),
		ec.Client.Search.WithContext(context.Background()),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching document in Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var errMsg map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			return nil, fmt.Errorf("error parsing Elasticsearch response: %w", err)
		}
		return nil, fmt.Errorf("elasticsearch retrieval error: %v", errMsg)
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source model.Response `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("error decoding search result: %w", err)
	}

	if len(searchResult.Hits.Hits) == 0 {
		return nil, fmt.Errorf("document with URL '%s' not found in index '%s'", docID, ec.Index)
	}

	return &searchResult.Hits.Hits[0].Source, nil
}

func (ec *ElasticClient) GetAllFromElastic() error {
	log.Println("fetching all documents")
	searchRequest := esapi.SearchRequest{
		Index: []string{ec.Index},
		Body:  strings.NewReader(`{"query": {"match_all": {}}}`),
	}

	res, err := searchRequest.Do(context.Background(), ec.Client)
	if err != nil {
		return fmt.Errorf("error searching documents in Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var errMsg map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			return fmt.Errorf("error parsing Elasticsearch response: %w", err)
		}
		return fmt.Errorf("elasticsearch search error: %v", errMsg)
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source model.Response `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return fmt.Errorf("error decoding search result: %w", err)
	}

	if len(searchResult.Hits.Hits) > 0 {
		for _, hit := range searchResult.Hits.Hits {
			log.Printf("Found document: %+v", hit.Source)
		}
	} else {
		log.Println("No documents found.")
	}

	return nil
}
