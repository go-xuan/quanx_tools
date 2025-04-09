package dao

import (
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/gormx"
	"github.com/go-xuan/quanx/types/stringx"

	"quanx_tools/common/model"
)

func TableList(source string, tables ...string) ([]*model.TableQuery, error) {
	var config = gormx.GetConfig(source)
	source = config.Source
	switch config.Type {
	case gormx.MYSQL:
		return MysqlTableList(source, config.Database, tables...)
	case gormx.POSTGRES, gormx.PGSQL:
		return PgTableList(source, config.Database, config.Schema, tables...)
	default:
		return nil, errorx.Errorf("not support type: %s", config.Type)
	}
}

func TableFieldList(source string, tables ...string) ([]*model.TableField, error) {
	var config = gormx.GetConfig(source)
	source = config.Source
	switch config.Type {
	case gormx.MYSQL:
		return MysqlTableFieldList(source, config.Database, tables...)
	case gormx.POSTGRES, gormx.PGSQL:
		return PgTableFieldList(source, config.Database, config.Schema, tables...)
	default:
		return nil, errorx.Errorf("not support type: %s", config.Type)
	}
}

func TableQuery(source string, table string) (*model.TableQuery, error) {
	var config = gormx.GetConfig(source)
	source = config.Source
	switch config.Type {
	case gormx.MYSQL:
		return MysqlTableQuery(source, config.Database, table)
	case gormx.POSTGRES, gormx.PGSQL:
		return PgTableQuery(source, config.Database, config.Schema, table)
	default:
		return nil, errorx.Errorf("not support type: %s", config.Type)
	}
}

// SqlExec 执行sql
func SqlExec(name, sql string) error {
	if err := gormx.GetInstance(name).Exec(sql).Error; err != nil {
		return errorx.Wrap(err, "执行sql失败")
	}
	return nil
}

// GetDBFieldDataList 查询表字段数据
func GetDBFieldDataList(param map[string]string) ([]string, error) {
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
		return nil, errorx.Wrap(err, "database connection failed")
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
