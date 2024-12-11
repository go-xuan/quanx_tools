package gen

import (
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
	"quanx_tools/command/gen/internal"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Gen, "代码生成工具",
		flagx.StringOption("config", "配置文件", "gen.yaml"),
		flagx.BoolOption("check", "模板检测", false),
	).SetExecutor(executor)
}

func executor() error {
	// 读取配置文件
	var config = &internal.Config{}
	var configPath = Command.GetOptionValue("config").String()
	if err := marshalx.UnmarshalFromFile(configPath, config); err != nil {
		return errorx.Wrap(err, "读取配置文件失败:"+configPath)
	}

	if check := Command.GetOptionValue("check").Bool(); check {
		if err := config.CheckTemplate(); err != nil {
			return errorx.Wrap(err, "检测模板文件失败")
		}
		return nil
	}
	// 代码生成
	if err := config.Generator().Execute(); err != nil {
		return errorx.Wrap(err, "代码生成失败")
	}
	return nil
}
