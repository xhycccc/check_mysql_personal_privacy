package aliyun

import (
	"fmt"
	"regexp"
)

type RDSResponse struct{
	TotalRecordCount int
	PageRecordCount int
	RequestId string
	PageNumber int
	Items interface{}
}

type DBResponse struct {
	Databases map[string]interface{}
	RequestId string
}

type DBNetResponse struct {
	RequestId string
	DBInstanceNetInfos map[string]interface{}
	SecurityIPMode string
	InstanceNetworkType string
}

func printArray(array []string){
	for _,i := range array {
		fmt.Println(i)
	}
}

func CheckArray(array []string, key string) bool {
	for _,str := range array {
		if key == str{
			return true
		}
	}
	return false
}

//取字符串除去尾部的数字的部分作为前缀
func Get_pre_name(name string) (string) {
	re, err := regexp.Compile("\\d+$")
	if err != nil{
		fmt.Println(err)
	}
	pre_name := re.ReplaceAllString(name, "")
	return pre_name
}
