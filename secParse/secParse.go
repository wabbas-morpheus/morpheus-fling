package secparse

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Mysql struct {
	RootPassword     string `json:"root_password"`
	MorpheusPassword string `json:"morpheus_password"`
	OpsPassword      string `json:"ops_password"`
}

type Rabbitmq struct {
	MorpheusPassword  string `json:"morpheus_password"`
	QueueUserPassword string `json:"queue_user_password"`
	Cookie            string `json:"cookie"`
}

type Secret struct {
	Mysql    Mysql
	Rabbitmq Rabbitmq
}

func ParseSecrets(secfilePtr string) Secret {
	jsonFile, err := os.Open(secfilePtr)
	if err != nil {
		log.Fatalf("Error reading secrets file: %s", err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var secrets Secret
	if err := json.Unmarshal(byteValue, &secrets); err != nil {
		log.Fatalf("Error unmarshalling secrets file: %s", err)
	}
	return secrets
}