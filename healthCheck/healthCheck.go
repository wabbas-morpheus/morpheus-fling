package healthCheck

import (
	"encoding/json"
	"fmt"
	elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
	"io/ioutil"
	"os"
	"strconv"
	//"golang.org/x/text/cases"
	//"golang.org/x/text/language"
	//sysgatherer "github.com/wabbas-morpheus/morpheus-fling/sysGatherer"
)

type HealthChecks struct {
	HealthCheckName   string  `json:"healthCheckName"`
	HealthCheckStatus bool    `json:"healthCheckStatus"`
	Checks            []Check `json:"checks"`
}

type Check struct {
	CheckName   string `json:"checkName"`
	CheckStatus bool   `json:"checkStatus"`
	CheckInfo   string `json:"checkInfo"`
}

type FlingSettings struct {
	MorpheusApiToken string `json:"morpheus_api_token"`
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func CheckHealth(flingSettings string) {

	var allHealthChecks []HealthChecks
	var allESChecks []Check

	// Open our jsonFile
	jsonFile2, err := os.Open(flingSettings)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened morpheus-fling-settings.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile2.Close()

	byteValue2, _ := ioutil.ReadAll(jsonFile2)

	var flSettings FlingSettings

	json.Unmarshal(byteValue2, &flSettings)

	fmt.Println(prettyPrint(flSettings))

	allESChecks = append(allESChecks, checkESWatermarkThreshold())
	allESChecks = append(allESChecks, checkESStats())

	allHealthChecks = append(allHealthChecks, setHealthCheckStatus(allESChecks, "Elasticsearch"))

	e, err := json.Marshal(allHealthChecks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))

}

func setHealthCheckStatus(checks []Check, checkHeading string) HealthChecks {

	status := true

	for _, c := range checks {
		if !c.CheckStatus {
			status = false
		}
		// fmt.Printf("Name: %s\n", c.CheckName)
		// fmt.Printf("Status: %t\n",c.CheckStatus)
	}

	hc := HealthChecks{
		HealthCheckName:   checkHeading,
		HealthCheckStatus: status,
		Checks:            []Check{},
	}
	hc.Checks = checks

	return hc

}

func checkESStats() Check {

	Esstats := elasticing.ElasticHealth()

	cluster_status := (*Esstats)[0].Status

	node_total, err := strconv.Atoi((*Esstats)[0].NodeTotal) //Remove percent sign and convert to int
	if err != nil {
		fmt.Println(err)
	}

	healthy := true
	checkInfo := ""

	if node_total >= 1 && cluster_status == "red" {

		healthy = false
		checkInfo = "Elasticsearch cluster in a unhealthy state - " + cluster_status

	} else if node_total > 1 && cluster_status == "yellow" {

		healthy = false
		checkInfo = "Elasticsearch cluster in a unhealthy state - " + cluster_status

	} else {
		checkInfo = "Cluster is healthy"
	}

	c := Check{
		CheckName:   "Cluster status",
		CheckStatus: healthy,
		CheckInfo:   checkInfo,
	}

	// fmt.Println("Check name = "+c.CheckName)
	// fmt.Printf("Check status = %t",c.CheckStatus)
	// fmt.Println("Check Info = "+c.CheckInfo)

	return c
}

func checkESWatermarkThreshold() Check {

	//Get water settings from elasticsearch
	esWaterMarkSettings := elasticing.ElasticWatermarkSettings()

	low := esWaterMarkSettings.Low
	lowNumberOnly, err := strconv.Atoi(low[0 : len(low)-1]) //Remove percent sign and convert to int
	if err != nil {
		fmt.Println(err)
	}

	high := esWaterMarkSettings.High
	highNumberOnly, err := strconv.Atoi(high[0 : len(high)-1]) //Remove percent sign and convert to int
	if err != nil {
		fmt.Println(err)
	}

	flood := esWaterMarkSettings.FloodStage
	floodNumberOnly, err := strconv.Atoi(flood[0 : len(flood)-1]) //Remove percent sign and convert to int
	if err != nil {
		fmt.Println(err)
	}

	//get total used storage from the app node. Need to find another way to get current storage as the sysinfo not
	//compatible with other OS
	currentStorage := 0 //sysgatherer.GetStorageUsed()

	// fmt.Println("Low = " + strconv.Itoa(lowNumberOnly))
	// fmt.Println("High = " + strconv.Itoa(highNumberOnly))
	// fmt.Println("Flood Stage = " + strconv.Itoa(floodNumberOnly))
	// fmt.Println("Storage Used = " + strconv.Itoa(sysgatherer.GetStorageUsed()))

	//Check if elasticsearch watermark thresholds has been reached
	healthy := true
	checkInfo := ""
	if currentStorage >= lowNumberOnly && currentStorage < highNumberOnly {

		healthy = false
		checkInfo = "Low (" + strconv.Itoa(lowNumberOnly) + ") watermark threshold has been reached"
	} else if currentStorage >= highNumberOnly && currentStorage < floodNumberOnly {
		healthy = false
		checkInfo = "High (" + strconv.Itoa(highNumberOnly) + ") watermark threshold has been reached"
	} else if currentStorage >= floodNumberOnly {
		healthy = false
		checkInfo = "Flood (" + strconv.Itoa(floodNumberOnly) + ") watermark threshold has been reached"
	} else {
		checkInfo = "Watermark threshold has not been reached"
	}

	c := Check{
		CheckName:   "Watermark",
		CheckStatus: healthy,
		CheckInfo:   checkInfo,
	}

	return c

}
