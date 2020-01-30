package elasticing

import (
	"context"
	"encoding/json"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"bytes"
	//"github.com/mitchellh/mapstructure"
)

//ElasticIndices Cats the active ES indices found
func ElasticIndices() string {

	var r map[string]interface{}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	req := esapi.CatIndicesRequest{
		Format: "json",
		Pretty: false,
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	}

	b, _ := json.Marshal(r)
	return string(b)

}

// ElasticHealth returns a esapi.Response of Health
func ElasticHealth() string {

	//var r map[string]interface{}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	req := esapi.CatHealthRequest{
		Format: "json",
		Pretty: false,
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	//if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	//	log.Printf("Error parsing the response body: %s", err)
	//}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	s := buf.String()

	return s
}