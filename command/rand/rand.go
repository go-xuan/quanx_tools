package rand

import (
	"fmt"
	"github.com/go-xuan/quanx/types/stringx"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/utils/randx"

	"quanx_tools/command"
	"quanx_tools/common"
	"quanx_tools/common/dao"
	"quanx_tools/common/enums"
	"quanx_tools/common/utils"
)

var (
	Command           = flagx.NewCommand(command.Rand, "随机数生成器")
	allRandTypes      = enumx.NewStringEnum[string]()
	allRandArgs       = enumx.NewStringEnum[string]()
	argsExamples      = enumx.NewStringEnum[string]()
	MustArgsRandTypes = enumx.NewStringEnum[*enumx.StringEnum[string]]()
)

func init() {
	Command.AddOption(
		flagx.StringOption("type", "数据类型", ""),
		flagx.IntOption("size", "生成数量", 1),
		flagx.StringOption("args", "约束参数", ""),
		flagx.StringOption("default", "默认值", ""),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetExecutor(executor)

	allRandTypes.
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
	allRandArgs.
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

	argsExamples.
		Add(common.Enum, `-args="enums=1,2,3"`).
		Add(common.Database, `-args="type=postgres&host=localhost&port=5432&username=postgres&password=username&database=demo&table=t_user&field=user_name"`)

	MustArgsRandTypes.
		Add(common.Enum,
			enumx.NewStringEnum[string]().
				Add("enums", `枚举取值范围，多个值以","分隔`)).
		Add(common.Database,
			enumx.NewStringEnum[string]().
				Add("type", `数据库类型，必填，可选值：mysql/postgres"`).
				Add("host", `数据库Host，必填`).
				Add("port", `数据库port，必填`).
				Add("username", `数据库用户名，必填`).
				Add("password", `数据库密码，必填`).
				Add("database", `数据库名，必填`).
				Add("table", `数据库表名，必填`).
				Add("field", `字段名，必填`))
}

// 随机数生成器
func executor() error {
	randType := Command.GetOptionValue("type").String()
	if randType == "" && Command.NeedHelp() {
		fmtx.Cyan.XPrintf("执行%s命令时，可用的-type参数值列表：\n", command.Rand)
		enums.Print(fmtx.Green, allRandTypes)
		return nil
	}

	if randType != "" && !allRandTypes.Exist(randType) {
		fmtx.Red.XPrintf("当前输入命令 -type=%s 暂不支持，以下是可用的-type参数值：\n", randType)
		enums.Print(fmtx.Green, allRandTypes)
		return nil
	}

	argsStr := Command.GetOptionValue("args").String()
	args := stringx.ParseUrlParams(argsStr)
	if len(args) > 0 {
		if randType != "" {
			if enum := MustArgsRandTypes.Get(randType); enum != nil {
				fmtx.Magenta.XPrintf("当-type=%s时，-args参数不能为空!\n", randType)
				fmt.Println("-args参数示例：", argsExamples.Get(randType))
				fmt.Println("-args参数说明：")
				enums.Print(fmtx.Green, enum)
				return nil
			}
		} else if Command.NeedHelp() {
			fmtx.Red.Println(`args参数可用于约束随机值的生成条件，参数格式为-args="key1=value1&key2=value2"`)
			enums.Print(fmtx.Green, allRandArgs)
			return nil
		}
	}

	options := randx.Options{
		Type:    randType,
		Default: Command.GetOptionValue("default").String(),
		Args:    randx.NewArgs(args),
	}
	if options.Type == common.Database {
		if data, err := dao.GetDBFieldDataList(args); err != nil {
			return errorx.Wrap(err, "failed to query database table field data")
		} else {
			options.Enums = data
		}
	}
	if Command.GetOptionValue("copy").Bool() {
		data := options.NewString()
		fmtx.Magenta.Println(data)
		if err := utils.WriteToClipboard(data); err != nil {
			return errorx.Wrap(err, "copy value to be pasted failed")
		}
		fmtx.Magenta.XPrintf("当前值%s已复制到粘贴板\n", data)
	} else {
		if size := Command.GetOptionValue("size").Int(); size > 0 {
			for i := 0; i < size; i++ {
				options.Offset = i
				fmtx.Magenta.Println(options.NewString())
			}
		}
	}
	return nil
}
