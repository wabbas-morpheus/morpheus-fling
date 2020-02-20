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
)

const helpText = `morpheus-fling [options]
Options:
-infile     The source file for network port scanning.  If none is provided port scans will be skipped.
-outfile    The destination directory of the generated package, "output.txt" by default.
-ulimit     Ulimit of the system, defaults to 1024.
-logfile    Logfile to add to the bundle.  Defaults to "/var/log/morpheus/morpheus-ui/current".
-bundler    Path and file to bundle into.  Defaults to "/tmp/bundler.zip".

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
	//if *infilePtr != "" {
		FileWrtr("PORT SCANS:\n", *outfilePtr)
		psArray := filereader.FileToStructArray(*infilePtr, *uLimit)

		destArray := portscanner.Start(psArray, 500*time.Millisecond)
		thisisjson, err := json.MarshalIndent(destArray, "", " ")
		if err != nil {
			log.Fatal("Can't encode to JSON", err)
		}
		fmt.Println(destArray)
		fmt.Fprintf(os.Stdout, "%s", thisisjson)
		FileWrtr(string(thisisjson), *outfilePtr)


	sysStats := sysgatherer.SysGather()
	haveajson, err := json.MarshalIndent(sysStats, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}
	FileWrtr("\n\nOS STATS:\n" + string(haveajson), *outfilePtr)
	esStuff := elasticing.ElasticHealth()
	esIndices := elasticing.ElasticIndices()
	morejson, err := json.MarshalIndent(esStuff, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}
	yetmorejson, err := json.MarshalIndent(esIndices, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}
	FileWrtr("\n\nES STATS:\n" + string(morejson), *outfilePtr)
	FileWrtr("\n\nES INDICES:\n" + string(yetmorejson), *outfilePtr)
	if err := archiver.Archive([]string{*outfilePtr, *logfilePtr}, *bundlerPtr); err != nil {
		log.Fatal(err)
	}
	lolol := Results{
		ElasticStats:   esStuff,
		ElasticIndices: esIndices,
		System:         sysStats,
		Scans:          destArray,
	}

	ultimatejson, err := json.MarshalIndent(lolol, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}

	FileWrtr("\n\nULTIMATE:\n" + string(ultimatejson), *outfilePtr)
	nonSense := encryptText.EncryptItAll("/tmp/this.pub", string(ultimatejson))
	nonText := nonSense.Ciphertext
	nonKey := nonSense.EncryptedKey
	FileWrtr(string(nonText), *outfilePtr)
	FileWrtr(string(nonKey), *outfilePtr)
}

func help() {
	fmt.Println(helpText)
	os.Exit(1)
}
