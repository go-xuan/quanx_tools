package restic

import (
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/flagx"
	"github.com/go-xuan/quanx/utils/marshalx"

	"quanx_tools/command"
	"quanx_tools/command/restic/internal"
)

var Command *flagx.Command

func init() {
	Command = flagx.NewCommand(command.Restic, "restic工具",
		flagx.StringOption("config", "配置文件", "restic.yaml"),
		flagx.BoolOption("backup", "备份快照", false),
		flagx.BoolOption("restore", "恢复快照", false),
		flagx.BoolOption("forget", "删除快照", false),
		flagx.BoolOption("init", "初始化存储库", false),
	).SetHandler(handler)
}

// restic脚本执行
func handler() error {
	// 读取配置文件
	var config = &internal.Config{}
	var configPath = Command.GetOptionValue("config").String()
	if err := marshalx.UnmarshalFromFile(configPath, config); err != nil {
		return errorx.Wrap(err, "读取配置文件失败:"+configPath)
	}

	var executor internal.ResticExecutor
	if Command.GetOptionValue("backup").Bool() {
		executor = &internal.BackupExecutor{
			Backup:     config.Backup,
			Repository: config.Repository,
			Datasource: config.Datasource,
		}
	} else if Command.GetOptionValue("restore").Bool() {
		executor = &internal.RestoreExecutor{
			Restore:    config.Restore,
			Repository: config.Repository,
		}
	} else if Command.GetOptionValue("forget").Bool() {
		executor = &internal.ForgetExecutor{
			Forget:     config.Forget,
			Repository: config.Repository,
		}
	} else if Command.GetOptionValue("init").Bool() {
		executor = &internal.InitRepoExecutor{
			Repository: config.Repository,
		}
	} else {
		executor = &internal.SnapshotsExecutor{
			Repository: config.Repository,
		}
	}

	// 直接执行
	if _, err := executor.Execute(); err != nil {
		return errorx.Wrap(err, "执行restic命令失败")
	}
	return nil
}
