package restic

import (
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
)

type ForgetCommander struct {
	Forget     *Forget     `json:"forget" yaml:"forget"`         // 删除配置
	Repository *Repository `json:"repository" yaml:"repository"` // 存储库配置
}

func (e *ForgetCommander) Execute() (string, error) {
	// 执行restic命令
	var cmd = strings.Builder{}
	cmd.WriteString(e.Repository.SetEnv())
	cmd.WriteString(e.Forget.Command())
	if out, err := ExecCommandAndLog(cmd.String(), "删除restic快照"); err != nil {
		return out, errorx.Wrap(err, "删除restic快照失败")
	} else {
		return out, nil
	}
}

// Forget 删除命令
type Forget struct {
	Snapshot Snapshot  `json:"snapshot" yaml:"snapshot"` // 删除快照
	Strategy *Strategy `json:"strategy" yaml:"strategy"` // 删除策略
}

func (r *Forget) Command() string {
	var cmd = strings.Builder{}
	cmd.WriteString(`restic forget `)
	if r.Snapshot.ShortId != "" {
		cmd.WriteString(r.Snapshot.ShortId)
	} else if r.Snapshot.Id != "" {
		cmd.WriteString(r.Snapshot.Id)
	} else if r.Strategy != nil {
		cmd.WriteString(r.Strategy.Options())
	} else {
		cmd.WriteString(` --keep-last 3 `)
	}
	cmd.WriteString(` --verbose --prune `)
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

// Strategy 删除策略
// --keep-last n   keep the last n snapshots (use 'unlimited' to keep all snapshots)
// --keep-hourly n   keep the last n hourly snapshots (use 'unlimited' to keep all hourly snapshots)
// --keep-daily n   keep the last n daily snapshots (use 'unlimited' to keep all daily snapshots)
// --keep-weekly n   keep the last n weekly snapshots (use 'unlimited' to keep all weekly snapshots)
// --keep-monthly n   keep the last n monthly snapshots (use 'unlimited' to keep all monthly snapshots)
// --keep-yearly n   keep the last n yearly snapshots (use 'unlimited' to keep all yearly snapshots)
// --keep-within duration  keep snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-hourly duration keep hourly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-daily duration keep daily snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-weekly duration keep weekly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-monthly duration keep monthly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-yearly duration keep yearly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-tag tag-list  keep snapshots with this tag0list (can be specified multiple times) (default [])
type Strategy struct {
	KeepLast          string `json:"keep-last" yaml:"keep-last"`                     // 保留最近n个快照
	KeepHourly        string `json:"keep-hourly" yaml:"keep-hourly"`                 // 保留近n小时内的快照（使用unlimited保留所有小时快照）
	KeepDaily         string `json:"keep-daily" yaml:"keep-daily"`                   // 保留近n天内的快照（使用unlimited保留所有每日快照）
	KeepWeekly        string `json:"keep-weekly" yaml:"keep-weekly"`                 // 保留近n周内的快照（使用unlimited保留所有每周快照）
	KeepMonthly       string `json:"keep-monthly" yaml:"keep-monthly"`               // 保留近n月内的快照（使用unlimited保留所有每周快照）
	KeepYearly        string `json:"keep-yearly" yaml:"keep-yearly"`                 // 保留近n年内的快照（使用unlimited保留所有每周快照）
	KeepTag           string `json:"keep-tag" yaml:"keep-tag"`                       // 根据此标签列表保留快照（可以指定多个，默认[]）
	KeepWithin        string `json:"keep-within" yaml:"keep-within"`                 // 相对于最新快照，保留时间范围内（例如1y5m7d2h）的快照
	KeepWithinHourly  string `json:"keep-within-hourly" yaml:"keep-within-hourly"`   // 相对于最新快照，保留时间范围内（例如1y5m7d2h）的每小时快照
	KeepWithinDaily   string `json:"keep-within-daily" yaml:"keep-within-daily"`     // 相对于最新快照，保留时间范围内（例如1y5m7d2h）的每日快照
	KeepWithinWeekly  string `json:"keep-within-weekly" yaml:"keep-within-weekly"`   // 相对于最新快照，保留时间范围内（例如1y5m7d2h）的每周快照
	KeepWithinMonthly string `json:"keep-within-monthly" yaml:"keep-within-monthly"` // 相对于最新快照，保留时间范围内（例如1y5m7d2h）的每月快照
	KeepWithinYearly  string `json:"keep-within-yearly" yaml:"keep-within-yearly"`   // 相对于最新快照，保留时间范围内（例如1y5m7d2h）的每年快照
}

func (fs *Strategy) Options() string {
	var sb = strings.Builder{}
	if fs.KeepLast != "" {
		sb.WriteString(` --keep-last `)
		sb.WriteString(fs.KeepLast)
	}
	if fs.KeepHourly != "" {
		sb.WriteString(` --keep-hourly `)
		sb.WriteString(fs.KeepHourly)
	}
	if fs.KeepDaily != "" {
		sb.WriteString(` --keep-daily `)
		sb.WriteString(fs.KeepDaily)
	}
	if fs.KeepWeekly != "" {
		sb.WriteString(` --keep-weekly `)
		sb.WriteString(fs.KeepWeekly)
	}
	if fs.KeepMonthly != "" {
		sb.WriteString(` --keep-monthly `)
		sb.WriteString(fs.KeepMonthly)
	}
	if fs.KeepYearly != "" {
		sb.WriteString(` --keep-yearly `)
		sb.WriteString(fs.KeepYearly)
	}
	if fs.KeepTag != "" {
		sb.WriteString(` --keep-tag `)
		sb.WriteString(fs.KeepTag)
	}
	if fs.KeepWithin != "" {
		sb.WriteString(` --keep-within `)
		sb.WriteString(fs.KeepWithin)
	}
	if fs.KeepWithinHourly != "" {
		sb.WriteString(` --keep-within-hourly `)
		sb.WriteString(fs.KeepWithinHourly)
	}
	if fs.KeepWithinDaily != "" {
		sb.WriteString(` --keep-within-daily `)
		sb.WriteString(fs.KeepWithinDaily)
	}
	if fs.KeepWithinWeekly != "" {
		sb.WriteString(` --keep-within-weekly `)
		sb.WriteString(fs.KeepWithinWeekly)
	}
	if fs.KeepWithinMonthly != "" {
		sb.WriteString(` --keep-within-monthly `)
		sb.WriteString(fs.KeepWithinMonthly)
	}
	if fs.KeepWithinYearly != "" {
		sb.WriteString(` --keep-within-yearly `)
		sb.WriteString(fs.KeepWithinYearly)
	}
	return sb.String()
}
