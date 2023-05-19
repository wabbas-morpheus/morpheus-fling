package main

import (
	"encoding/json"
	"flag"
	"fmt"
	elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
	rabbiting "github.com/wabbas-morpheus/morpheus-fling/rabbitIng"
	secparse "github.com/wabbas-morpheus/morpheus-fling/secParse"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	encryptText "github.com/wabbas-morpheus/morpheus-fling/encryptText"
	filereader "github.com/wabbas-morpheus/morpheus-fling/fileReader"
	portscanner "github.com/wabbas-morpheus/morpheus-fling/portScanner"
	sysgatherer "github.com/wabbas-morpheus/morpheus-fling/sysGatherer"
	"github.com/mholt/archiver"
	"github.com/zcalusic/sysinfo"
)

var (
	infilePtr  = flag.String("infile", "", "a string")
	secfilePtr = flag.String("secfile", "/etc/morpheus/morpheus-secrets.json", "a string")
	outfilePtr = flag.String("outfile", path.Join(".", "output.json"), "a string")
	uLimit     = flag.Int64("ulimit", 1024, "an integer")
	logfilePtr = flag.String("logfile", "/var/log/morpheus/morpheus-ui/current", "a string")
	bundlerPtr = flag.String("bundler", "/tmp/bundler.zip", "a string")
	keyfilePtr = flag.String("keyfile", "/tmp/bundlerkey.enc", "a string")
	pubPtr     = flag.String("pubfile", "/tmp/morpheus.pub", "a string")
	privatekeyPtr = flag.String("privatefile", "/root/morpheus.pem", "a string")
	extractPtr    = flag.Bool("extract",false,"a bool")
	healthPtr    = flag.Bool("health",false,"a bool")
)

const helpText = `morpheus-fling [options]
Options:
-infile     The source file for network port scanning.  If none is provided port scans will be skipped.
-secfile    The morpheus secrets file.  Defaults to "/etc/morpheus/morpheus-secrets.json".
-outfile    The destination directory of the generated package, "output.txt" by default.
-ulimit     Ulimit of the system, defaults to 1024.
-logfile    Logfile to add to the bundle.  Defaults to "/var/log/morpheus/morpheus-ui/current".
-bundler    Path and file to bundle into.  Defaults to "/tmp/bundler.zip".
-keyfile    Path and file to put the public key encrypted AES-GCM key into.  Defaults to "/tmp/bundlerkey.enc"
-pubfile    Path and file for the public key used for encrypting the AES-GCM key.  Defaults to "/tmp/morpheus.pub"

-help    Prints this text.
Examples:
Generates a bundle with port scans, system stats, elasticsearch results and morpheus logs
   $ ./morpheus-fling -infile="/home/slimshady/network.txt"

Generates a bundle with no portscans in it at /tmp/bundler.zip
   $ ./morpheus-fling

Specify current directory for bundler and keyfile path
   $ go run morpheus_fling.go -bundler bundler.zip -keyfile bundlerkey.enc
`

type Results struct {
	ElasticStats     *elasticing.Esstats       `json:"es_stats"`
	ElasticIndices   []elasticing.Esindices    `json:"es_indices"`
	System           *sysinfo.SysInfo          `json:"system_stats"`
	Scans            []portscanner.ScanResult  `json:"port_scans,omitempty"`
	RabbitStatistics []rabbiting.RabbitResults `json:"rabbit_stats"`
	MorphLogs        string                    `json:"morpheus_logs"`
}



