package lib

var (
	logger = GetLogger()
	COLUMN_WHITE_LIST = map[string]string{						//column白名单，不会被匹配到结果中
		"Timestamp": ".*time$",									//timestamp会和和其它正则冲突
		"Hash": "(?i).*cid$|.*rid$|.*biz_id$|.*gid$|.*tid$|.*url$|^webhook$|.*_id$|^img$|.*icon$|^images$|.*file$",		//32位hash会和passwd正则冲突
	}
	COLUMN_RE = map[string]string{								//根据列名匹配
		"真实姓名": "(?i)^truename$|^consignee$|^name_cn$|^user_real_name$|^person_name$",
		"密码": "(?i).*password.*|.*psw.*|.*pwd.*|.*passwd.*",
		"email": "(?i).*mail.*",
		"电话号码": "(?i).*phone.*|.*mobile.*",
		"性别": "^男$|^女$",
		"身份证号": "(?i).*idcard.*|^credit_number$|^person_id_no$",
		"银行卡号": "(?i)^cash_account$|^regCardno$",
		"QQ/微信": "(?i)^qq$|^weixin$",
		"登录账号": "(?i)^account$|^account_name$|^user_show$|^username$|^user_id$",
		"生日": "(?i)^birthday$",
		"地址": "(?i).*addr.*",
		"device_id": "(?i)^device_id$|^devid$|^deviceid$",
		"device_mac": "(?i)^device_mac$|^devicemac$",
		"device_imei": "(?i)^device_imei$|^deviceimei$",
		"device_sn": "(?i)^device_sn$|^devicesn$",
		"android_id": "(?i)^android_id$|^androidid$",
		"问题答案": "(?i)answer.*",
	}
	VALUE_RE = map[string]string{
		"密码": "[0-9a-fA-F]{32}",
		"email": "\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*",
		"电话号码": "^1([358][0-9]|4[579]|66|7[0135678]|9[89])[0-9]{8}$",
		"性别": "^男$|^女$",
		"身份证号": "^[1-9]\\d{5}(19|20)\\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]",
		"银行卡号": "^([1-9]{1})(\\d{14}|\\d{18})$",
		"device_mac": "^[A-F0-9]{2}([-:]?[A-F0-9]{2})([-:.]?[A-F0-9]{2})([-:]?[A-F0-9]{2})([-:.]?[A-F0-9]{2})([-:]?[A-F0-9]{2})$",
	}
)

type instance struct {			//实例
	id string					//实例id
	name string					//实例名
	region string				//地区
	cs string					//实例连接地址串
}
type database struct {			//数据库
	instance_info instance		//实例信息
	name string					//数据库名
	tables []table				//表信息
}

type table struct {				//表
	name string					//表名
	data []interface{}			//检查结果
}

type Config struct {
	USER                       string   `yaml:"USER"`
	PASSWD                     string   `yaml:"PASSWD"`
	RDS_FILEPATH               string   `yaml:"RDS_FILEPATH"`
	TIMEOUT                    int		`yaml:"TIMEOUT"`
	LIMIT_COUNT                int      `yaml:"LIMIT_COUNT"`
	EXCEL_NAME                 string   `yaml:"EXCEL_NAME"`
	LOGPATH                    string   `yaml:"LOGPATH"`
	DB_WHITE_LIST              []string `yaml:"DB_WHITE_LIST"`
	RDS_INFO_FILEPATH          string	`yaml:"RDS_INFO_FILEPATH"`
	REGIONS                    []string `yaml:"REGIONS"`
	ACCESSKEYID                string	`yaml:"ACCESSKEYID"`
	ACCESSKEYSECRET            string	`yaml:"ACCESSKEYSECRET"`
}