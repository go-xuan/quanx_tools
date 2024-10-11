package enums

import (
	"fmt"

	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/enumx"

	"quanx_tools/common"
)

func Print(color fmtx.Color, enum *enumx.StringEnum[string]) {
	for _, k := range enum.Keys() {
		fmt.Printf("%-30s %s\n", color.String(k), enum.Get(k))
	}
}

var (
	RandTypeExplain = enumx.NewStringEnum[string]().
			Add(common.String, "字符串，默认长度为10。可通过args中的length进行约束，例：length=14").
			Add(common.Int, "整数，默认取值范围为1~9999，可通过args中的min、max进行约束，例：min=1&max=100").
			Add(common.Float, "浮点数，可通过args中的min、max、prec(精度)参数进行约束，例：min=1&max=100&prec=6").
			Add(common.Sequence, "序列，起始值取值于args中的min值").
			Add(common.Time, "时间，格式为：YYYY-MM-DD hh:mm:ss，默认取值范围为近30天，可通过args中的min、max进行约束，例：min=2024-01-01").
			Add(common.Date, "日前，格式为：YYYY-MM-DD，默认取值范围为1970-01-01至今天，可通过args中的min、max进行约束，例：min=2024-01-01").
			Add(common.Password, "密码，可通过args中的length、level（难度级别：1/2/3）参数进行约束").
			Add(common.Enum, `枚举取值，取值范围可通过args中的enums（多个以","分隔）进行约束`).
			Add(common.Database, "从数据库取值，需要在args里配置数据库配置信息，仅会取数据库中该字段去重后的前100条数据").
			Add(common.Uuid, "UUID").
			Add(common.Phone, "手机号").
			Add(common.Name, "姓名").
			Add(common.IdCard, "身份证").
			Add(common.PlateNo, "车牌号").
			Add(common.Email, "邮箱").
			Add(common.IP, "IP").
			Add(common.Province, "省").
			Add(common.City, "城市，目前仅支持湖北省内城市")

	RandArgsExplain = enumx.NewStringEnum[string]().
			Add("prefix", "前缀，所有type通用，非必要参数").
			Add("suffix", "后缀，所有type通用，非必要参数").
			Add("upper", "大写，true为开启，所有type通用，非必要参数").
			Add("lower", "小写，true为开启，所有type通用，非必要参数").
			Add("old", "用于替换的旧值，必须和new搭配使用，所有type通用，非必要参数").
			Add("new", "用于替换的新值，必须和old搭配使用，所有type通用，非必要参数").
			Add("min", "取值范围下限，部分type可用，非必要参数").
			Add("min", "取值范围上限，部分type可用，非必要参数").
			Add("format", "格式话，用于time和date类型的格式化").
			Add("length", "字符长度，用于控制string类型的长度").
			Add("prec", "浮点精度，用于指定float精度").
			Add("level", "密码级别，用于指定password难度级别").
			Add("enums", `枚举取值，多个以","分隔`).
			Add("type", `数据库类型，仅type=database时使用`).
			Add("host", `数据库Host，仅type=database时使用`).
			Add("port", `数据库port，仅type=database时使用`).
			Add("username", `数据库用户名，仅type=database时使用`).
			Add("password", `数据库密码，仅type=database时使用`).
			Add("database", `数据库名，仅type=database时使用`).
			Add("field", `字段名，仅type=database时使用`).
			Add("table", `数据库表名，仅type=database时使用`)

	RequiredArgsExamples = enumx.NewStringEnum[string]().
				Add(common.Enum, `-args="enums=1,2,3"`).
				Add(common.Database, `-args="type=postgres&host=localhost&port=5432&username=postgres&password=username&database=demo&table=t_user&field=user_name"`)

	RequiredArgsExplain = enumx.NewStringEnum[*enumx.StringEnum[string]]().
				Add(common.Enum, EnumArgs).
				Add(common.Database, DatabaseArgs)

	EnumArgs = enumx.NewStringEnum[string]().
			Add("enums", `枚举取值范围，多个值以","分隔`)

	DatabaseArgs = enumx.NewStringEnum[string]().
			Add("type", `数据库类型，必填，可选值：mysql/postgres"`).
			Add("host", `数据库Host，必填`).
			Add("port", `数据库port，必填`).
			Add("username", `数据库用户名，必填`).
			Add("password", `数据库密码，必填`).
			Add("database", `数据库名，必填`).
			Add("table", `数据库表名，必填`).
			Add("field", `字段名，必填`)

	CryptoFunc = enumx.NewStringEnum[string]().
			Add("upper", `转大写，例如：-formula=upper(abc)，将字符串“abc”转为大写“ABC”`).
			Add("lower", `转小写，例如：-formula=upper(ABC)，将字符串“ABC”转为小写“abc”`).
			Add("reverse", `反转字符串，例如：-formula=reverse(abc)，将字符串“abc”转为“cba”`).
			Add("md5", `md5加密，例如：-formula=md5(abc)，将字符串“abc”进行md5加密`).
			Add("base64", `base64加密，例如：-formula=base64(abc)，将字符串“abc”进行base64加密`)

	CryptoVariable = enumx.NewStringEnum[string]().
			Add("general", `普通变量，变量名可自定义。例如：-variables="key=123"，将formula中的变量{key}赋值为123`).
			Add("uuid", `特殊变量，变量名可自定义。例如：-variables="key=uuid"，将formula中的变量{key}赋值为随机生成的UUID`).
			Add("timestamp", `特殊变量，变量名可自定义。例如：-variables="key=timestamp"，将formula中的变量{key}赋值为当前unix时间戳"`)

	PasswordEnum = enumx.NewStringEnum[string]().
			Add("postgres", `
source: pgsql
enable: true
type: postgres
host: localhost
port: 5432
username: postgres
password: postgres
database: 
debug: true
`).
		Add("mysql", `
source: mysql
enable: false
type: mysql
host: localhost
port: 3306
username: root
password: mysql123qc
database: 
debug: true
`).
		Add("redis", `
source: mysql
enable: false
host: localhost
port: 6379
username: 
password: Init@123
`)
	LinuxEnum = enumx.NewStringEnum[string]().
			Add("", "")

	MacEnum = enumx.NewStringEnum[string]().
		Add(`ifconfig en0 | grep "inet " | awk '{print $2}'`, "查看本机IP")
)
