package rbParse

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getChar(str string, index int) rune {
	return []rune(str)[index]
}

func ParseRb(rbfilePtr string) {
	//rbSettings := make(map[string]string)
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
		//fmt.Println(sc.Text()) // GET the line string

		rbLine := strings.Trim(sc.Text(), " ")
		var firstChar string = ""
		if rbLine != "" {
			firstChar = string(getChar(rbLine, 0))
			if firstChar != "#" { //skip comment line
				rbLine = strings.ReplaceAll(rbLine, "=>", "-")
				rbLine = strings.ReplaceAll(rbLine, "{", "")
				rbLine = strings.ReplaceAll(rbLine, "}", "")
				s := strings.Split(rbLine, "=")
				if len(s) == 2 {
					fmt.Printf("s key = %s value = %s\n", s[0], strings.ReplaceAll(s[1], "'", ""))
				}
			}

			if strings.Count(rbLine, "appliance_url") == 1 && strings.Count(rbLine, "=") == 0 && firstChar != "#" {
				s := strings.Split(rbLine, "'")
				if len(s) == 3 {
					fmt.Printf("s key = %s value = %s\n", s[0], strings.ReplaceAll(s[1], "'", ""))
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}

}
