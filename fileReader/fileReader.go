package filereader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	portscanner "github.com/gomorpheus/morpheus-fling/portScanner"
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
