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

func CheckHealth (logfile string){

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
	
	checkESWatermarkThreshold()
	
}

func checkESWatermarkThreshold(){

	esWaterMarkSettings := elasticing.ElasticWatermarkSettings()
	 
	low := esWaterMarkSettings.Low
	lowNumberOnly, err := strconv.Atoi(low[0:len(low)-1]) //Remove percent sign and convert to int
	if err != nil {
        fmt.Println(err)
    }

	high := esWaterMarkSettings.High
	highNumberOnly, err := strconv.Atoi(high[0:len(high)-1]) //Remove percent sign and convert to int64
	if err != nil {
        fmt.Println(err)
    }

	flood := esWaterMarkSettings.FloodStage
	floodNumberOnly, err := strconv.Atoi(flood[0:len(flood)-1]) //Remove percent sign and convert to int
	if err != nil {
        fmt.Println(err)
    }

    currentStorage := sysgatherer.GetStorageUsed()




	fmt.Println("Low = " + strconv.Itoa(lowNumberOnly))
	fmt.Println("High = " + strconv.Itoa(highNumberOnly))
	fmt.Println("Flood Stage = " + strconv.Itoa(floodNumberOnly))
	fmt.Println("Storage Used = " + strconv.Itoa(sysgatherer.GetStorageUsed()))

	health := true 
	if (currentStorage >= lowNumberOnly && currentStorage < highNumberOnly){
    	fmt.Println("Low watermark threshould has been reached")
    	health = false
    } else if (currentStorage >= highNumberOnly && currentStorage < floodNumberOnly){
    	fmt.Println("High watermark threshould has been reached")
    	health = false
    } else if (currentStorage >= floodNumberOnly){
    	fmt.Println("Flood watermark threshould has been reached")
    	health = false
    } else{
    	fmt.Println("Watermark threshold has not been reached")
    }




	// checkResults := make(map[string]interface{})

	// checkResults["Elasticsearch Health"] = map[string]interface{}{
	// 	"status":health,
	// 	"checks": map[string]interface{}{
	// 		"watermark threshold": health,
	// 	},
	// }

	// for _,m := range checkResults{

	// 	for k, v := range m {
	// 		fmt.Println(k," value is ", v)
	// 	}
		

	// }

	// fmt.Println(checkResults)
     

}

