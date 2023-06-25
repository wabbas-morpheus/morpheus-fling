package healthCheck


import(
elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
"strconv"
"fmt"
"io/ioutil"
"encoding/json"
)

func checkHealth (*logfile ){

	// Open our jsonFile
	jsonFile, err := os.Open(*logfile)
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
	//elasticing.ElasticWatermarkSettings()
	//esWaterMarkSettings := elasticing.ElasticWatermarkSettings()
	//fmt.Println("Watermark = "+string(esWaterMarkSettings))
	esWaterMarkSettings := elasticing.ElasticWatermarkSettings()
	 
	low := esWaterMarkSettings.Low
	lowNumberOnly, err := strconv.Atoi(low[0:len(low)-1]) //Remove percent sign and convert to int

	high := esWaterMarkSettings.High
	highNumberOnly, err := strconv.Atoi(high[0:len(high)-1]) //Remove percent sign and convert to int64

	flood := esWaterMarkSettings.FloodStage
	floodNumberOnly, err := strconv.Atoi(flood[0:len(flood)-1]) //Remove percent sign and convert to int

	fmt.Println("Low = " + strconv.Itoa(lowNumberOnly))
	fmt.Println("High = " + strconv.Itoa(highNumberOnly))
	fmt.Println("Flood Stage = " + strconv.Itoa(floodNumberOnly))
	fmt.Println("Storage Used = " + strconv.Itoa(sysgatherer.GetStorageUsed()))

}