// FileWrtr takes content and an outfile and appends content to the outfile
func FileWrtr(content string, fileName string) {
		//Remove existing files
		if (fileExists(fileName)){
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

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func createBundle(){

		//Remove exiting bundle file
		if (fileExists(*bundlerPtr)){
			fmt.Println("Bundler File already exists. Replacing file")
			e := os.Remove(*bundlerPtr)
			if e != nil {
				log.Fatal(e)
			}
		}

		files := []string{
			*outfilePtr,
			*keyfilePtr,
		}
		// Bundle the whole shebang
		if err := archiver.Archive(files, *bundlerPtr); err != nil {
			log.Fatal(err)
		}
		
	
}

func extractBundle(){

	// Extract the encrypted bundle
	t := time.Now()
	timeStamp := t.Format("20060102150405")
	folderName := "extracted_"+timeStamp
	if err := archiver.Unarchive(*bundlerPtr,folderName+"/"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Extracting Bundle File")
	nonText, err := os.ReadFile(folderName+"/output.json")
	if err != nil {
		log.Fatal("Can't load output file", err)
	}

	nonKey, err := os.ReadFile(folderName+"/bundlerkey.enc")
	if err != nil {
		log.Fatal("Can't load key file", err)
	}

	decryptedText := encryptText.DecryptItAll(*privatekeyPtr, nonText,nonKey)
	//fmt.Println("Decrypted Text = ",decryptedText)
	FileWrtr(decryptedText, folderName+"/morpheus_log.json")
	
}

func checkHealth(){
	fmt.Println("Checking health status")

	if *infilePtr != "" {
	// Open our jsonFile
	jsonFile, err := os.Open(*infilePtr)
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

	type ESstats []struct {
	Epoch string `json:"epoch"`
	Timestamp string `json:"timestamp"`
	Cluster string `json:"cluster"`
	Status string `json:"status"`
	NodeTotal string `json:"node.total"`
	NodeData string `json:"node.data"`
	Shards string `json:"shards"`
	Pri string `json:"pri"`
	Relo string `json:"relo"`
	Init string `json:"init"`
	Unassign string `json:"unassign"`
	PendingTasks string `json:"pending_tasks"`
	MaxTaskWaitTime string `json:"max_task_wait_time"`
	ActiveShardsPercent string `json:"active_shards_percent"`
	}
	type R struct {
		ElasticStats     ESstats       `json:"es_stats"`
		
	}

	var es R

	json.Unmarshal(byteValue, &es)
	fmt.Println(es)
}
}



// Need to initialize the ini file and pass into another function to iterate?
func main() {

	flag.Usage = help
	flag.Parse()

	
	if *extractPtr{ //Extract 

		extractBundle()


	} else if *healthPtr { //check health from log files

		checkHealth()

	} else { // Encrypt and bundle log file

	

	

	// Initialize an empty ScanResult slice, omitted from result if empty
	var destArray []portscanner.ScanResult
	if *infilePtr != "" {
		psArray := filereader.FileToStructArray(*infilePtr, *uLimit)
		destArray = portscanner.Start(psArray, 500*time.Millisecond)
	}

	superSecrets := secparse.ParseSecrets(*secfilePtr)
	rmqpassword := superSecrets.Rabbitmq.MorpheusPassword

	// Gather system stats into a si array
	sysStats := sysgatherer.SysGather()

	// Gather elasticsearch health and indices into structs for results
	esHealth := elasticing.ElasticHealth()
	esIndices := elasticing.ElasticIndices()
	rabbitStuff := rabbiting.RabbitStats("morpheus", rmqpassword)

	morpheus, err := ioutil.ReadFile(*logfilePtr)
	if err != nil {
		log.Fatalf("Error reading public key file: %s", err)
	}

	// Create instance of results struct from packages returns
	results := Results{
		ElasticStats:     esHealth,
		ElasticIndices:   esIndices,
		System:           sysStats,
		Scans:            destArray,
		RabbitStatistics: rabbitStuff,
		MorphLogs:        string(morpheus),
	}

	resultjson, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}

	//fmt.Fprintf(os.Stdout, "%s", resultjson)
	//FileWrtr("\nULTIMATE:\n" + string(resultjson), *outfilePtr)

	// Base resultjson into Encryption package and write encrypted file and key
	nonSense := encryptText.EncryptItAll(*pubPtr, string(resultjson))
	nonText := nonSense.Ciphertext
	nonKey := nonSense.EncryptedKey
	_ = nonText
	_ = nonKey
	FileWrtr(string(nonText), *outfilePtr)
	FileWrtr(string(nonKey), *keyfilePtr)


	createBundle()
	
}

	
	
}

func help() {
	fmt.Println(helpText)
	os.Exit(1)
}
