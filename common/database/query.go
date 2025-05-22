package database

import (
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/gormx"
	"github.com/go-xuan/quanx/types/stringx"

	"quanx_tools/common/model"
)

func GetQuery(source string) (Query, error) {
	db, config := gormx.GetInstance(source), gormx.GetConfig(source)
	switch config.Type {
	case gormx.POSTGRES, gormx.PGSQL:
		return &Pgsql{
			Database: config.Database,
			Schema:   config.Schema,
			DB:       db,
		}, nil
	case gormx.MYSQL:
		return &Mysql{
			Database: config.Database,
			Schema:   config.Schema,
			DB:       db,
		}, nil
	default:
		return nil, errorx.Errorf("not support type: %s", config.Type)
	}
}

type Query interface {
	GetTable(table string) (*model.Table, error)
	Tables(tables ...string) ([]*model.Table, error)
	Columns(tables ...string) ([]*model.Column, error)
}

// SqlExec 执行sql
func SqlExec(source, sql string) error {
	if err := gormx.GetInstance(source).Exec(sql).Error; err != nil {
		return errorx.Wrap(err, "执行sql失败")
	}
	return nil
}

// GetColumnValues 查询表字段数据
func GetColumnValues(param map[string]string) ([]string, error) {
	conf := &gormx.Config{
		Enable:   true,
		Type:     param["db_type"],
		Host:     param["db_host"],
		Port:     stringx.ParseInt(param["db_port"]),
		Username: param["db_user"],
		Password: param["db_pwd"],
		Database: param["db_name"],
	}
	//初始化数据库连接
	if err := conf.Execute(); err != nil {
		return nil, errorx.Wrap(err, "Database connection failed")
	}
	defer gormx.Close()

	sb := strings.Builder{}
	sb.WriteString(`select distinct `)
	sb.WriteString(param["db_field"])
	sb.WriteString(` from `)
	sb.WriteString(param["db_table"])
	sb.WriteString(" limit 100")
	var result []string
	if err := gormx.GetInstance().Raw(sb.String()).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询sql失败")
	}
	return result, nil
}
