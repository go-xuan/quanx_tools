package restic

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/fmtx"
)

type BackupCommander struct {
	Backup     *Backup     `json:"backup" yaml:"backup"`         // 备份配置
	Repository *Repository `json:"repository" yaml:"repository"` // 存储库配置
	Datasource *Datasource `json:"datasource" yaml:"datasource"` // 备份数据源
}

func (e *BackupCommander) Execute() (string, error) {
	var dumpCmd, dumpPath string
	if dumpCmd, dumpPath = e.Datasource.Dump(); dumpCmd != "" {
		filex.CreateDir(dumpPath)
		if out, err := ExecCommandAndLog(dumpCmd, "根据数据库配置执行dump"); err != nil {
			return out, errorx.Wrap(err, "执行数据库dump失败")
		}
	}

	backupPath := filepath.Join(e.Backup.Path, dumpPath)
	// 检查备份路径
	if !filex.Exists(backupPath) {
		filex.CreateDir(backupPath)
	} else if !filex.IsDir(backupPath) {
		_ = os.Remove(backupPath)
	}
	// 将dump文件移动至备份路径
	mvCmd := fmt.Sprintf("mv %s/* %s", dumpPath, backupPath)
	if out, err := ExecCommandAndLog(mvCmd, "移动dump文件至备份路径"); err != nil {
		return out, errorx.Wrap(err, "移动dump文件至备份路径失败")
	}

	// 执行restic命令
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(e.Backup.Command())
	if out, err := ExecCommandAndLog(cmd.String(), "创建restic快照"); err != nil {
		return out, errorx.Wrap(err, "创建restic快照失败")
	} else {
		if match := regexp.MustCompile(`snapshot (\w+) saved`).FindStringSubmatch(out); len(match) > 1 {
			fmt.Println("生成快照ID >>> ", fmtx.Red.String(match[1]))
		}
		return out, nil
	}
}

// Backup 备份参数
type Backup struct {
	Host string `json:"host" yaml:"host"` // 主机名称
	Tag  string `json:"tag"  yaml:"tag" ` // 标签
	Path string `json:"path" yaml:"path"` // 备份路径
}

func (b *Backup) Command() string {
	var cmd = strings.Builder{}
	cmd.WriteString(`restic backup `)
	cmd.WriteString(filex.Pwd(b.Path))
	cmd.WriteString(` --verbose `)
	if b.Host != "" {
		cmd.WriteString(` --host `)
		cmd.WriteString(b.Host)
	}
	if b.Tag != "" {
		cmd.WriteString(` --tag `)
		cmd.WriteString(b.Tag)
	}
	return cmd.String()
}
