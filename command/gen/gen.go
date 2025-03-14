package gen

import (
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
	"quanx_tools/command/gen/internal"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Gen, "代码生成工具",
		flagx.StringOption("config", "配置文件", "gen.yaml"),
		flagx.BoolOption("check", "检测外置模板", false),
	).SetExecutor(executor)
}

func executor() error {
	// 读取配置文件
	var config = &internal.Config{}
	var configPath = Command.GetOptionValue("config").String()
	if err := marshalx.Apply(configPath).Read(configPath, config); err != nil {
		return errorx.Wrap(err, "配置文件读取失败:"+configPath)
	}
	// 检测外置模板
	if Command.GetOptionValue("check").Bool() {
		if err := config.ExternalTemplateCheck(); err != nil {
			return errorx.Wrap(err, "外置模板检测失败")
		}
		return nil
	}
	// 代码生成
	if err := config.Generator().Execute(); err != nil {
		return errorx.Wrap(err, "代码生成失败")
	}
	return nil
}
