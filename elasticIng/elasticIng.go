package elasticing

import (
	"fmt"
	"context"
	"encoding/json"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"github.com/mitchellh/mapstructure"
)


type Esstats []struct {
	Epoch string `json:"epoch"`
	Timestamp string `json:"timestamp"`
	Cluster string `json:"cluster"`
	Status string `json:"status"`
	NodeTotal string `json:"node.total"`
	NodeData string `json:"node.data"`
	Shards string `json:"shards"`
	Pri string `json:"pri"`
	Relo string `json:"relo"`
	Init string `json:"init"`
	Unassign string `json:"unassign"`
	PendingTasks string `json:"pending_tasks"`
	MaxTaskWaitTime string `json:"max_task_wait_time"`
	ActiveShardsPercent string `json:"active_shards_percent"`
}

type Esindices struct {
	Health string `json:"health"`
	Status string `json:"status"`
	Index string `json:"index"`
	Uuid string `json:"uuid"`
	Pri int `json:"pri"`
	Rep int `json:"rep"`
	DocsCount int `json:"docs.count"`
	DocsDeleted int `json:"docs.deleted"`
	StoreSize string `json:"store.size"`
	PriStoreSize string `json:"pri.store.size"`
}


//ElasticIndices Cats the active ES indices found
func ElasticIndices() []*Esindices {

	var r []map[string]interface{}

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

	indexSlice := make([]*Esindices, len(r))

	for i, element := range r {
		result := &Esindices{}
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &result,
			TagName:  "json",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		decoder.Decode(element)
		indexSlice[i] = result
	}

	return indexSlice

}

// ElasticHealth returns a esapi.Response of Health
func ElasticHealth() *Esstats {
	var r []map[string]interface{}
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	req := esapi.CatHealthRequest{
		Format: "json",
		Pretty:	false,
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	defer res.Body.Close()

	result := &Esstats{}

	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &result,
		TagName:  "json",
	}

	decoder, _ := mapstructure.NewDecoder(cfg)

	decoder.Decode(r)

	data, err := json.MarshalIndent(&result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))

	return result
}