package rbParse

import (
	"strings"
)

func GetApplianceInstallType(rbPtr string) string {
	installType := "All In One"
	rb := ParseRb(rbPtr)
	for k, v := range rb {
		//fmt.Printf("setting = %s value = %s\n", k, v)
		if k == "mysql[enable]" && strings.ToLower(v) == "false" {
			installType = "HA"
			//fmt.Println("Found setting")
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

func ExternalDB(rbPtr string) bool {
	externalDB := false
	rb := ParseRb(rbPtr)
	for k, v := range rb {
		//fmt.Printf("setting = %s value = %s\n", k, v)
		if k == "mysql[enable]" && strings.ToLower(v) == "false" {
			externalDB = true
			//fmt.Println("Found setting")
			break
		}
	}

	return externalDB

}

func ExternalRabbit(rbPtr string) bool {
	externalRabbit := false
	rb := ParseRb(rbPtr)
	for k, v := range rb {
		//fmt.Printf("setting = %s value = %s\n", k, v)
		if k == "rabbitmq[enable]" && strings.ToLower(v) == "false" {
			externalRabbit = true
			//fmt.Println("Found setting")
			break
		}
	}

	return externalRabbit

}

func ExternalElastic(rbPtr string) bool {
	externalElastic := false
	rb := ParseRb(rbPtr)
	for k, v := range rb {
		//fmt.Printf("setting = %s value = %s\n", k, v)
		if k == "elasticsearch[enable]" && strings.ToLower(v) == "false" {
			externalElastic = true
			//fmt.Println("Found setting")
			break
		}
	}

	return externalElastic

}
