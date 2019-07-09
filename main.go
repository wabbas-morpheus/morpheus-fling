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

	"github.com/zcalusic/sysinfo"
	"golang.org/x/sync/semaphore"
)

// PortScanner is the struct of the ip objects and weighted semaphore arguments for cpu arch
type PortScanner struct {
	ip   string
	port int
	lock *semaphore.Weighted
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
			fmt.Println(port, "closed")
			return ip + ":" + strconv.Itoa(port) + " closed\n"
		}
	}

	conn.Close()
	fmt.Println(port, "open")
	return ip + ":" + strconv.Itoa(port) + " open\n"
}

// Start is the method that runs the goroutines the ScanPort is called from
func Start(fn string, fo string, uLimit int64, timeout time.Duration) {

	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wg := sync.WaitGroup{}
	defer wg.Wait()

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
		wg.Add(1)
		ps.lock.Acquire(context.TODO(), 1)
		go func(porti int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			scanResult := ScanPort(ps.ip, ps.port, timeout)
			FileWrtr(scanResult, fo)
		}(ps.port)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
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
	flag.Parse()

	Start(*infilePtr, *outfilePtr, *uLimit, 500*time.Millisecond)
	sysStats := SysGather()
	FileWrtr(sysStats, *outfilePtr)

}
