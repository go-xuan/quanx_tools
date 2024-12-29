package internal

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-xuan/quanx/core/gormx"
	"github.com/go-xuan/quanx/os/filex"
)

type Config struct {
	Backup     *Backup     `json:"backup"        yaml:"backup"        comment:"备份"`
	Restore    *Restore    `json:"restore"       yaml:"restore"       comment:"恢复"`
	Forget     *Forget     `json:"forget"        yaml:"forget"        comment:"删除"`
	Repository *Repository `json:"repository"    yaml:"repository"    comment:"存储库信息"`
	Datasource *Datasource `json:"datasource"    yaml:"datasource"    comment:"备份数据源"`
	//RetryTimes    int               `json:"retryTimes"    yaml:"retryTimes"    comment:"重试次数"`
	//RetryInterval int               `json:"retryInterval" yaml:"retryInterval" comment:"重试间隔(秒)"`
	//Cron          string            `json:"cron"          yaml:"cron"          comment:"定时任务"`
}

// Backup 备份命令
type Backup struct {
	Host string `json:"host" yaml:"host" comment:"主机名称"`
	Tag  string `json:"tag"  yaml:"tag"  comment:"标签"`
	Path string `json:"path" yaml:"path" comment:"备份路径"`
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

// Restore 恢复命令
type Restore struct {
	Host     string `json:"host"     yaml:"host"     comment:"主机名称"`
	Tag      string `json:"tag"      yaml:"tag"      comment:"标签"`
	Path     string `json:"path"     yaml:"path"     comment:"备份路径"`
	Target   string `json:"target"   yaml:"target"   comment:"恢复路径"`
	Snapshot string `json:"snapshot" yaml:"snapshot" comment:"恢复快照"`
}

func (r *Restore) Command() string {
	var cmd = strings.Builder{}
	cmd.WriteString(`restic restore `)
	if r.Snapshot != "" {
		cmd.WriteString(r.Snapshot)
	} else {
		cmd.WriteString(` latest `)
	}
	cmd.WriteString(` --target `)
	cmd.WriteString(filex.Pwd(r.Target))
	cmd.WriteString(` --verbose `)
	if r.Host != "" {
		cmd.WriteString(` --host `)
		cmd.WriteString(r.Host)
	}
	if r.Tag != "" {
		cmd.WriteString(` --tag `)
		cmd.WriteString(r.Tag)
	}
	return cmd.String()
}

// Forget 删除命令
type Forget struct {
	Host     string    `json:"host"     yaml:"host"     comment:"主机名称"`
	Tag      string    `json:"tag"      yaml:"tag"      comment:"标签"`
	Snapshot string    `json:"snapshot" yaml:"snapshot" comment:"删除快照"`
	Strategy *Strategy `json:"strategy" yaml:"strategy" comment:"删除策略"`
}

func (r *Forget) Command() string {
	var cmd = strings.Builder{}
	cmd.WriteString(`restic forget `)
	if r.Snapshot != "" {
		cmd.WriteString(r.Snapshot)
	} else if r.Strategy != nil {
		cmd.WriteString(r.Strategy.Options())
	} else {
		cmd.WriteString(` --keep-last 3 `)
	}
	cmd.WriteString(` --verbose --prune `)
	if r.Host != "" {
		cmd.WriteString(` --host `)
		cmd.WriteString(r.Host)
	}
	if r.Tag != "" {
		cmd.WriteString(` --tag `)
		cmd.WriteString(r.Tag)
	}
	return cmd.String()
}

// Strategy 删除策略
// --keep-last n                    keep the last n snapshots (use 'unlimited' to keep all snapshots)
// --keep-hourly n                  keep the last n hourly snapshots (use 'unlimited' to keep all hourly snapshots)
// --keep-daily n                   keep the last n daily snapshots (use 'unlimited' to keep all daily snapshots)
// --keep-weekly n                  keep the last n weekly snapshots (use 'unlimited' to keep all weekly snapshots)
// --keep-monthly n                 keep the last n monthly snapshots (use 'unlimited' to keep all monthly snapshots)
// --keep-yearly n                  keep the last n yearly snapshots (use 'unlimited' to keep all yearly snapshots)
// --keep-within duration           keep snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-hourly duration    keep hourly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-daily duration     keep daily snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-weekly duration    keep weekly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-monthly duration   keep monthly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-within-yearly duration    keep yearly snapshots that are newer than duration (eg. 1y5m7d2h) relative to the latest snapshot
// --keep-tag tag-list              keep snapshots with this tag0list (can be specified multiple times) (default [])
type Strategy struct {
	KeepLast          string `json:"keep-last"           comment:"保留最近n个快照"`
	KeepHourly        string `json:"keep-hourly"         comment:"保留近n小时内的快照（使用“unlimited”保留所有小时快照）"`
	KeepDaily         string `json:"keep-daily"          comment:"保留近n天内的快照（使用“unlimited”保留所有每日快照）"`
	KeepWeekly        string `json:"keep-weekly"         comment:"保留近n周内的快照（使用“unlimited”保留所有每周快照）"`
	KeepMonthly       string `json:"keep-monthly"        comment:"保留近n月内的快照（使用“unlimited”保留所有每周快照）"`
	KeepYearly        string `json:"keep-yearly"         comment:"保留近n年内的快照（使用“unlimited”保留所有每周快照）"`
	KeepTag           string `json:"keep-tag"            comment:"根据此标签列表保留快照（可以指定多个）（默认[]）"`
	KeepWithin        string `json:"keep-within"         comment:"相对于最新快照，保留时间范围内（例如1y5m7d2h）的快照"`
	KeepWithinHourly  string `json:"keep-within-hourly"  comment:"相对于最新快照，保留时间范围内（例如1y5m7d2h）的每小时快照"`
	KeepWithinDaily   string `json:"keep-within-daily"   comment:"相对于最新快照，保留时间范围内（例如1y5m7d2h）的每日快照"`
	KeepWithinWeekly  string `json:"keep-within-weekly"  comment:"相对于最新快照，保留时间范围内（例如1y5m7d2h）的每周快照"`
	KeepWithinMonthly string `json:"keep-within-monthly" comment:"相对于最新快照，保留时间范围内（例如1y5m7d2h）的每月快照"`
	KeepWithinYearly  string `json:"keep-within-yearly"  comment:"相对于最新快照，保留时间范围内（例如1y5m7d2h）的每年快照"`
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

// Repository 存储库
type Repository struct {
	Uri         string `json:"uri"          yaml:"uri"          comment:"存储库uri"`
	Password    string `json:"password"     yaml:"password"     comment:"存储库密码"`
	AwsUser     string `json:"aws_user"     yaml:"aws_user"     comment:"AmazonS3秘钥ID"`
	AwsPassword string `json:"aws_password" yaml:"aws_password" comment:"AmazonS3秘钥"`
}

// SetEnv 设置环境变量
func (r *Repository) SetEnv() string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf(
			`export RESTIC_REPOSITORY="%s";export RESTIC_PASSWORD="%s";export AWS_ACCESS_KEY_ID="%s";export AWS_SECRET_ACCESS_KEY="%s";`,
			r.Uri, r.Password, r.AwsUser, r.AwsPassword)
	case "windows":
		return fmt.Sprintf(
			`set RESTIC_REPOSITORY=%s&&set RESTIC_PASSWORD=%s&&set AWS_ACCESS_KEY_ID=%s&&set AWS_SECRET_ACCESS_KEY=%s&&`,
			r.Uri, r.Password, r.AwsUser, r.AwsPassword)
	default:
		return ""
	}
}

