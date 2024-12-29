package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/execx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/fmtx"
	log "github.com/sirupsen/logrus"
)

type ResticExecutor interface {
	Execute() (string, error)
}

type BackupExecutor struct {
	Backup     *Backup     `json:"backup"        yaml:"backup"        comment:"备份"`
	Repository *Repository `json:"repository"    yaml:"repository"    comment:"存储库信息"`
	Datasource *Datasource `json:"datasource"    yaml:"datasource"    comment:"备份数据源"`
}

type RestoreExecutor struct {
	Restore    *Restore    `json:"restore"       yaml:"restore"       comment:"恢复"`
	Repository *Repository `json:"repository"    yaml:"repository"    comment:"存储库信息"`
}

type ForgetExecutor struct {
	Forget     *Forget     `json:"forget"        yaml:"forget"        comment:"删除"`
	Repository *Repository `json:"repository"    yaml:"repository"    comment:"存储库信息"`
}

type InitRepoExecutor struct {
	Repository *Repository `json:"repository"    yaml:"repository"    comment:"存储库信息"`
}

type SnapshotsExecutor struct {
	Repository *Repository `json:"repository"    yaml:"repository"    comment:"存储库信息"`
}

func (e *BackupExecutor) Execute() (string, error) {
	var dumpCmd, dumpPath string
	if dumpCmd, dumpPath = e.Datasource.Dump(); dumpCmd != "" {
		dir, _ := filepath.Split(dumpPath)
		filex.CreateDir(dir)
		if _, err := ExecCommandAndLog(dumpCmd, "根据数据库配置执行dump"); err != nil {
			return "", errorx.Wrap(err, "执行数据库dump失败")
		}
	}
	// 将dump文件移动至备份路径
	if backupPath := e.Backup.Path; filex.IsDir(backupPath) {
		backupPath = backupPath + string(os.PathSeparator)
	} else {
		mvCmd := fmt.Sprintf("mv %s %s", dumpPath, backupPath)
		if _, err := ExecCommandAndLog(mvCmd, "移动dump文件至备份路径"); err != nil {
			return "", errorx.Wrap(err, "移动dump文件至备份路径失败")
		}
		// 检查备份路径是否存在或者为空
		if checkPath := e.Backup.Path; !filex.Exists(checkPath) {
			return "", errorx.Errorf("此备份文件或者文件夹不存在: %s", checkPath)
		} else if filex.IsDir(checkPath) && filex.IsEmptyDir(checkPath) {
			return "", errorx.Errorf("此备份文件夹为空: %s", checkPath)
		}
	}

	// 执行restic命令
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(e.Backup.Command())
	out, err := ExecCommandAndLog(cmd.String(), "创建restic快照")
	if err != nil {
		return "", errorx.Wrap(err, "创建restic快照失败")
	}
	match := regexp.MustCompile(`snapshot\s+(\w+)\s+saved`).FindStringSubmatch(out)
	fmtx.Red.XPrintf("生成快照ID：%s", match[0])
	return out, nil
}

func (e *RestoreExecutor) Execute() (string, error) {
	// 检查恢复路径是否存在
	if checkPath := e.Restore.Target; filex.Exists(checkPath) {
		// 检查恢复路径是否是文件夹
		if !filex.IsDir(checkPath) {
			return "", errorx.Errorf("此恢复路径不是文件夹: %s", checkPath)
		}
	} else {
		filex.CreateDir(checkPath)
	}

	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(e.Restore.Command())
	out, err := ExecCommandAndLog(cmd.String(), "恢复restic快照")
	if err != nil {
		return "", errorx.Wrap(err, "恢复restic快照失败")
	}

	// restic restore命令生成的恢复文件需要移动至所需恢复路径
	if backupPath := e.Restore.Path; backupPath != "" {
		restorePath := e.Restore.Target
		backupPath = filepath.Join(restorePath, filex.Pwd(backupPath))
		mvCmd := fmt.Sprintf("mv %s %s/", backupPath, restorePath)
		_, _ = ExecCommandAndLog(mvCmd, "移动恢复文件")
	}
	return out, nil
}

func (e *ForgetExecutor) Execute() (string, error) {
	// 执行restic命令
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(e.Forget.Command())
	if out, err := ExecCommandAndLog(cmd.String(), "删除restic快照"); err != nil {
		return "", errorx.Wrap(err, "删除restic快照失败")
	} else {
		return out, nil
	}
}

func (e *InitRepoExecutor) Execute() (string, error) {
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(`restic init --verbose`)
	if out, err := ExecCommandAndLog(cmd.String(), "初始化restic存储库"); err != nil {
		return "", errorx.Wrap(err, "初始化restic存储库失败")
	} else {
		return out, nil
	}
}

func (e *SnapshotsExecutor) Execute() (string, error) {
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(`restic snapshots --json`)

	if out, err := ExecCommandAndLog(cmd.String(), "执行restic快照查询命令"); err != nil {
		return "", errorx.Wrap(err, "查询restic快照失败")
	} else {
		return out, nil
	}
}

// ExecCommandAndLog 执行命令并记录日志
func ExecCommandAndLog(cmd string, msg string) (string, error) {
	log.WithField(`cmd`, cmd).Info(msg)
	if stdout, stderr, err := execx.Command(cmd).Run(); err != nil {
		log.WithField(`cmd`, cmd).Error(err)
		return stdout + "\n" + stderr, err
	} else {
		return stdout, nil
	}
}
