package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"rds/src/aliyun"
	"rds/src/lib"
)

func testCase(str string){
	lib.Check_data(str, lib.COLUMN_RE)
	lib.Check_data(str, lib.VALUE_RE)
}

var (
	ifInit = flag.Bool("init", false, "从阿里云获取rds实例信息，如手动导入可忽略")
	run = flag.Bool("run", false, "开始检查")
	test = flag.String("test", "", "输入一个字符串，检查是否属于个人隐私数据")
	help = flag.Bool("h", false, "帮助信息")
)

func main() {
	flag.Parse()
	logger := lib.GetLogger()
	if *help == true {
		fmt.Printf("Usage: go run main.go [-h] [-init] [-run]")
		flag.PrintDefaults()
	}else if(*ifInit == true){
		logger.Println("初始化获取阿里云RDS实例信息")
		aliyun.Init()
	}else if(*test != ""){
		testCase(*test)
	}else if(*run == true){
		logger.Println("开始检查")
		lib.Check()
	}else{
		fmt.Printf("Usage: go run main.go [-h] [-init] [-run]")
		flag.PrintDefaults()
	}
}
