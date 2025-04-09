package rand

import (
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
	Command        = flagx.NewCommand(command.Rand, "随机数生成器")
	types          = enumx.NewStringEnum[string]()
	params         = enumx.NewStringEnum[string]()
	needParamTypes = enumx.NewStringEnum[*enumx.Enum[string, string]]()
)

func init() {
	types.Add(common.String, "随机数类型-字符串（字符长度为可通过-length参数进行约束，例：length=16）").
		Add(common.Int, "随机数类型-整数（默认取值范围为1~9999，可通过-min、-max参数进行约束，例：min=1&max=100）").
		Add(common.Float, "随机数类型-浮点数（可通过-min、-max、-prec参数数进行约束，例：min=1&max=100&prec=6）").
		Add(common.Sequence, "随机数类型-序列（可通过-min参数约束起始值）").
		Add(common.Time, "随机数类型-时间（格式为YYYY-MM-DD hh:mm:ss，取值范围可通过-min、-max参数进行约束，例：min=1970-01-01）").
		Add(common.Date, "随机数类型-日期（格式为YYYY-MM-DD，取值范围可通过-min、-max参数进行约束，例：min=1970-01-01）").
		Add(common.Password, "随机数类型-密码（可通过-length、-level参数进行约束）").
		Add(common.Enum, `随机数类型-枚举取值（取值范围可通过-enums参数进行约束）`).
		Add(common.Database, "随机数类型-数据库（需要通过约束参数配置数据库连接信息，且仅会取数据库中目标字段去重后的前100条数据）").
		Add(common.Uuid, "随机数类型-UUID").
		Add(common.Phone, "随机数类型-手机号").
		Add(common.Name, "随机数类型-姓名").
		Add(common.IdCard, "随机数类型-身份证").
		Add(common.PlateNo, "随机数类型-车牌号").
		Add(common.Email, "随机数类型-邮箱").
		Add(common.IP, "随机数类型-IP").
		Add(common.Province, "随机数类型-省名（包括全国所有省级）").
		Add(common.City, "随机数类型-城市（目前仅支持湖北省地级市）")

	params.Add("prefix", "约束参数-前缀（所有类型通用，非必要参数）").
		Add("suffix", "约束参数-后缀（所有类型通用，非必要参数）").
		Add("upper", "约束参数-大写（true/false，所有类型通用，非必要参数）").
		Add("lower", "约束参数-小写（true/false，所有类型通用，非必要参数）").
		Add("old", "约束参数-replace旧值（必须和-new搭配使用，所有类型通用，非必要参数）").
		Add("new", "约束参数-replace新值（必须和-old搭配使用，所有类型通用，非必要参数）").
		Add("min", "约束参数-取值范围下限（部分类型可用，非必要参数）").
		Add("max", "约束参数-取值范围上限（部分类型可用，非必要参数）").
		Add("format", "约束参数-格式化（用于time和date类型的格式化）").
		Add("length", "约束参数-字符长度（用于约束string类型的长度）").
		Add("prec", "约束参数-浮点精度（用于指定float类型的精度）").
		Add("level", "约束参数-密码级别（用于指定password难度级别，1/2/3）").
		Add("enums", `约束参数-枚举取值（多个以","分隔）`).
		Add("db_type", `约束参数-数据库类型（仅-database时使用）`).
		Add("db_host", `约束参数-数据库Host（仅-database时使用）`).
		Add("db_port", `约束参数-数据库port（仅-database时使用）`).
		Add("db_user", `约束参数-数据库用户名（仅-database时使用）`).
		Add("db_pwd", `约束参数-数据库密码（仅-database时使用）`).
		Add("db_name", `约束参数-数据库名（仅-database时使用）`).
		Add("db_table", `约束参数-数据库表名（仅-database时使用）`).
		Add("db_field", `约束参数-数据库表字段（仅-database时使用）`)

	needParamTypes.
		Add(common.Enum, enumx.NewStringEnum[string]().
			Add("enums", `枚举取值，多个值以","分隔`)).
		Add(common.Password, enumx.NewStringEnum[string]().
			Add("length", `密码长度`).
			Add("level", "密码难度级别，1-仅数字/2-数字+字母/3-数字+字母+特殊符号")).
		Add(common.Database, enumx.NewStringEnum[string]().
			Add("db_type", `数据库类型，必填，可选值：mysql/postgres`).
			Add("db_host", `数据库Host，必填`).
			Add("db_port", `数据库port，必填`).
			Add("db_user", `数据库用户名，必填`).
			Add("db_pwd", `数据库密码，必填`).
			Add("db_name", `数据库名，必填`).
			Add("db_table", `数据库表名，必填`).
			Add("db_field", `数据库表字段，必填`))

	Command.AddOption(
		flagx.IntOption("size", "生成数量", 1),
		flagx.StringOption("default", "默认值", ""),
		flagx.BoolOption("copy", "复制结果值", false),
	)

	// 添加type选项
	for _, key := range types.Keys() {
		Command.AddOption(flagx.BoolOption(key, types.Get(key), false))
	}

	// 添加param选项
	for _, key := range params.Keys() {
		Command.AddOption(flagx.StringOption(key, params.Get(key), ""))
	}

	Command.SetExecutor(executor)
}

// 随机数生成器
func executor() error {
	// 获取随机数类型
	var randType string
	for _, key := range types.Keys() {
		if Command.GetOptionValue(key).Bool() {
			randType = key
			break
		}
	}

	if randType == "" && Command.NeedHelp() {
		fmtx.Cyan.Println("可支持的随机数类型：")
		enums.Print(fmtx.Green, types)
		return nil
	} else if randType != "" && !types.Exist(randType) {
		fmtx.Red.Xprintf("当前随机数类型-%s暂不支持，以下是可支持选项：\n", randType)
		enums.Print(fmtx.Green, types)
		return nil
	}

	// 收集param参数
	var param = make(map[string]string)
	for _, key := range params.Keys() {
		if value := Command.GetOptionValue(key).String(); value != "" {
			param[key] = value
		}
	}

	if len(param) == 0 {
		if randType != "" {
			if needParamHelp := needParamTypes.Get(randType); needParamHelp != nil {
				fmtx.Magenta.Xprintf("当-%s时，需要带上约束参数！相关约束参数说明：\n", randType)
				enums.Print(fmtx.Green, needParamHelp)
				return nil
			}
		} else if Command.NeedHelp() {
			fmtx.Red.Println(`参数可用于约束随机值的生成条件，以下是可支持参数"`)
			enums.Print(fmtx.Green, params)
			return nil
		}
	}

	options := randx.Options{
		Type:    randType,
		Default: Command.GetOptionValue("default").String(),
		Param:   randx.NewParam(param),
	}
	if options.Type == common.Database {
		if data, err := dao.GetDBFieldDataList(param); err != nil {
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
		fmtx.Magenta.Xprintf("当前值%s已复制到粘贴板\n", data)
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
