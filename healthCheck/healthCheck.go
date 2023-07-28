package healthCheck


import(
elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
"strconv"
"fmt"
"io/ioutil"
"encoding/json"
"os"
//"golang.org/x/text/cases"
//"golang.org/x/text/language"
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

func prettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}

func CheckHealth (flingSettings string){

	// var allChecks []HealthChecks
	var allESChecks []Check
	
	// Open our jsonFile
	// jsonFile, err := os.Open(logfile)
	// if err != nil {
    //     fmt.Println(err)
    // }
    // fmt.Println("Successfully Opened logfile")
    // // defer the closing of our jsonFile so that we can parse it later on
    // defer jsonFile.Close()

    // byteValue, _ := ioutil.ReadAll(jsonFile)

	
	// type APP_DATA struct {
	// 	ElasticStats     elasticing.Esstats       `json:"es_stats"`
	// }


	// var appData APP_DATA

	// json.Unmarshal(byteValue, &appData)
	// caser := cases.Title(language.English) //Capitalise first letter
	// fmt.Println("Elasticsearch-> \n\t\tStatus: "+caser.String(appData.ElasticStats[0].Status) + "\n\t\tTotal Nodes: "+appData.ElasticStats[0].NodeTotal)
	
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


	// allChecks = append(allChecks,esChecks)


	allESChecks = append(allESChecks,checkESWatermarkThreshold())
	// allChecks = append(allESChecks,checkESStats())
	setHealthCheckStatus(allESChecks)


	// e, err := json.Marshal(allChecks)
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // fmt.Println(string(e))

    

	
}

func setHealthCheckStatus(checks []Check){

	for _, c := range checks {
        fmt.Printf("Name: %s\n", c.CheckName)
        fmt.Printf("Status: %s\n",c.CheckStatus)
    }

// 	hc := HealthChecks{
// 	HealthCheckName: "Elasticsearch",
// 	HealthCheckStatus: healthy,
// 	Checks: []Check{
// 		},
// 	}
// hc.Checks = append(hc.Checks,c)



}

func checkESStats(){

	Esstats := elasticing.ElasticHealth()

	cluster_status := (*Esstats)[0].Status

	node_total, err := strconv.Atoi((*Esstats)[0].NodeTotal) //Remove percent sign and convert to int
	if err != nil {
        fmt.Println(err)
    }

	fmt.Printf("%+v\n", Esstats)
	//fmt.Println(prettyPrint(Esstats))
	fmt.Printf("Cluster Status = %s",cluster_status)
	fmt.Printf("Total Nodes = %d",node_total)

	healthy := true
	checkInfo := ""


	if (node_total >= 1 && cluster_status=="red"){

		healthy = false
    	checkInfo = "Elasticsearch cluster in a unhealthy state - "+cluster_status

	} else if (node_total > 1 && cluster_status =="yellow"){
    
    	healthy = false
    	checkInfo = "Elasticsearch cluster in a unhealthy state - "+cluster_status
    
    } else{
    	checkInfo = "Cluster is healthy"
    }

    c := Check {
    	CheckName: "Cluster status",
    	CheckStatus:healthy,
    	CheckInfo:checkInfo,
    }

    fmt.Println("Check name = "+c.CheckName)
    fmt.Println("Check status = "+c.CheckStatus)
    fmt.Println("Check Info = "+c.CheckInfo)
}

func checkESWatermarkThreshold() Check{

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







return c

}

