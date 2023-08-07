package rbParse

import (
	"fmt"
)

func GetApplianceInstallType(string rbPtr) string {
	installType := "AIO"
	rb := ParseRb(rbPtr)
	for k, v := range rb {
		fmt.Printf("setting = %s value = %s\n", k, v)
		if k == "mysql['enable']" && v == "true" {
			installType = "HA"
		}
	}

	return installType
}
