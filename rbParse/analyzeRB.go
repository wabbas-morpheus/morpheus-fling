package rbParse

import (
	"fmt"
	"strings"
)

func GetApplianceInstallType(rbPtr string) string {
	installType := "All In One"
	rb := ParseRb(rbPtr)
	for k, v := range rb {
		//fmt.Printf("setting = %s value = %s\n", k, v)
		if k == "mysql[enable]" && v == "true" {
			installType = "HA"
			fmt.Println("Found setting")
			break
		}
	}

	return installType
}

func GetTotalNumberOfDBNodes(rbPtr string) int {
	totalNodes := 1
	rb := ParseRb(rbPtr)
	for k, v := range rb {
		//fmt.Printf("setting = %s value = %s\n", k, v)
		if k == "mysql['host']" {
			totalNodes = len(strings.Split(v, ","))
			break
		}
	}

	return totalNodes
}
