package dao

import (
	"strings"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/gormx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/stringx"

	"quanx_tools/common/model"
)

func TableList(source string, tables ...string) ([]*model.TableQuery, error) {
	var config = gormx.GetConfig(source)
	source = config.Source
	switch config.Type {
	case gormx.Mysql:
		return MysqlTableList(source, config.Database, tables...)
	case gormx.Postgres:
		return PgTableList(source, config.Database, config.Schema, tables...)
	default:
		return nil, errorx.Errorf("not support type: %s", config.Type)
	}
}

func TableFieldList(source string, tables ...string) ([]*model.TableField, error) {
	var config = gormx.GetConfig(source)
	source = config.Source
	switch config.Type {
	case gormx.Mysql:
		return MysqlTableFieldList(source, config.Database, tables...)
	case gormx.Postgres:
		return PgTableFieldList(source, config.Database, config.Schema, tables...)
	default:
		return nil, errorx.Errorf("not support type: %s", config.Type)
	}
}

func TableQuery(source string, table string) (*model.TableQuery, error) {
	var config = gormx.GetConfig(source)
	source = config.Source
	switch config.Type {
	case gormx.Mysql:
		return MysqlTableQuery(source, config.Database, table)
	case gormx.Postgres:
		return PgTableQuery(source, config.Database, config.Schema, table)
	default:
		return nil, errorx.Errorf("not support type: %s", config.Type)
	}
}

// SqlExec 执行sql
func SqlExec(name, sql string) (err error) {
	if err = gormx.DB(name).Exec(sql).Error; err != nil {
		return errorx.Wrap(err, "执行sql失败")
	}
	return
}

// GetDBFieldDataList 查询表字段数据
func GetDBFieldDataList(args string) ([]string, error) {
	params := stringx.ParseUrlParams(args)

	//初始化数据库连接
	if err := configx.Execute(&gormx.Config{
		Enable:   true,
		Type:     params["type"],
		Host:     params["host"],
		Port:     stringx.ParseInt(params["port"]),
		Username: params["username"],
		Password: params["password"],
		Database: params["database"],
	}); err != nil {
		return nil, errorx.Wrap(err, "database connection failed")
	}
	defer gormx.Close()

	sb := strings.Builder{}
	sb.WriteString(`select distinct `)
	sb.WriteString(params["field"])
	sb.WriteString(` from `)
	sb.WriteString(params["table"])
	sb.WriteString(" limit 100")
	var result []string
	if err := gormx.DB().Raw(sb.String()).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询sql失败")
	}
	return result, nil
}
