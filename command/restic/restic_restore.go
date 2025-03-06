package restic

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
)

type RestoreCommander struct {
	Restore    *Restore    `json:"restore" yaml:"restore"`       // 恢复配置
	Repository *Repository `json:"repository" yaml:"repository"` // 存储库配置
}

func (e *RestoreCommander) Execute() (string, error) {
	// 检查恢复路径
	if checkPath := e.Restore.Path; filex.Exists(checkPath) {
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
		return out, errorx.Wrap(err, "恢复restic快照失败")
	}
	// restic restore命令恢复的文件需要移动至恢复路径
	if snapshotPath := e.Restore.Snapshot.Path; snapshotPath != "" {
		restorePath := e.Restore.Path
		mvCmd := fmt.Sprintf("mv -f %s/* %s/", filepath.Join(restorePath, snapshotPath), restorePath)
		_, _ = ExecCommandAndLog(mvCmd, "移动恢复文件")

	}
	return out, nil
}

// Restore 恢复命令
type Restore struct {
	Snapshot Snapshot `json:"snapshot" yaml:"snapshot"` // 恢复快照
	Path     string   `json:"path" yaml:"path"`         // 恢复路径
}

func (r *Restore) Command() string {
	var cmd = strings.Builder{}
	cmd.WriteString(`restic restore `)
	if r.Snapshot.ShortId != "" {
		cmd.WriteString(r.Snapshot.ShortId)
	} else if r.Snapshot.Id != "" {
		cmd.WriteString(r.Snapshot.Id)
	} else {
		cmd.WriteString(` latest `)
	}
	cmd.WriteString(` --target `)
	cmd.WriteString(filex.Pwd(r.Path))
	cmd.WriteString(` --verbose `)
	if r.Snapshot.Host != "" {
		cmd.WriteString(` --host `)
		cmd.WriteString(r.Snapshot.Host)
	}
	if r.Snapshot.Tag != "" {
		cmd.WriteString(` --tag `)
		cmd.WriteString(r.Snapshot.Tag)
	}
	return cmd.String()
}
