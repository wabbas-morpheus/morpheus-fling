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

func cleanRBLine(str string) string {
	str = strings.Trim(str, " ")
	str = strings.ReplaceAll(str, "=>", "-")
	str = strings.ReplaceAll(str, "{", "")
	str = strings.ReplaceAll(str, "}", "")
	str = strings.ReplaceAll(str, "'", "")
	str = strings.ReplaceAll(str, "\"", "")
	return str
}

func GetMorpheusRBFile(rbfilePtr string) string {
	rbLine := ""
	morpheusRBFile, err := os.Open(rbfilePtr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened Morpheus rb file @ %s\n", rbfilePtr)
	// defer the closing of our jsonFile so that we can parse it later on
	defer func(morpheusRBFile *os.File) {
		err := morpheusRBFile.Close()
		if err != nil {
			log.Fatalf("unable to open rb file: %v", err)
		}
	}(morpheusRBFile)

	sc := bufio.NewScanner(morpheusRBFile)
	for sc.Scan() {
		foundPassword := strings.Contains(sc.Text(), "password")
		if foundPassword {
			fmt.Printf("Password - %s\n", sc.Text())
			s := strings.Split(sc.Text(), "=")
			rbLine = rbLine + s[0] + " = " + "Password Redacted"
		} else {
			rbLine = rbLine + sc.Text()
		}

	}
	return rbLine
}

func ParseRb(rbfilePtr string) map[string]string {
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

	sc := bufio.NewScanner(morpheusRBFile)
	for sc.Scan() {

		rbLine := cleanRBLine(sc.Text()) //Remove unwanted characters from rb line
		var firstChar string = ""
		if rbLine != "" {
			firstChar = string(getChar(rbLine, 0))
			if firstChar != "#" { //skip comment line
				i := strings.Index(rbLine, "#")
				if i > 0 { //Remove inline comments
					rbLine = rbLine[0:i]
				}

				s := strings.Split(rbLine, "=") //Get setting name and value
				if len(s) == 2 {
					//fmt.Printf("s key = %s value = %s\n", s[0], s[1])
					rbSettings[strings.Trim(s[0], " ")] = strings.Trim(s[1], " ")
				}
			}
			//obtain appliance url setting
			if strings.Count(rbLine, "appliance_url") == 1 && strings.Count(rbLine, "=") == 0 && firstChar != "#" {
				s := strings.Split(rbLine, " ")
				//fmt.Printf("s key = %s value = %s\n", s[0], s[1])
				rbSettings[strings.Trim(s[0], " ")] = strings.Trim(s[1], " ")
			}
		}
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
	}

	return rbSettings

}
