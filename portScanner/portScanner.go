package portscanner

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// PortScanner is the struct of the ip objects and weighted semaphore arguments for cpu arch
type PortScanner struct {
	Ip   string
	Port int
	Lock *semaphore.Weighted
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

// Start is blah
func Start(psEntity []*PortScanner, timeout time.Duration) []ScanResult {

	scanSlice := make([]ScanResult, len(psEntity))

	var wg sync.WaitGroup
	defer wg.Wait()
	for i, element := range psEntity {
		wg.Add(1)
		element.Lock.Acquire(context.TODO(), 1)
		go func(i int, element *PortScanner) {
			defer element.Lock.Release(1)
			defer wg.Done()
			scanOut := ScanPort(element.Ip, element.Port, timeout)
			sr := ScanResult{
				Ip:     element.Ip,
				Port:   element.Port,
				Status: scanOut,
			}
			scanSlice[i] = sr
		}(i, element)
	}

	return scanSlice
}
