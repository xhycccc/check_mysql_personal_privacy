## 背景

数据安全项目，检查mysql数据库中的个人信息敏感字段。

## 原理说明
1. 连接数据库实例
2. 执行sql语句"show databases"，获取到数据库名
3. 遍历数据库并连接，对数据库名进行检查，多个库名前缀相同（分库）只取一个
4. 执行sql语句"show tables"，获取表名
5. 遍历表，对表名进行检查，多个表前缀相同（分表）只取一个。
6. 执行sql语句`select count(*) from $table`，获取表的行数。若行数小`MIN_ROW`则不进行检查。
7. 执行sql语句`select * from $table limit $LIMIT_COUNT`，获取`LIMIT_COUNT`行数据进行检查。
8. 遍历列名和数据，使用正则逐一匹配。
9. 返回匹配结果，写入excel

## 配置文件
1. 在instance_info.txt中填写数据库实例相关信息，分别为实例ID、实例名、所在区域、连接地址串，用于程序获取数据库连接地址以及生成报表。格式为（字段之间以tab分隔）：
    ```
    test_instance_id  test_instance_nam test_region 127.0.0.1:3306
    127.0.0.1 test  test  127.0.0.1:3307
    ```
   如检查阿里云上的rds，可参考第四步自动获取。

2. 在check.txt中填写要检查的rds或mysql实例的id，一行一个，程序根据这个文件中的id去instance_info.txt中获取对应的数据库连接串。
    ```
    127.0.0.1
    127.0.0.2
    ```

3. 在配置文件config.yml填写数据库账号密码
    ```
   USER: 数据库账号
   PASSWD: 数据库密码 # 要求检查的数据库要有统一的账号密码
   ```
4. 运行命令
    ```
    Usage: go run main.go [-h] [-init] [-run]  -h   帮助信息
      -init
            从阿里云获取rds实例信息，如手动填写可忽略
      -run
            开始检查
      -test string
            输入一个字符串，检查是否属于个人隐私数据

    ```
   如果有授权过的accessKey，可以使用-init参数可从阿里云自动获取该账号下的所有rds信息。
   ```
   go run main.go -init
   ```
   开始检查：
   ```
   go run main.go -run
   ```
