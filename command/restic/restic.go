package restic

import (
	"fmt"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/execx"
	"github.com/go-xuan/quanx/base/flagx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Restic, "restic工具",
		flagx.StringOption("config", "配置文件", "restic.yaml"),
		flagx.BoolOption("init", "初始化存储库", false),
		flagx.BoolOption("snapshots", "查询快照", false),
		flagx.BoolOption("backup", "备份快照", false),
		flagx.BoolOption("restore", "恢复快照", false),
		flagx.BoolOption("forget", "删除快照", false),
	).SetExecutor(executor)
}

// restic脚本执行
func executor() error {
	// 读取配置文件
	var config = &Config{}
	if configPath := Command.GetOptionValue("config").String(); configPath != "" {
		if err := marshalx.Apply(configPath).Read(configPath, config); err != nil {
			return errorx.Wrap(err, "读取配置文件失败:"+configPath)
		}
	} else {
		Command.OptionsHelp()
		return nil
	}
	var commander Commander
	if Command.GetOptionValue("backup").Bool() {
		commander = config.BackupCommander()
	} else if Command.GetOptionValue("restore").Bool() {
		commander = config.RestoreCommander()
	} else if Command.GetOptionValue("forget").Bool() {
		commander = config.ForgetCommander()
	} else if Command.GetOptionValue("init").Bool() {
		commander = config.InitRepoCommander()
	} else {
		commander = config.SnapshotsCommander()
	}
	// 直接执行
	if out, err := commander.Execute(); err != nil {
		return errorx.Wrap(err, "执行restic命令失败")
	} else {
		fmt.Println(out)
		return nil
	}
}

type Commander interface {
	Execute() (string, error)
}

// ExecCommandAndLog 执行命令并记录日志
func ExecCommandAndLog(cmd string, msg string) (string, error) {
	fmt.Println(msg, ">>>", cmd)
	if stdout, stderr, err := execx.Command(cmd).Run(); err != nil {
		return stdout + "\n" + stderr, errorx.New(stderr)
	} else {
		return stdout, nil
	}
}
