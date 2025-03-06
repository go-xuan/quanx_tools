package restic

import (
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
)

type InitRepoCommander struct {
	Repository *Repository `json:"repository" yaml:"repository"` // 存储库配置
}

func (e *InitRepoCommander) Execute() (string, error) {
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(`restic init --verbose`)
	if out, err := ExecCommandAndLog(cmd.String(), "初始化restic存储库"); err != nil {
		return out, errorx.Wrap(err, "初始化restic存储库失败")
	} else {
		return out, nil
	}
}
