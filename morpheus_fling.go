package main

import (
	"encoding/json"
	"flag"
	"fmt"
	elasticing "github.com/wabbas-morpheus/morpheus-fling/elasticIng"
	rabbiting "github.com/wabbas-morpheus/morpheus-fling/rabbitIng"
	secparse "github.com/wabbas-morpheus/morpheus-fling/secParse"
	"log"
	"morpheus-fling/rbParse"
	"os"
	"path"
	"time"

	"github.com/mholt/archiver"
	encryptText "github.com/wabbas-morpheus/morpheus-fling/encryptText"
	filereader "github.com/wabbas-morpheus/morpheus-fling/fileReader"
	portscanner "github.com/wabbas-morpheus/morpheus-fling/portScanner"
	//sysgatherer "github.com/wabbas-morpheus/morpheus-fling/sysGatherer"
	//"github.com/zcalusic/sysinfo"
	"morpheus-fling/healthCheck"
)

var (
	defaultPath = "."
	infilePtr   = flag.String("infile", "", "a string")
	secfilePtr  = flag.String("secfile", "/etc/morpheus/morpheus-secrets.json", "a string")
	outfilePtr  = flag.String("outfile", path.Join(".", "encrypted_logs.json"), "a string")
	uLimit      = flag.Int64("ulimit", 1024, "an integer")
	logfilePtr  = flag.String("logfile", "/var/log/morpheus/morpheus-ui/current", "a string")
	bundlerPtr  = flag.String("bundler", path.Join(defaultPath, ""), "a string")
	keyfilePtr  = flag.String("keyfile", "/tmp/bundlerkey.enc", "a string")
	pubPtr      = flag.String("pubkey", path.Join(defaultPath, "morpheus.pub"), "a string")
	//privatekeyPtr    = flag.String("privkey", path.Join(defaultPath, "morpheus.pem"), "a string")
	//extractPtr       = flag.Bool("extract", false, "a bool")
	healthPtr        = flag.Bool("health", false, "a bool")
	flingsettingsPtr = flag.String("token", "/etc/morpheus/morpheus-fling-settings.json", "a string")
	rbfilePtr        = flag.String("rbfile", "/etc/morpheus/morpheus2.rb", "a string")
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
	ElasticStats    *elasticing.Esstats             `json:"es_stats"`
	ElasticIndices  []elasticing.Esindices          `json:"es_indices"`
	ElasticSettings *elasticing.ESWaterMarkSettings `json:"es_settings"`
	//System           *sysinfo.SysInfo                `json:"system_stats"`
	Scans            []portscanner.ScanResult  `json:"port_scans,omitempty"`
	RabbitStatistics []rabbiting.RabbitResults `json:"rabbit_stats"`
	MorphLogs        string                    `json:"morpheus_logs"`
	MorphRB          string                    `json:"morpheus_rb"`
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createBundle() {
	var bundleName = ""
	files := []string{
		*outfilePtr,
		*keyfilePtr,
	}
	if *bundlerPtr != "." { //If filename specified use that
		//Remove exiting bundle file
		if fileExists(*bundlerPtr) {
			fmt.Println("Bundler File already exists. Replacing file")
			e := os.Remove(*bundlerPtr)
			if e != nil {
				log.Fatal(e)
			}
		}
		bundleName = *bundlerPtr
	} else { //Otherwise assign hostname to bundle file
		t := time.Now()
		timeStamp := t.Format("20060102150405")
		bundleName = "bundle_" + getHostName() + "_" + timeStamp + ".zip"
		//fmt.Printf("bundleName = %s", bundleName)

	}
	// Bundle the whole shebang
	if err := archiver.Archive(files, bundleName); err != nil {
		log.Fatal(err)
	}

}

func getHostName() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	//fmt.Println("hostname:", name)
	return name
}

func runHealthCheck() {
	fmt.Println("Checking health status")
	healthCheck.CheckHealth(*flingsettingsPtr)

	//fmt.Printf("Install Type = %s\n", rbParse.GetApplianceInstallType(*rbfilePtr))
	//fmt.Printf("Total DB Nodes = %d\n", rbParse.GetTotalNumberOfDBNodes(*rbfilePtr))

}

// Need to initialize the ini file and pass into another function to iterate?
func main() {

	flag.Usage = help
	flag.Parse()

	if *healthPtr { //check health from log files

		runHealthCheck()

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
		//sysStats := sysgatherer.SysGather()

		// Gather elasticsearch health and indices into structs for results
		esHealth := elasticing.ElasticHealth()
		esIndices := elasticing.ElasticIndices()
		esWaterMarkSettings := elasticing.ElasticWatermarkSettings()

		rabbitStuff := rabbiting.RabbitStats("morpheus", rmqpassword)

		morpheus, err := os.ReadFile(*logfilePtr)
		if err != nil {
			log.Fatalf("Error reading morpheus current log file key file: %s", err)
		}

		morpheusRb := rbParse.GetMorpheusRBFile(*rbfilePtr)

		// Create instance of results struct from packages returns
		results := Results{
			ElasticStats:    esHealth,
			ElasticIndices:  esIndices,
			ElasticSettings: esWaterMarkSettings,
			//System:           sysStats,
			Scans:            destArray,
			RabbitStatistics: rabbitStuff,
			MorphLogs:        string(morpheus),
			MorphRB:          morpheusRb,
		}

		//fmt.Printf("%+v", results.MorphLogs)

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
