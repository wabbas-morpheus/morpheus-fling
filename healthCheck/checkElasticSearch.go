package healthCheck

import (
	"bufio"
	"fmt"
	elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
	"log"
	"os/exec"
	"strings"

	//sysgatherer "github.com/wabbas-morpheus/morpheus-fling/sysGatherer"
	"strconv"
)

func runESChecks() []Check {
	var allESChecks []Check
	allESChecks = append(allESChecks, checkESStats())
	allESChecks = append(allESChecks, checkESWatermarkThreshold())
	return allESChecks
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
	currentStorage := GetStorageUsed()

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

func GetStorageUsed() int {

	//get path of the elasticsearch data
	outESFile, err := exec.Command("grep", "data", "/opt/morpheus/embedded/elasticsearch/config/elasticsearch.yml").Output()

	if err != nil {
		log.Fatal(err)
	}

	outputES := string(outESFile[:])

	fPath := ""
	scannerES := bufio.NewScanner(strings.NewReader(outputES))
	for scannerES.Scan() { //iterate of each line
		line := strings.Fields(scannerES.Text()) //convert line text in a list
		fPath = line[1]                          //get storage used info
		//fmt.Printf("fpath = %s\n",fPath)

	}
	if err := scannerES.Err(); err != nil {
		log.Fatal(err)
	}

	//Get available storage info on elasticsearch mount point
	//fmt.Printf("fPath path =%s", fPath)
	out, err := exec.Command("df", "-h", fPath).Output()

	// if there is an error with our execution
	// handle it here
	if err != nil {
		log.Fatal(err)
	}
	// as the out variable defined above is of type []byte we need to convert
	// this to a string or else we will see garbage printed out in our console
	// this is how we convert it to a string
	output := string(out[:])
	used := ""
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() { //iterate of each line
		line := strings.Fields(scanner.Text()) //convert line text in a list
		storageUsedPercent := line[4]          //get storage used info
		used = storageUsedPercent[0 : len(storageUsedPercent)-1]

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	rtn, err := strconv.Atoi(used) //Convert to integer
	fmt.Println("Storage used= " + used)
	return rtn
}
