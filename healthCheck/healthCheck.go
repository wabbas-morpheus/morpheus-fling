package healthCheck


import(
elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
"strconv"
"fmt"
"io/ioutil"
"encoding/json"
"os"
"golang.org/x/text/cases"
"golang.org/x/text/language"
sysgatherer "github.com/wabbas-morpheus/morpheus-fling/sysGatherer"
)


type HealthChecks struct {

	HealthCheckName string `json:"healthCheckName"`
	HealthCheckStatus bool `json:"healthCheckStatus"`
	Checks []Check `json:"checks"`

}

type Check struct {

	CheckName string `json:"checkName"`
	CheckStatus bool `json:"checkStatus"`
	CheckInfo string `json:"checkInfo"`

}

type FlingSettings struct {

	MorpheusApiToken string `json:"morpheus_api_token"`
}

func CheckHealth (logfile string,flingSettings string){

	var allChecks []HealthChecks
	
	// Open our jsonFile
	jsonFile, err := os.Open(logfile)
	if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Successfully Opened users.json")
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    // var result map[string]interface{}
    // json.Unmarshal([]byte(byteValue), &result)

	// for key, value := range result {
	// 	fmt.Printf("Key: %s === Value: %s \n", key, value)
	// }
    // fmt.Println(result["es_stats"][0])

	
	type APP_DATA struct {
		ElasticStats     elasticing.Esstats       `json:"es_stats"`
	}


	var appData APP_DATA

	json.Unmarshal(byteValue, &appData)
	caser := cases.Title(language.English) //Capitalise first letter
	fmt.Println("Elasticsearch-> \n\t\tStatus: "+caser.String(appData.ElasticStats[0].Status) + "\n\t\tTotal Nodes: "+appData.ElasticStats[0].NodeTotal)
	
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

	fmt.Println("access token = "+flSettings.MorpheusApiToken)

	esChecks := checkESWatermarkThreshold()
	allChecks = append(allChecks,esChecks)

	e, err := json.Marshal(allChecks)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(e))

	
}

func checkESWatermarkThreshold() HealthChecks{

	//Get water settings from elasticsearch
	esWaterMarkSettings := elasticing.ElasticWatermarkSettings()


	 
	low := esWaterMarkSettings.Low
	lowNumberOnly, err := strconv.Atoi(low[0:len(low)-1]) //Remove percent sign and convert to int
	if err != nil {
        fmt.Println(err)
    }

	high := esWaterMarkSettings.High
	highNumberOnly, err := strconv.Atoi(high[0:len(high)-1]) //Remove percent sign and convert to int
	if err != nil {
        fmt.Println(err)
    }

	flood := esWaterMarkSettings.FloodStage
	floodNumberOnly, err := strconv.Atoi(flood[0:len(flood)-1]) //Remove percent sign and convert to int
	if err != nil {
        fmt.Println(err)
    }

    //get total used storage from the app node
    currentStorage := sysgatherer.GetStorageUsed()




	// fmt.Println("Low = " + strconv.Itoa(lowNumberOnly))
	// fmt.Println("High = " + strconv.Itoa(highNumberOnly))
	// fmt.Println("Flood Stage = " + strconv.Itoa(floodNumberOnly))
	// fmt.Println("Storage Used = " + strconv.Itoa(sysgatherer.GetStorageUsed()))


//Check if elasticsearch watermark thresholds has been reached
	healthy := true
	checkInfo := ""
	if (currentStorage >= lowNumberOnly && currentStorage < highNumberOnly){
    
    	healthy = false
    	checkInfo = "Low ("+strconv.Itoa(lowNumberOnly)+") watermark threshold has been reached"
    } else if (currentStorage >= highNumberOnly && currentStorage < floodNumberOnly){
    	healthy = false
    	checkInfo = "High ("+strconv.Itoa(highNumberOnly)+") watermark threshold has been reached"
    } else if (currentStorage >= floodNumberOnly){
    	healthy = false
    	checkInfo = "Flood ("+strconv.Itoa(floodNumberOnly)+") watermark threshold has been reached"
    } else{
    	checkInfo = "Watermark threshold has not been reached"
    }

    c := Check {
    	CheckName: "Watermark",
    	CheckStatus:healthy,
    	CheckInfo:checkInfo,
    }


hc := HealthChecks{
		HealthCheckName: "Elasticsearch",
		HealthCheckStatus: healthy,
		Checks: []Check{
			},
		}
hc.Checks = append(hc.Checks,c)




return hc

}

