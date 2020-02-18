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

	filereader "github.com/gomorpheus/morpheus-fling/fileReader"
	portscanner "github.com/gomorpheus/morpheus-fling/portScanner"
	sysgatherer "github.com/gomorpheus/morpheus-fling/sysGatherer"
	encryptText "github.com/gomorpheus/morpheus-fling/encryptText"
	"github.com/mholt/archiver"
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

Generates a statik package only with the ".js" files
from the ./public directory.
   $ statik -include=*.js
`

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
	if *infilePtr != "" {
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

	}

	sysStats := sysgatherer.SysGather()
	FileWrtr("\n\nOS STATS:\n"+sysStats, *outfilePtr)
	nonSense := encryptText.EncryptItAll("/tmp/this.pub", sysStats)
	nonText := nonSense.Ciphertext
	nonKey := nonSense.EncryptedKey
	FileWrtr(string(nonText), *outfilePtr)
	FileWrtr(string(nonKey), *outfilePtr)
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

func help() {
	fmt.Println(helpText)
	os.Exit(1)
}
