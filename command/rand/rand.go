package rand

import (
	"fmt"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/slicex"
	"github.com/go-xuan/quanx/utils/randx"

	"quanx_tools/command"
	"quanx_tools/common"
	"quanx_tools/common/dao"
	"quanx_tools/common/enums"
	"quanx_tools/common/utils"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Rand, "随机数生成器",
		flagx.StringOption("type", "数据类型", ""),
		flagx.IntOption("size", "生成数量", 1),
		flagx.StringOption("args", "约束参数", ""),
		flagx.StringOption("default", "默认值", ""),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetExecutor(executor)
}

// 随机数生成器
func executor() error {
	if Command.GetHelpOptionValue().Bool() {
		fmtx.Cyan.XPrintf("执行%s命令时，", command.Rand)
		Command.OptionsHelp()
		return nil
	}

	randType := Command.GetOptionValue("type").String()
	if randType == "" && Command.GetHelpOptionValue().Bool() {
		fmtx.Cyan.XPrintf("执行%s命令时，可用的-type参数值列表：\n", command.Rand)
		enums.Print(fmtx.Green, enums.RandTypes)
		return nil
	}

	if randType != "" && !slicex.Contains(enums.RandTypes.Keys(), randType) {
		fmtx.Red.XPrintf("当前输入命令 -type=%s 暂不支持，以下是可用的-type参数值：\n", randType)
		enums.Print(fmtx.Green, enums.RandTypes)
		return nil
	}

	args := Command.GetOptionValue("args").String()
	if args == "" {
		if randType != "" {
			if enum := enums.MustArgsRandTypes.Get(randType); enum != nil {
				fmtx.Magenta.XPrintf("当-type=%s时，-args参数不能为空!\n", randType)
				fmt.Println("-args参数示例：", enums.RandArgsExamples.Get(randType))
				fmt.Println("-args参数说明：")
				enums.Print(fmtx.Green, enum)
				return nil
			}
		} else if Command.GetHelpOptionValue().Bool() {
			fmtx.Red.Println(`args参数可用于约束随机值的生成条件，参数格式为-args="key1=value1&key2=value2"`)
			enums.Print(fmtx.Green, enums.RandArgs)
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
		data := options.RandDataString()
		fmtx.Magenta.Println(data)
		if err := utils.CopyTobePasted(data); err != nil {
			return errorx.Wrap(err, "copy value to be pasted failed")
		}
	} else {
		size := Command.GetOptionValue("size").Int()
		for i := 0; i < size; i++ {
			options.Offset = i
			fmtx.Magenta.Println(options.RandDataString())
		}
	}
	return nil
}
