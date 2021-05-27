package lib

import (
	"database/sql"
	"fmt"
	"strings"
)

func get_table_data(db *sql.DB, table_name string)([]interface{}) {		//检查table_name中的所有字段，以及取LIMIT_COUNT条数据，看是否包含个人隐私数据
	config := GetConfig()
	var result []interface{}

	sqlData := fmt.Sprintf("select * from `%s` limit %d", table_name, config.LIMIT_COUNT)
	rows, err := db.Query(sqlData)				//查询数据
	if err != nil{
		logger.Println(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logger.Println(err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]  							//二维数组
	}
	var value string
	column_flag := true										//column标志位，只检查一次
	var value_flag []bool									//value标志位，20条数据若第一条检查出来，后面的19条不再检查。
	for i:=0; i<len(columns); i++ {
		value_flag = append(value_flag, false)
	}
	for rows.Next() {										//遍历每一行数据
		err = rows.Scan(scanArgs...)
		if err != nil {
			logger.Println(err)
		}
		for i, col := range values {						//遍历每一列数据
			if col == nil {
				value = "NULL"
			}else{
				value = string(col)
			}
			if column_flag {
				result_column := Check_data(columns[i], COLUMN_RE)			//检查字段名
				if len(result_column) > 0 {
					result_column["flag"] = "Column"						//若匹配到，则设置标志位为 Column
					result_column["column"] = columns[i]
					result_column["value"] = value
					result = append(result, result_column)
					value_flag[i] = true
				}else{
					fmt.Println("未匹配到 column：", columns[i])
				}

			}
			if !value_flag[i] {													//limit的条数均进行检查，若匹配到一条则后面的value不再检查
				result_value := Check_data(value, VALUE_RE)						//检查字段value
				if len(result_value) > 0 {
					result_column := Check_data(columns[i], COLUMN_WHITE_LIST)	//检查列名是否在白名单中
					if len(result_column) > 0 {
						value_flag[i] = true
						continue
					}
					result_value["flag"] = "Value"								//若匹配到则设置标志位为 Value
					result_value["column"] = columns[i]
					value_flag[i] = true										//value标志位设为true，同列的value不再匹配
					result = append(result, result_value)
				}else{
					//fmt.Println("未匹配到 value： %", value[:100])
				}
			}
		}
		column_flag = false			//column只检查一次，置为false后面的循环不再检查
	}
	return result
}

func get_tables(instance_info instance, db_name string)([]table) {
	config := GetConfig()
	instance_id := instance_info.id
	instance_cs := instance_info.cs
	var connStr string
	connStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=5s&readTimeout=%ds", config.USER, config.PASSWD, instance_cs, db_name, config.TIMEOUT)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		logger.Println(err)
	}
	defer db.Close()

	sqlTxt := "show tables"
	rows, err := db.Query(sqlTxt)
	if err != nil {
		logger.Println(err)
	}
	var table_name string
	var tables = []table{}
	var table_pre_name = make(map[string]struct{})
	i := 0
	for rows.Next() {			//遍历所有表
		i++
		_ = rows.Scan(&table_name)
		pre_name := get_pre_name(table_name)
		if _, ok := table_pre_name[pre_name]; !ok{
			table_pre_name[pre_name] = struct{}{}
			fmt.Printf("正在检查 RDS 实例: %s, 数据库: %s, 表: %s\n", instance_id, db_name, table_name)
			logger.Printf("正在检查 RDS 实例: %s, 数据库: %s, 表: %s\n", instance_id, db_name, table_name)
			results := get_table_data(db, table_name)
			if len(results) > 0 {
				//for _, result := range results{
				//	fmt.Println(db_name, table_name, result.(map[string]string)["flag"], result.(map[string]string)["rule"], result.(map[string]string)["value"], )
				//}
				tmp_table := table {
					table_name,
					results,
				}
				tables = append(tables, tmp_table)
			}
		}
	}
	return tables
}

//根据rds实例查询数据库
func get_dbs(instance_info instance)( []database ) {
	config := GetConfig()
	instanceConnectStrings := instance_info.cs
	sqlTxt := "show databases"
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/?timeout=5s&readTimeout=%ds", config.USER, config.PASSWD, instanceConnectStrings, config.TIMEOUT)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		errs := fmt.Sprintf("Open mysql error: %s", err)
		logger.Println(errs)
		if strings.Contains(errs, "i/o timeout") {
			fmt.Printf("Instance %s connect timeout", instance_info.id)
			logger.Println(errs, connStr)
		}
		return []database{}
	}
	defer db.Close()

	rows, err := db.Query(sqlTxt)
	if err != nil {
		fmt.Println(err)
		logger.Println(err)
	}
	var db_name string
	var dbs []database
	db_pre_name := make(map[string]struct{})
	for rows.Next() {
		_ = rows.Scan(&db_name)
		//if _, isWhiteDb := db_white_list[db_name]; isWhiteDb{
		//	continue
		//}
		if CheckArray(config.DB_WHITE_LIST, db_name) {			//如果数据库在白名单中则跳过
			continue
		}
		pre_name := get_pre_name(db_name)

		if _, ok := db_pre_name[pre_name]; ok {					//判断数据库前缀是否出现过，出现则跳过
			continue
		}
		db_pre_name[pre_name] = struct{}{}						//未出现则将pre_name加入字典中
		tables := get_tables(instance_info, db_name)
		_database := database {
			instance_info,
			db_name,
			tables,
		}
		dbs = append(dbs, _database)
	}
	return dbs
}