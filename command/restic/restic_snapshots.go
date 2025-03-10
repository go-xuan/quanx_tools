package restic

import (
	"github.com/go-xuan/quanx/utils/marshalx"
	"strings"
	"time"

	"github.com/go-xuan/quanx/base/errorx"
)

type Snapshot struct {
	Id      string `json:"id"`       // 快照ID
	ShortId string `json:"short_id"` // 快照短ID
	Host    string `json:"host"`     // 主机名
	Tag     string `json:"tag"`      // 标签
	Path    string `json:"path"  `   // 路径
}

type SnapshotsCommander struct {
	Repository *Repository `json:"repository" yaml:"repository"` //存储库信息
}

func (e *SnapshotsCommander) Execute() (string, error) {
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(`restic snapshots --json`)
	if out, err := ExecCommandAndLog(cmd.String(), "执行restic快照查询命令"); err != nil {
		return out, errorx.Wrap(err, "查询restic快照失败")
	} else {
		var snapshots []*SnapshotOut
		m := marshalx.Json("    ")
		_ = m.Unmarshal([]byte(out), &snapshots)
		bytes, _ := m.Marshal(&snapshots)
		out = string(bytes)
		return out, nil
	}
}

type SnapshotOut struct {
	Id             string    `json:"id"`                // 快照ID
	ShortId        string    `json:"short_id"`          // 快照短ID
	Time           time.Time `json:"time"`              // 快照时间
	Hostname       string    `json:"hostname"`          // 主机名
	Tags           []string  `json:"tags"`              // 标签
	Paths          []string  `json:"paths"  `           // 路径
	Tree           string    `json:"tree"`              // 树
	Username       string    `json:"username"`          // 用户名
	ProgramVersion string    `json:"program_version"  ` // restic版本
}