type Datasource struct {
	Type     string `json:"type"     yaml:"type"     description:"数据库类型（mysql/postgres/mongo）"`
	Host     string `json:"host"     yaml:"host"     description:"主机host"`
	Port     int    `json:"port"     yaml:"port"     description:"端口"`
	Username string `json:"username" yaml:"username" description:"用户名"`
	Password string `json:"password" yaml:"password" description:"密码"`
	Name     string `json:"name"     yaml:"name"     description:"数据库名"`
	Schema   string `json:"schema"   yaml:"schema"   description:"模式名"`
	Table    string `json:"table"    yaml:"table"    description:"表名"`
}

func (d *Datasource) Dump(dir ...string) (dump string, path string) {
	path = filepath.Join("", d.Type)
	if len(dir) > 0 {
		path = filepath.Join(path, dir[0])
	}
	switch d.Type {
	case gormx.MYSQL:
		dump, path = d.mysqlDump(path)
	case gormx.POSTGRES, gormx.PGSQL:
		dump, path = d.pgsqlDump(path)
	case "mongo":
		dump, path = d.mongoDump(path)
	}
	return
}

func (d *Datasource) mysqlDump(path ...string) (string, string) {
	var filePath = d.Name
	var dump = strings.Builder{}
	dump.WriteString(`mysqldump`)
	dump.WriteString(` -h `)
	dump.WriteString(d.Host)
	dump.WriteString(` -P `)
	dump.WriteString(strconv.Itoa(d.Port))
	dump.WriteString(` -u `)
	dump.WriteString(d.Username)
	dump.WriteString(` --password=`)
	dump.WriteString(d.Password)
	dump.WriteString(` `)
	dump.WriteString(d.Name)
	if d.Table != "" {
		dump.WriteString(` `)
		dump.WriteString(d.Table)
		filePath = filepath.Join(filePath, d.Table)
	}
	if len(path) > 0 {
		filePath = filepath.Join(path[0], filePath)
	}
	filePath = filePath + ".sql"
	dump.WriteString(` > `)
	dump.WriteString(filePath)
	return dump.String(), filePath
}

