package lib

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"github.com/tealeg/xlsx/v3"
)

func GetLogger()(*log.Logger){
	config := GetConfig()
	logFile, err := os.OpenFile(config.LOGPATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println(err)
	}
	logger := log.New(logFile, "", log.LstdFlags|log.Llongfile)
	return logger
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
func get_pre_name(name string) (string) {
	re, err := regexp.Compile("[\\d+$]")
	if err != nil{
		fmt.Println(err)
	}
	pre_name := re.ReplaceAllString(name, "")
	return pre_name
}

func Check_data(value string, rule map[string]string) (map[string]string){
	result := map[string]string{}
	//正则匹配敏感信息
	for rule_id, rule_re := range rule {
		re, err := regexp.Compile(rule_re)
		if err != nil{
			fmt.Println(err)
			logger.Println(err)
		}
		v_flag := re.FindString(value)
		if v_flag != "" {
			fmt.Println(rule_id, v_flag)
			result["rule"] = rule_id
			result["value"] = v_flag
		}
	}
	return result
}


func write_excel(sh *xlsx.Sheet,dbs []database){

	for _, db := range dbs{
		for _, table := range db.tables{
			for _, data := range table.data{
				tmp := []string{}
				tmp = append(tmp, db.instance_info.region)
				tmp = append(tmp, db.instance_info.id)
				tmp = append(tmp, db.instance_info.name)
				tmp = append(tmp, db.name)
				tmp = append(tmp, table.name)
				tmp = append(tmp, data.(map[string]string)["column"])
				tmp = append(tmp, data.(map[string]string)["flag"])
				tmp = append(tmp, data.(map[string]string)["rule"])
				tmp = append(tmp, data.(map[string]string)["value"])
				row := sh.AddRow()
				row.SetHeight(15)
				for _,i := range tmp{
					cell := row.AddCell()
					cell.Value = i
				}
			}
		}
	}
}

func GetConfig() (Config){
	var config Config
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return config
}

func Check(){
	config := GetConfig()
	instances := GetInstances()					//获取所有RDS实例id
	totalInstance := len(instances)
	fmt.Printf("共获取到 %d 个待检查实例\n", totalInstance)
	logger.Printf("共获取到 %d 个待检查实例", totalInstance)
	wb := xlsx.NewFile()
	sh, err := wb.AddSheet("Sheet")
	if err != nil {
		fmt.Println(err)
	}
	column_title := []string{"可用区", "RDS实例", "实例名", "数据库", "表名", "列名", "匹配依据", "隐私类别", "值"}
	row := sh.AddRow()
	row.SetHeight(15)
	for _,i := range column_title{
		cell := row.AddCell()
		cell.Value = i
	}
	var totalDB = 0
	for i:=0; i<len(instances); i++ {
		dbs := get_dbs(instances[i])				//获取每个RDS实例中的个人隐私字段
		if len(dbs) > 0 {
			write_excel(sh, dbs)			//将结果写入excel
		}else{
			totalInstance -= 1
		}
		totalDB += len(dbs)
	}
	wb.Save(config.EXCEL_NAME)
	fmt.Printf("检查完毕，共检查 %d 个实例，%d个数据库，结果保存至excel：%s", totalInstance, totalDB, config.EXCEL_NAME)
	logger.Printf("检查完毕，共检查 %d 个实例，%d个数据库，结果保存至excel：%s", totalInstance, totalDB, config.EXCEL_NAME)
}