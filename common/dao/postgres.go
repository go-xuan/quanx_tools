package dao

import (
	"strings"

	"github.com/go-xuan/quanx/core/gormx"
	"github.com/go-xuan/quanx/os/errorx"

	"quanx_tools/common/model"
)

// PgTableList 查询表列表
func PgTableList(source, database, schema string, tables ...string) ([]*model.TableQuery, error) {
	sql := strings.Builder{}
	sql.WriteString(`
select table_catalog as database,
       table_schema as schema,
       table_name as table,
       string_agg(column_name,',' order by ordinal_position) as columns
  from information_schema.columns
 where table_catalog = ?
   and table_schema = ? `)
	if len(tables) > 0 {
		sql.WriteString(` and table_name in ('`)
		sql.WriteString(strings.Join(tables, `','`))
		sql.WriteString(`') `)
	}
	sql.WriteString(` group by table_catalog, table_schema, table_name order by table_name `)
	var result []*model.TableQuery
	if err := gormx.DB(source).Raw(sql.String(), database, schema).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表列表失败")
	}
	return result, nil

}

// PgTableQuery 查询表基本信息
func PgTableQuery(source, database, schema string, table string) (*model.TableQuery, error) {
	sql := strings.Builder{}
	sql.WriteString(`
select table_catalog as database,
       table_schema as schema,
       table_name as table,
       string_agg(column_name,',' order by ordinal_position) as columns
  from information_schema.columns
 where table_catalog = ?
   and table_schema = ?`)
	if len(table) > 0 {
		sql.WriteString(` and table_name = '`)
		sql.WriteString(table)
		sql.WriteString(`' `)
	}
	sql.WriteString(` group by table_catalog, table_schema, table_name order by table_name `)
	var result *model.TableQuery
	if err := gormx.DB(source).Raw(sql.String(), database, schema).Scan(result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表基本信息失败")
	}
	return result, nil
}

// PgTableFieldList 查询表字段列表
func PgTableFieldList(source, database, schema string, tables ...string) ([]*model.TableField, error) {
	sql := strings.Builder{}
	sql.WriteString(`
select t1.column_name as name,
       t1.table_name as table,
       t1.table_schema as schema,
       t1.table_catalog as database,
       t1.udt_name as type,
       t1.column_default as default,
       obj_description(t3.oid) as table_comment,
       t5.description as comment,
       case when t1.numeric_precision is null then t1.character_maximum_length else t1.numeric_precision end as precision,
       t1.numeric_scale as scale,
       t1.is_nullable = 'YES' as nullable
  from information_schema.columns t1
  left join pg_namespace t2 on t1.table_schema = t2.nspname
  left join pg_class t3 on t3.relname = t1.table_name and t3.relnamespace = t2.oid
  left join pg_attribute t4 on t4.attname = t1.column_name and t4.attrelid = t3.oid
  left join pg_description t5 on t5.objoid = t4.attrelid and t5.objsubid = t4.attnum`)
	sql.WriteString(` where t1.table_catalog = '`)
	sql.WriteString(database)
	sql.WriteString(`' `)
	if schema != "" {
		sql.WriteString(` and t1.table_schema = '`)
		sql.WriteString(schema)
		sql.WriteString(`' `)
	}
	if len(tables) > 0 {
		sql.WriteString(` and t1.table_name in ('`)
		sql.WriteString(strings.Join(tables, `','`))
		sql.WriteString(`') `)
	}
	sql.WriteString(` order by t1.table_schema, t1.table_name, t1.ordinal_position`)

	var result []*model.TableField
	if err := gormx.DB(source).Raw(sql.String()).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表字段列表失败")
	}
	return result, nil
}
