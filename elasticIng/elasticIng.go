package elasticing

import (
	"context"
	"encoding/json"
	"fmt"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/mitchellh/mapstructure"
	"log"
	"sync"
	"net/http"
	"io/ioutil"
	"os"
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

type Eswater_mark_settings struct {

	MaxHeadRoom string `json:"flood_stage.frozen.max_headroom"`
	FloodStage string `json:"flood_stage"`
	High string `json:"high"`
	Low string `json:"low"`
	EnableSDN string `json:"enable_for_single_data_node"`
	FloodStageFrozen string `json:"flood_stage.frozen"`

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
func ElasticIndices() []Esindices {

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

	indexSlice := make([]Esindices, len(r))

	var wg sync.WaitGroup
	for i, element := range r {
		wg.Add(1)
		go func(i int, element map[string]interface{}) {
			defer wg.Done()
			result := Esindices{}
			cfg := &mapstructure.DecoderConfig{
				Metadata: nil,
				Result:   &result,
				TagName:  "json",
			}
			decoder, _ := mapstructure.NewDecoder(cfg)
			decoder.Decode(element)
			indexSlice[i] = result
		}(i, element)
	}
	wg.Wait()

	return indexSlice

}



func ElasticWatermarkSettings(){

	response, err := http.Get("http://localhost:9200/_cluster/settings?pretty&include_defaults")

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

	type APP_DATA struct {
		ElasticWaterMark elasticing.Eswater_mark_settings `json:"watermark"`
	}
	var appData APP_DATA

	json.Unmarshal(responseData, &appData)

    fmt.Println(string(responseData))

	fmt.Println(string(appData.ElasticWaterMark.High))


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