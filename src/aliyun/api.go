package aliyun

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"os"
	"rds/src/lib"
	"strings"
)

type Instance struct {
	id string
	name string
	region string
	cs string
}

var config = lib.GetConfig()
var logger = lib.GetLogger()

func getDatabaseByInstance(instance string, client *rds.Client) ([]string) {

	request := rds.CreateDescribeDatabasesRequest()
	request.Scheme = "https"

	request.DBInstanceId = instance

	response, err := client.DescribeDatabases(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	var db_response DBResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &db_response)
	if err != nil {
		fmt.Println(err)
	}
	dbs := []string{}
	for _,db := range db_response.Databases["Database"].([]interface{}) {
		dbName := db.(map[string]interface{})["DBName"].(string)
		preDbName := Get_pre_name(dbName)
		if !lib.CheckArray(dbs, preDbName) {
			dbs = append(dbs, preDbName)
		}
	}
	return dbs
}

func getInstanceConnectAddress(client *rds.Client, instance_id string) (string) {

	request := rds.CreateDescribeDBInstanceNetInfoRequest()
	request.Scheme = "https"

	request.DBInstanceId = instance_id

	response, err := client.DescribeDBInstanceNetInfo(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	var dbNetResponse DBNetResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &dbNetResponse)
	if err != nil {
		fmt.Println(err)
	}
	dbnet := dbNetResponse.DBInstanceNetInfos["DBInstanceNetInfo"].([]interface{})[0]
	connectionString := dbnet.(map[string]interface{})["ConnectionString"]
	port := dbnet.(map[string]interface{})["Port"]

	return fmt.Sprintf("%s:%s", connectionString, port)
}

func getInstanceByPage(pageNumber int, client *rds.Client) ([]Instance) {

	request := rds.CreateDescribeDBInstancesRequest()
	request.Scheme = "https"
	request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNumber))
	response, err := client.DescribeDBInstances(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	if strings.Contains(response.GetHttpContentString(), "SDK.ServerError") {
		fmt.Println(response.GetHttpContentString())
		os.Exit(0)
	}
	var rds_response RDSResponse
	err = json.Unmarshal(response.GetHttpContentBytes(), &rds_response)
	if err != nil {
		fmt.Println(err)
	}
	DBInstances := rds_response.Items.(map[string]interface{})["DBInstance"].([]interface{})
	instances := []Instance{}
	for _,dbinstance := range DBInstances{
		tmp_instance := Instance{
			dbinstance.(map[string]interface{})["DBInstanceId"].(string),
			dbinstance.(map[string]interface{})["DBInstanceDescription"].(string),
			dbinstance.(map[string]interface{})["RegionId"].(string),
			getInstanceConnectAddress(client, dbinstance.(map[string]interface{})["DBInstanceId"].(string)),
		}
		instances = append(instances, tmp_instance)
	}
	return instances
}

func Init() {		//获取阿里云rds实例信息和rds连接串，写入文件中
	f, err := os.OpenFile(config.RDS_INFO_FILEPATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	if config.ACCESSKEYID == "" || config.ACCESSKEYSECRET == "" {
		fmt.Println("请在配置文件 config.yml 中填写 accessKey")
		logger.Println("请在配置文件中 config.yml 填写 accessKey")
		os.Exit(1)
	}
	for _,region := range config.REGIONS{
		client, err := rds.NewClientWithAccessKey(region, config.ACCESSKEYID, config.ACCESSKEYSECRET)
		if err != nil {
			fmt.Print(err.Error())
		}
		request := rds.CreateDescribeDBInstancesRequest()
		request.Scheme = "https"
		pageSize := 30
		request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))
		response, err := client.DescribeDBInstances(request)
		if err != nil {
			fmt.Print(err.Error())
		}
		if strings.Contains(response.GetHttpContentString(), "SDK.ServerError") {
			fmt.Println(response.GetHttpContentString())
			os.Exit(0)
		}

		var rds_response RDSResponse
		err = json.Unmarshal(response.GetHttpContentBytes(), &rds_response)
		if err != nil {
			fmt.Println(err)
		}
		totalPageCount := (rds_response.TotalRecordCount / pageSize) + 1
		for i:=1; i<=totalPageCount; i++ {
			fmt.Printf("%s 区共有 %d 个rds，正在获取第 %d 页 rds, 一共有 %d 页\n", region, rds_response.TotalRecordCount, i, totalPageCount)
			instances := getInstanceByPage(i, client)
			//将rds基础信息写入文件
			for _,instance := range instances{
				str := fmt.Sprintf("%s	%s	%s	%s\n", instance.id, instance.name, instance.region, instance.cs)
				fmt.Print(str)
				if _, err := f.WriteString(str); err != nil {
					if err := f.Close(); err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
	fmt.Println("获取 rds 实例信息完成。")

}