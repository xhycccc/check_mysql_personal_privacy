package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var config = GetConfig()
var instanceInfo = GetInstancesInfo()

func GetInstancesInfo() (map[string]instance) {
	file, err := os.Open(config.RDS_INFO_FILEPATH)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	var rds = make(map[string]instance)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		arr := strings.Split(sc.Text(), "\t")
		var tmp_instance instance
		if len(arr) > 0 {
			if len(arr) == 4 {
				tmp_instance = instance{
					arr[0],
					arr[1],
					arr[2],
					arr[3],
				}
			} else{
				tmp_instance = instance{
					arr[0],
					"",
					"",
					"",
				}
			}
			rds[tmp_instance.id] = tmp_instance
		}
	}
	return rds
}

//从RDS_FILEPATH文件中读取rds
func GetInstances()( []instance ) {
	instances := []instance{}
	file, err := os.Open(config.RDS_FILEPATH)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		rds := sc.Text()
		if len(rds) == 0 {
			continue
		}
		instance := instance{
			rds,
			instanceInfo[rds].name,
			instanceInfo[rds].region,
			instanceInfo[rds].cs,
		}
		instances = append(instances, instance)
	}
	return instances
}