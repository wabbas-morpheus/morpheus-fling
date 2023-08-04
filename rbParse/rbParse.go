package rbParse

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func ParseRb(rbfilePtr string) {
	rbSettings := make(map[string]string)
	// Open our rbfile
	morpheusRBFile, err := os.Open(rbfilePtr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened Morpheus rb file @ %s\n", rbfilePtr)
	// defer the closing of our jsonFile so that we can parse it later on
	defer func(morpheusRBFile *os.File) {
		err := morpheusRBFile.Close()
		if err != nil {
			log.Fatalf("defer file error: %v", err)
		}
	}(morpheusRBFile)

	// byteValue, _ := ioutil.ReadAll(morpheusRBFile)

	sc := bufio.NewScanner(morpheusRBFile)
	for sc.Scan() {
		fmt.Println(sc.Text()) // GET the line string
		rbLine := strings.Trim(sc.Text(), " ")
		s := strings.Split(rbLine, "'")
		fmt.Printf("s key = %s value = %s", s[0], s[1])
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}

}
