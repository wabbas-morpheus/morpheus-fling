package main

import (
	"encoding/json"
	"flag"
	"fmt"
	elasticing "github.com/gomorpheus/morpheus-fling/elasticIng"
	"log"
	"os"
	"path"
	"time"

	encryptText "github.com/gomorpheus/morpheus-fling/encryptText"
	filereader "github.com/gomorpheus/morpheus-fling/fileReader"
	portscanner "github.com/gomorpheus/morpheus-fling/portScanner"
	sysgatherer "github.com/gomorpheus/morpheus-fling/sysGatherer"
	"github.com/mholt/archiver"
	"github.com/zcalusic/sysinfo"
)


var (
	infilePtr = flag.String("infile", "", "a string")
	outfilePtr = flag.String("outfile", path.Join(".", "output.txt"), "a string")
	uLimit = flag.Int64("ulimit", 1024, "an integer")
	logfilePtr = flag.String("logfile", "/var/log/morpheus/morpheus-ui/current", "a string")
	bundlerPtr = flag.String("bundler", "/tmp/bundler.zip", "a string")
	keyfilePtr = flag.String("keyfile", "/tmp/bundlerkey.enc", "a string")
)

const helpText = `morpheus-fling [options]
Options:
-infile     The source file for network port scanning.  If none is provided port scans will be skipped.
-outfile    The destination directory of the generated package, "output.txt" by default.
-ulimit     Ulimit of the system, defaults to 1024.
-logfile    Logfile to add to the bundle.  Defaults to "/var/log/morpheus/morpheus-ui/current".
-bundler    Path and file to bundle into.  Defaults to "/tmp/bundler.zip".
-keyfile    Path and file to put the public key encrypted AES-GCM key into.  Defaults to "/tmp/bundlerkey.enc"

-help    Prints this text.
Examples:
Generates a bundle with port scans, system stats, elasticsearch results and morpheus logs
   $ ./morpheus-fling -infile="/home/slimshady/network.txt"

Generates a bundle with no portscans in it at /tmp/bundler.zip
   $ ./morpheus-fling
`

type Results struct {
	ElasticStats	*elasticing.Esstats	`json:"es_stats"`
	ElasticIndices	[]elasticing.Esindices	`json:"es_indices"`
	System	*sysinfo.SysInfo	`json:"system_stats"`
	Scans 	[]portscanner.ScanResult	`json:"port_scans,omitempty"`
}

// FileWrtr takes content and an outfile and appends content to the outfile
func FileWrtr(content string, fileName string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		log.Println(err)
	}
}

// Need to initialize the ini file and pass into another function to iterate?
func main() {

	flag.Usage = help
	flag.Parse()

	// Initialize an empty ScanResult slice, omitted from result if empty
	var destArray []portscanner.ScanResult
	if *infilePtr != "" {
		psArray := filereader.FileToStructArray(*infilePtr, *uLimit)
		destArray = portscanner.Start(psArray, 500*time.Millisecond)
	}


	// Gather system stats into a si array
	sysStats := sysgatherer.SysGather()

	// Gather elasticsearch health and indices into structs for results
	esHealth := elasticing.ElasticHealth()
	esIndices := elasticing.ElasticIndices()

	// Create instance of results struct from packages returns
	results := Results{
		ElasticStats:   esHealth,
		ElasticIndices: esIndices,
		System:         sysStats,
		Scans:          destArray,
	}

	resultjson, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}

	fmt.Fprintf(os.Stdout, "%s", resultjson)
	//FileWrtr("\nULTIMATE:\n" + string(resultjson), *outfilePtr)

	// Base resultjson into Encryption package and write encrypted file and key
	nonSense := encryptText.EncryptItAll("/tmp/this.pub", string(resultjson))
	nonText := nonSense.Ciphertext
	nonKey := nonSense.EncryptedKey
	FileWrtr(string(nonText), *outfilePtr)
	FileWrtr(string(nonKey), *keyfilePtr)

	// Bundle the whole shebang
	if err := archiver.Archive([]string{*outfilePtr, *keyfilePtr, *logfilePtr}, *bundlerPtr); err != nil {
		log.Fatal(err)
	}
}

func help() {
	fmt.Println(helpText)
	os.Exit(1)
}
