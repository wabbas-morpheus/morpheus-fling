package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/mholt/archiver"
	filereader "github.com/tadamhicks/morpheus-fling/fileReader"
	portscanner "github.com/tadamhicks/morpheus-fling/portScanner"
	"github.com/zcalusic/sysinfo"
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

// SysGather gathers system statistics and returns them as a string
func SysGather() string {
	current, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if current.Uid != "0" {
		log.Fatal("requires superuser privilege")
	}

	var si sysinfo.SysInfo

	si.GetSysInfo()

	data, err := json.MarshalIndent(&si, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
	return string(data)
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

	sysStats := SysGather()
	FileWrtr("\n\nOS STATS:\n"+sysStats, *outfilePtr)
	if err := archiver.Archive([]string{*outfilePtr, *logfilePtr}, *bundlerPtr); err != nil {
		log.Fatal(err)
	}

}