func (d *Datasource) pgsqlDump(path ...string) (string, string) {
	var filePath = d.Name
	var dump = strings.Builder{}
	dump.WriteString(`pg_dump`)
	dump.WriteString(fmt.Sprintf(` "host=%s port=%d user=%s password=%s dbname=%s" `, d.Host, d.Port, d.Username, d.Password, d.Name))
	if d.Schema != "" {
		dump.WriteString(` -n `)
		dump.WriteString(d.Schema)
		filePath = filepath.Join(filePath, d.Schema)
		if d.Table != "" {
			dump.WriteString(` -t `)
			dump.WriteString(d.Table)
			filePath = filepath.Join(filePath, d.Table)
		}
	}
	if len(path) > 0 {
		filePath = filepath.Join(path[0], filePath)
	}
	filePath = filePath + ".sql"
	dump.WriteString(` -f `)
	dump.WriteString(filePath)
	return dump.String(), filePath
}

func (d *Datasource) mongoDump(path ...string) (string, string) {
	var filePath = d.Name
	var dump = strings.Builder{}
	dump.WriteString(`mongodump`)
	dump.WriteString(` -h `)
	dump.WriteString(d.Host)
	dump.WriteString(` --port `)
	dump.WriteString(strconv.Itoa(d.Port))
	dump.WriteString(` -u `)
	dump.WriteString(d.Username)
	dump.WriteString(` -p `)
	dump.WriteString(d.Password)
	dump.WriteString(` -d `)
	dump.WriteString(d.Name)
	if d.Table != "" {
		dump.WriteString(` -c `)
		dump.WriteString(d.Table)
		filePath = filepath.Join(filePath, d.Table)
	}
	if len(path) > 0 {
		filePath = filepath.Join(path[0], filePath)
	}
	dump.WriteString(` -o `)
	dump.WriteString(filePath)
	dump.WriteString(` --gzip`)
	return dump.String(), filePath
}
