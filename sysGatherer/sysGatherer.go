package sysgatherer

import (
	"log"
	"os/user"
	"os/exec"
    "bufio"
    "strings"
    "strconv"
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

	
}

func GetStorageUsed() int{

    //get path of the elasticsearch data  
    outESFile, err := exec.Command("grep","data","/opt/morpheus/embedded/elasticsearch/config/elasticsearch.yml").Output()

    if err != nil {
        log.Fatal(err)
    }


    outputES := string(outESFile[:])
    fPath := ""
    scannerES := bufio.NewScanner(strings.NewReader(outputES))
        for scannerES.Scan() { //iterate of each line
                line := strings.Fields(scannerES.Text())//convert line text in a list
                fPath = line[1] //get storage used info
                //fmt.Printf("fpath = %s\n",fPath)
                
                        
                
        }
        if err := scannerES.Err(); err != nil {
                log.Fatal(err)
        }

    //Get available storage info on elasticsearch mount point
    out, err := exec.Command("df","-h",fPath).Output()

    // if there is an error with our execution
    // handle it here
    if err != nil {
        log.Fatal(err)
    }
    // as the out variable defined above is of type []byte we need to convert
    // this to a string or else we will see garbage printed out in our console
    // this is how we convert it to a string
    output := string(out[:])
    used := ""
    scanner := bufio.NewScanner(strings.NewReader(output))
        for scanner.Scan() { //iterate of each line
                line := strings.Fields(scanner.Text())//convert line text in a list
                storageUsedPercent := line[4] //get storage used info
                used = storageUsedPercent[0:len(storageUsedPercent)-1]
                        
                
        }
        if err := scanner.Err(); err != nil {
                log.Fatal(err)
        }
        rtn, err :=  strconv.Atoi(used)//Convert to integer
return rtn
}
