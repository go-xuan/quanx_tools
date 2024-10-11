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
		flagx.StringOption("type", "数据类型", "string"),
		flagx.IntOption("size", "生成数量", 1),
		flagx.StringOption("args", "约束参数", ""),
		flagx.StringOption("default", "默认值", ""),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetHandler(handler)
}

// 随机数生成器
func handler() error {
	randType := Command.GetOptionValue("type").String()
	args := Command.GetOptionValue("args").String()
	if randType == "" && args == "" {
		fmtx.Cyan.XPrintf("执行%s命令时，参数不能为空！", command.Rand)
		Command.Help()
		return nil
	}

	if !slicex.Contains(enums.RandTypeExplain.Keys(), randType) {
		fmtx.Red.XPrintf("当前输入的 %s 不可用！", "-type="+randType)
		fmtx.Magenta.XPrintf("以下是可用的%s参数值：\n", "-type")
		enums.Print(fmtx.Green, enums.RandTypeExplain)
		return nil
	}

	if args == "" {
		if enum := enums.RequiredArgsExplain.Get(randType); enum != nil {
			fmtx.Magenta.XPrintf("当-type=%s时，args参数不能为空", randType)
			fmt.Println("\nargs参数示例：")
			fmt.Println(enums.RequiredArgsExamples.Get(randType))

			fmt.Println("\nargs参数说明：")
			enums.Print(fmtx.Green, enum)
			return nil
		}
	} else if args == "explain" {
		Command.Help()
		fmtx.Red.Println(`args参数可用于约束随机值的生成条件，参数格式为-args="key=value&key=value"`)
		enums.Print(fmtx.Green, enums.RandArgsExplain)
		return nil
	}

	options := randx.Options{
		Type:    randType,
		Default: Command.GetOptionValue("default").String(),
		Param:   randx.NewParam(args),
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
