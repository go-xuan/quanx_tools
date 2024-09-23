package rand

import (
	"fmt"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/types/slicex"
	"github.com/go-xuan/quanx/utils/fmtx"
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
		flagx.StringOption("args", "约束参数", ""),
		flagx.IntOption("size", "生成数量", 1),
		flagx.StringOption("default", "默认值", ""),
		flagx.BoolOption("copy", "复制粘贴", false),
	).SetHandler(handler)
}

// 随机数生成器
func handler() error {
	randType := Command.GetOptionValue("type").String()
	args := Command.GetOptionValue("args").String()
	if randType == "" && args == "" {
		Command.Help()
		return nil
	}
	if randType == "explain" || !slicex.Contains(enums.RandTypeExplain.Keys(), randType) {
		fmt.Println(`type参数可用于约束随机值的类型，以下是可支持的type类型：`)
		for _, k := range enums.RandTypeExplain.Keys() {
			fmt.Printf("%-30s %s\n", fmtx.Green.String(k), enums.RandTypeExplain.Get(k))
		}
		return nil
	}

	if args == "" {
		if enum := enums.RequiredArgsExplain.Get(randType); enum != nil {
			fmtx.Magenta.XPrintf("当-type=%s时，args参数不能为空", randType)
			fmt.Println("\nargs参数示例：")
			fmt.Println(enums.RequiredArgsExamples.Get(randType))
			fmt.Println("\nargs参数说明：")
			for _, key := range enum.Keys() {
				fmt.Printf("%s: %s\n", fmtx.Green.String(key), enum.Get(key))
			}
			return nil
		}
	} else if args == "explain" {
		Command.Help()
		fmtx.Red.Println(`args参数可用于约束随机值的生成条件，参数格式为-args="key=value&key=value"`)
		for _, k := range enums.RandArgsExplain.Keys() {
			fmt.Printf("%-30s %s\n", fmtx.Green.String(k), enums.RandArgsExplain.Get(k))
		}
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
