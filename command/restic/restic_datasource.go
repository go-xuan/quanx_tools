package restic

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Datasource struct {
	Type     string `json:"type" yaml:"type"`         // 数据库类型（mysql/postgres/mongo）
	Host     string `json:"host" yaml:"host"`         // 主机
	Port     int    `json:"port" yaml:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名
	Schema   string `json:"schema" yaml:"schema"`     // 模式名
	Table    string `json:"table" yaml:"table"`       // 表名
}

func (d *Datasource) Dump() (string, string) {
	var dumpCmd, dumpPath string
	switch d.Type {
	case "mysql":
		dumpCmd, dumpPath = d.mysqlDump()
	case "postgres", "pgsql":
		dumpCmd, dumpPath = d.pgsqlDump()
	case "mongo":
		dumpCmd, dumpPath = d.mongoDump()
	}
	// dumpPath是一个相对路径：{datasource.Type}/{datasource.Database}/{datasource.Table}
	return dumpCmd, dumpPath
}

func (d *Datasource) mysqlDump() (string, string) {
	var dumpPath = filepath.Join(d.Type, d.Database)
	var dumpCmd = strings.Builder{}
	dumpCmd.WriteString(`mysqldump`)
	dumpCmd.WriteString(` -h `)
	dumpCmd.WriteString(d.Host)
	dumpCmd.WriteString(` -P `)
	dumpCmd.WriteString(strconv.Itoa(d.Port))
	dumpCmd.WriteString(` -u `)
	dumpCmd.WriteString(d.Username)
	dumpCmd.WriteString(` --password=`)
	dumpCmd.WriteString(d.Password)
	dumpCmd.WriteString(` `)
	dumpCmd.WriteString(d.Database)
	if d.Table != "" {
		dumpCmd.WriteString(` `)
		dumpCmd.WriteString(d.Table)
		dumpPath = filepath.Join(dumpPath, d.Table)
	}
	dumpPath = dumpPath + ".sql"
	dumpCmd.WriteString(` > `)
	dumpCmd.WriteString(dumpPath)
	return dumpCmd.String(), dumpPath
}

func (d *Datasource) pgsqlDump() (string, string) {
	var dumpPath = filepath.Join(d.Type, d.Database)
	var dumpCmd = strings.Builder{}
	dumpCmd.WriteString(`pg_dump`)
	dumpCmd.WriteString(fmt.Sprintf(` "host=%s port=%d user=%s password=%s dbname=%s"`, d.Host, d.Port, d.Username, d.Password, d.Database))
	if d.Schema != "" {
		dumpCmd.WriteString(` -n `)
		dumpCmd.WriteString(d.Schema)
		dumpPath = filepath.Join(dumpPath, d.Schema)
		if d.Table != "" {
			dumpCmd.WriteString(` -t `)
			dumpCmd.WriteString(d.Table)
			dumpPath = filepath.Join(dumpPath, d.Table)
		}
	}
	dumpPath = dumpPath + ".sql"
	dumpCmd.WriteString(` -f `)
	dumpCmd.WriteString(dumpPath)
	return dumpCmd.String(), dumpPath
}

func (d *Datasource) mongoDump() (string, string) {
	var dumpCmd = strings.Builder{}
	dumpCmd.WriteString(`mongodump`)
	dumpCmd.WriteString(` -h `)
	dumpCmd.WriteString(d.Host)
	dumpCmd.WriteString(` --port `)
	dumpCmd.WriteString(strconv.Itoa(d.Port))
	dumpCmd.WriteString(` -u `)
	dumpCmd.WriteString(d.Username)
	dumpCmd.WriteString(` -p `)
	dumpCmd.WriteString(d.Password)
	dumpCmd.WriteString(` -d `)
	dumpCmd.WriteString(d.Database)
	if d.Table != "" {
		dumpCmd.WriteString(` -c `)
		dumpCmd.WriteString(d.Table)
		// mongodump时，会自动添加database层级文件夹
	}
	dumpCmd.WriteString(` -o `)
	dumpCmd.WriteString(d.Type)
	dumpCmd.WriteString(` --gzip`)
	return dumpCmd.String(), filepath.Join(d.Type, d.Database)
}
