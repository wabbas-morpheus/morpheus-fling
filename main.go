package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/mholt/archiver"
	portscanner "github.com/tadamhicks/morpheus-fling/portScanner"
	"github.com/zcalusic/sysinfo"
	"golang.org/x/sync/semaphore"
)

// FileToStructArray takes in a file with a list of ip:port and adds them to an array of structs
func FileToStructArray(fn string, uLimit int64) []*portscanner.PortScanner {
	var psArray []*portscanner.PortScanner
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		parts := strings.Split(s, ":")
		portInt, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println(portInt)
		}
		ps := &portscanner.PortScanner{
			Ip:   parts[0],
			Port: portInt,
			Lock: semaphore.NewWeighted(uLimit),
		}
		psArray = append(psArray, ps)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return psArray
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
	psArray := FileToStructArray(*infilePtr, *uLimit)

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
