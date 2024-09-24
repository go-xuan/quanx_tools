package dao

import (
	"strings"

	"github.com/go-xuan/quanx/core/gormx"
	"github.com/go-xuan/quanx/os/errorx"

	"quanx_tools/common/model"
)

// MysqlTableList 查询表列表
func MysqlTableList(source, database string, tables ...string) ([]*model.TableQuery, error) {
	sql := strings.Builder{}
	sql.WriteString(`
select table_catalog as "database",
       table_schema as "schema",
       table_name as "table",
       group_concat(column_name) as "columns"
  from information_schema.columns
 where  table_schema = ? `)
	if len(tables) > 0 {
		sql.WriteString(` and table_name in ('`)
		sql.WriteString(strings.Join(tables, `','`))
		sql.WriteString(`') `)
	}
	sql.WriteString(` group by table_catalog, table_schema, table_name order by table_name `)
	var result []*model.TableQuery
	if err := gormx.DB(source).Raw(sql.String(), database).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表列表失败")
	}
	return result, nil
}

// MysqlTableQuery 查询表基本信息
func MysqlTableQuery(source, database string, table string) (*model.TableQuery, error) {
	sql := strings.Builder{}
	sql.WriteString(`
select table_catalog as "database",
       table_schema as "schema",
       table_name as "table",
       group_concat(column_name) as "columns"
  from information_schema.columns
 where table_schema = ? `)
	if len(table) > 0 {
		sql.WriteString(` and table_name = '`)
		sql.WriteString(table)
		sql.WriteString(`' `)
	}
	sql.WriteString(` group by table_catalog, table_schema, table_name order by table_name`)
	var result *model.TableQuery
	if err := gormx.DB(source).Raw(sql.String(), database).Scan(result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表基本信息失败")
	}
	return result, nil
}

// MysqlTableFieldList 查询表字段列表
func MysqlTableFieldList(source, database string, tables ...string) ([]*model.TableField, error) {
	sql := strings.Builder{}
	sql.WriteString(`
select t1.column_name as "name",
       t1.table_name as "table",
       t1.table_schema as "database",
       t1.data_type as "type",
       if(t1.column_default is null, t1.extra, t1.column_default) as "default",
       t1.column_comment as "comment",
       t2.table_comment as "table_comment",
       if(t1.data_type = 'varchar', t1.character_maximum_length, t1.numeric_precision) as "precision",
       t1.numeric_scale as "scale",
       t1.is_nullable = 'YES' as "nullable"
  from information_schema.columns t1
  left join information_schema.tables t2
    on t1.table_schema = t2.table_schema
   and t1.table_name = t2.table_name`)
	sql.WriteString(` where t1.table_schema = '`)
	sql.WriteString(database)
	sql.WriteString(`' `)
	if len(tables) > 0 {
		sql.WriteString(` and t1.table_name in ('`)
		sql.WriteString(strings.Join(tables, `','`))
		sql.WriteString(`') `)
	}
	sql.WriteString(` order by t1.table_name, t1.ordinal_position`)
	var result []*model.TableField
	if err := gormx.DB(source).Raw(sql.String()).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表字段列表失败")
	}
	return result, nil
}
