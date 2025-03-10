package restic

type Config struct {
	Backup        *Backup     `json:"backup" yaml:"backup"`               // 备份
	Restore       *Restore    `json:"restore" yaml:"restore"`             // 恢复
	Forget        *Forget     `json:"forget" yaml:"forget"`               // 删除
	Repository    *Repository `json:"repository" yaml:"repository"`       // 存储库配置
	Datasource    *Datasource `json:"datasource" yaml:"datasource"`       // 备份数据源
	RetryTimes    int         `json:"retryTimes" yaml:"retryTimes"`       // 重试次数
	RetryInterval int         `json:"retryInterval" yaml:"retryInterval"` // 重试间隔(秒)
	Cron          string      `json:"cron" yaml:"cron"`                   // 定时任务
}

func (c *Config) InitRepoCommander() Commander {
	return &InitRepoCommander{
		Repository: c.Repository,
	}
}

func (c *Config) BackupCommander() Commander {
	return &BackupCommander{
		Backup:     c.Backup,
		Repository: c.Repository,
		Datasource: c.Datasource,
	}
}

func (c *Config) RestoreCommander() Commander {
	return &RestoreCommander{
		Restore:    c.Restore,
		Repository: c.Repository,
	}
}

func (c *Config) ForgetCommander() Commander {
	return &ForgetCommander{
		Forget:     c.Forget,
		Repository: c.Repository,
	}
}

func (c *Config) SnapshotsCommander() Commander {
	return &SnapshotsCommander{
		Repository: c.Repository,
	}
}
