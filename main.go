package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mholt/archiver"
	"github.com/zcalusic/sysinfo"
	"golang.org/x/sync/semaphore"
)

// PortScanner is the struct of the ip objects and weighted semaphore arguments for cpu arch
type PortScanner struct {
	ip   string
	port int
	lock *semaphore.Weighted
}

// ScanResult is the struct of ip objects and the status from running the scan
type ScanResult struct {
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Status string `json:"status"`
}

// ScanPort is the method that is called to test the port on target ips
func ScanPort(ip string, port int, timeout time.Duration) string {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		} else {
			fmt.Println(ip + ":" + strconv.Itoa(port) + " closed\n")
			return "closed"
		}
	}

	conn.Close()
	fmt.Println(ip + ":" + strconv.Itoa(port) + " open\n")
	return "open"
}

// FileToStructArray takes in a file with a list of ip:port and adds them to an array of structs
func FileToStructArray(fn string, uLimit int64) []*PortScanner {
	var psArray []*PortScanner
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
		ps := &PortScanner{
			ip:   parts[0],
			port: portInt,
			lock: semaphore.NewWeighted(uLimit),
		}
		psArray = append(psArray, ps)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return psArray
}

// Start is blah
func Start(psEntity []*PortScanner, timeout time.Duration) string {

	scanSlice := make([]ScanResult, len(psEntity))

	var wg sync.WaitGroup
	defer wg.Wait()
	for i, element := range psEntity {
		wg.Add(1)
		element.lock.Acquire(context.TODO(), 1)
		go func(i int, element *PortScanner) {
			defer element.lock.Release(1)
			defer wg.Done()
			scanOut := ScanPort(element.ip, element.port, timeout)
			sr := ScanResult{
				Ip:     element.ip,
				Port:   element.port,
				Status: scanOut,
			}
			scanSlice[i] = sr
		}(i, element)
	}

	scanJson, err := json.MarshalIndent(scanSlice, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}

	return string(scanJson)
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
	flag.Parse()

	FileWrtr("PORT SCANS:\n", *outfilePtr)
	psArray := FileToStructArray(*infilePtr, *uLimit)

	destArray := Start(psArray, 500*time.Millisecond)
	thisisjson, err := json.MarshalIndent(destArray, "", " ")
	if err != nil {
		log.Fatal("Can't encode to JSON", err)
	}

	fmt.Println(destArray)
	fmt.Fprintf(os.Stdout, "%s", thisisjson)
	FileWrtr(string(thisisjson), *outfilePtr)

	sysStats := SysGather()
	FileWrtr("\n\nOS STATS:\n"+sysStats, *outfilePtr)
	if err := archiver.Archive([]string{*logfilePtr}, "/tmp/bundler.zip"); err != nil {
		log.Fatal(err)
	}

}
