package rabbiting

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type RabbitResults struct {
	Name      string `json:name`
	VHost     string `json:vhost`
	Messages  int    `json:messages`
	Memory    int    `json:memory`
	Node      string `json:node`
	Policy    string `json:policy`
	Consumers int    `json:consumer`
}

func RabbitStats(user string, password string) []RabbitResults {
	value := make([]RabbitResults, 0)
	if RabbitManagementEnabled() {
		manager := "http://127.0.0.1:15672/api/queues/"
		client := &http.Client{}
		req, err := http.NewRequest("GET", manager, nil)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}

		req.SetBasicAuth(user, password)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}

		//value := &RabbitResults{}

		json.NewDecoder(resp.Body).Decode(&value)

		//resultjson, err := json.MarshalIndent(value, "", " ")
		//if err != nil {
		//	log.Fatal("Can't encode to JSON", err)
		//}

		//fmt.Fprintf(os.Stdout, "%s", resultjson)

	} else {
		var r RabbitResults
		r.Name = "API Error - Unable to retrieve rabbit stats. Make sure rabbitmq management plugin is enabled"
		value = append(value, r)

	}
	return value
}

func RabbitManagementEnabled() bool {
	mgmStatus := false
	data, err := os.ReadFile("/opt/morpheus/embedded/rabbitmq/etc/enabled_plugins")
	if err != nil {
		log.Fatal(err)
	}
	mgmStatus = strings.Contains(string(data), "rabbitmq_management")

	return mgmStatus
}
