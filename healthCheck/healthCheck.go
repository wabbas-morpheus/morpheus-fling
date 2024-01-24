package healthCheck

import (
	"encoding/json"
	"fmt"
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
	//jsonFile2, err := os.Open(flingSettings)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("Successfully Opened morpheus-fling-settings.json")
	// defer the closing of our jsonFile so that we can parse it later on
	//defer jsonFile2.Close()
	//
	//byteValue2, _ := ioutil.ReadAll(jsonFile2)
	//
	//var flSettings FlingSettings
	//
	//json.Unmarshal(byteValue2, &flSettings)

	//fmt.Println(prettyPrint(flSettings))

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
