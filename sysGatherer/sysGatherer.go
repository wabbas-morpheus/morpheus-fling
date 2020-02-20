package sysgatherer

import (
	"log"
	"os/user"

	"github.com/zcalusic/sysinfo"
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

	//data, err := json.MarshalIndent(&si, "", "  ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(string(data))
	//return string(data)
}
