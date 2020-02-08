package main

import (
	"encoding/json"
	"flag"
	"fmt"
	elasticing "github.com/gomorpheus/morpheus-fling/elasticIng"
	"log"
	"os"
	"time"

	filereader "github.com/gomorpheus/morpheus-fling/fileReader"
	portscanner "github.com/gomorpheus/morpheus-fling/portScanner"
	sysgatherer "github.com/gomorpheus/morpheus-fling/sysGatherer"
	encryptText "github.com/gomorpheus/morpheus-fling/encryptText"
	"github.com/mholt/archiver"
)

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

	infilePtr := flag.String("infile", "network.txt", "a string")
	outfilePtr := flag.String("outfile", "output.txt", "a string")
	uLimit := flag.Int64("ulimit", 1024, "an integer")
	logfilePtr := flag.String("logfile", "/var/log/morpheus/morpheus-ui/current", "a string")
	bundlerPtr := flag.String("bundler", "/tmp/bundler.zip", "a string")
	flag.Parse()

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
	FileWrtr("\n\nOS STATS:\n"+sysStats, *outfilePtr)
	esStuff := elasticing.ElasticHealth()
	morejson, err := json.MarshalIndent(esStuff, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}
	FileWrtr("\n\nES STATS:\n" + string(morejson), *outfilePtr)
	if err := archiver.Archive([]string{*outfilePtr, *logfilePtr}, *bundlerPtr); err != nil {
		log.Fatal(err)
	}
}
