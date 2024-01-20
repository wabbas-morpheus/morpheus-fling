package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mholt/archiver"
	elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
	encryptText "github.com/wabbas-morpheus/morpheus-fling/encryptText"
	portscanner "github.com/wabbas-morpheus/morpheus-fling/portScanner"
	rabbiting "github.com/wabbas-morpheus/morpheus-fling/rabbitIng"
	//"github.com/zcalusic/sysinfo"
	"log"
	"os"
	"path"
	"time"
)

var (
	defaultPath   = "."
	bundlerPtr    = flag.String("bundler", path.Join(defaultPath, "bundler.zip"), "a string")
	privatekeyPtr = flag.String("privkey", path.Join(defaultPath, "morpheus.pem"), "a string")
	extractPtr    = flag.Bool("extract", true, "a bool")
)

const helpText = `morpheus-fling [options]
Options:

-bundler	Path to bundled encrypted file.  Defaults to "./bundler.zip".
-privkey	Path to the private key file used for decryption.  Defaults to "./morpheus.pub"
-help		Prints this text.

Examples:
	Decrypt the encrypted bundled file at /Users/wabbas/Dev/tmp/bundler.zip
   		$ ./morpheus-fling-osx -bundler /Users/wabbas/Dev/tmp/bundler.zip

	Decrypt the encrypted bundled file at '/Users/wabbas/Dev/tmp/bundler.zip' with private key at '../../bin/morpheus.pem'
   		$ ./morpheus-fling-osx -privkey ../../bin/morpheus.pem -bundler /Users/wabbas/Dev/tmp/bundler.zip

`

type Results struct {
	ElasticStats    *elasticing.Esstats             `json:"es_stats"`
	ElasticIndices  []elasticing.Esindices          `json:"es_indices"`
	ElasticSettings *elasticing.ESWaterMarkSettings `json:"es_settings"`
	//System           *sysinfo.SysInfo                `json:"system_stats"`
	Scans            []portscanner.ScanResult  `json:"port_scans,omitempty"`
	RabbitStatistics []rabbiting.RabbitResults `json:"rabbit_stats"`
	MorphLogs        string                    `json:"morpheus_logs"`
}

type ESResults struct {
	ElasticStats    *elasticing.Esstats             `json:"es_stats"`
	ElasticIndices  []elasticing.Esindices          `json:"es_indices"`
	ElasticSettings *elasticing.ESWaterMarkSettings `json:"es_settings"`
}

type RabbitResults struct {
	RabbitStatistics []rabbiting.RabbitResults `json:"rabbit_stats"`
}

//type SystemResults struct {
//	System *sysinfo.SysInfo `json:"system_stats"`
//}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FileWrtr takes content and an outfile and appends content to the outfile
func FileWrtr(content string, fileName string) {
	//Remove existing files
	if fileExists(fileName) {
		e := os.Remove(fileName)
		if e != nil {
			log.Fatal(e)
		}
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		log.Println(err)
	}

}

func extractBundle() {

	// Extract the encrypted bundle
	t := time.Now()
	timeStamp := t.Format("20060102150405")
	folderName := "extracted_" + timeStamp
	if err := archiver.Unarchive(*bundlerPtr, folderName+"/"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Extracting Bundle File")
	nonText, err := os.ReadFile(folderName + "/encrypted_logs.json")
	if err != nil {
		log.Fatal("Can't load output file", err)
	}

	nonKey, err := os.ReadFile(folderName + "/bundlerkey.enc")
	if err != nil {
		log.Fatal("Can't load key file", err)
	}

	decryptedText := encryptText.DecryptItAll(*privatekeyPtr, nonText, nonKey)
	var jsonBlob = []byte(decryptedText)
	var results Results
	var esResults ESResults
	var rabbitResults RabbitResults
	//var system_results SystemResults

	err = json.Unmarshal(jsonBlob, &results)
	if err != nil {
		fmt.Println("error:", err)
	}

	esResults.ElasticStats = results.ElasticStats
	esResults.ElasticSettings = results.ElasticSettings
	esResults.ElasticIndices = results.ElasticIndices
	rabbitResults.RabbitStatistics = results.RabbitStatistics
	//system_results.System = results.System

	//fmt.Printf("%+v", results.MorphLogs)
	//fmt.Println("Decrypted Text = ",decryptedText)
	FileWrtr(decryptedText, folderName+"/all_logs.json")
	FileWrtr(results.MorphLogs, folderName+"/morpheus_current.log")
	FileWrtr(dumps(results.RabbitStatistics), folderName+"/rabbit_stats.log")
	FileWrtr(dumps(results.ElasticStats), folderName+"/elastic_status.log")
	FileWrtr(dumps(results.ElasticSettings), folderName+"/elastic_settings.log")
	FileWrtr(dumps(results.ElasticIndices), folderName+"/elastic_indices.log")
	//FileWrtr(dumps(results.System), folderName+"/system.log")

}

func dumps(data interface{}) string {
	jData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}
	return string(jData)
}

func main() {
	flag.Usage = help
	flag.Parse()

	if *extractPtr { //Extract

		extractBundle()

	}
}

func help() {
	fmt.Println(helpText)
	os.Exit(1)
}