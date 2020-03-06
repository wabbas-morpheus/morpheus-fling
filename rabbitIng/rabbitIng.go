package rabbiting

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type RabbitResults struct{
	Name      string `json:name`
	VHost     string `json:vhost`
	Messages  int    `json:messages`
	Memory    int    `json:memory`
	Node      string `json:node`
	Policy    string `json:policy`
	Consumers int    `json:consumer`
}

func RabbitStats() []RabbitResults {
	manager := "http://127.0.0.1:15672/api/queues/"
	client := &http.Client{}
	req, err := http.NewRequest("GET", manager, nil)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	req.SetBasicAuth("morpheus", "7f4a6fd594a8b962")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	//value := &RabbitResults{}
	value := make([]Queue, 0)

	json.NewDecoder(resp.Body).Decode(&value)

	resultjson, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}

	fmt.Fprintf(os.Stdout, "%s", resultjson)
	return value
}