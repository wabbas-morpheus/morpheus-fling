package sysgatherer

import (
	"github.com/zcalusic/sysinfo"
	"log"
	"os/user"
)

// SysGather gathers system statistics and returns them as a string
func SysGather() *sysinfo.SysInfo {
	current, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if current.Uid != "0" {
		log.Fatal("requires superuser privilege")
	}

	var si sysinfo.SysInfo

	si.GetSysInfo()

	return &si

}
