package database

import (
	"gorm.io/gorm"
	"strings"

	"github.com/go-xuan/quanx/base/errorx"

	"quanx_tools/common/model"
)

type Mysql struct {
	Database string
	Schema   string
	DB       *gorm.DB
}

func (q *Mysql) GetTable(table string) (*model.Table, error) {
	var result *model.Table
	if err := q.DB.Raw(`
select table_catalog as "database",
       table_schema as "schema",
       table_name as "name",
       table_comment as "comment"
  from information_schema.tables
 where table_schema = ? 
   and table_name = ?`, q.Database, table).Scan(result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表基本信息失败")
	}
	return result, nil
}

func (q *Mysql) Tables(tables ...string) ([]*model.Table, error) {
	sql := strings.Builder{}
	sql.WriteString(`
select t1.table_catalog as "database",
       t1.table_schema as "schema",
       t1.table_name as "name",
       t1.table_comment as "comment",
       t2.columns as "columns"
  from information_schema.tables t1
  left join (select table_schema,
                    table_name,
                    group_concat(column_name) as columns
               from information_schema.columns
              group by table_schema,table_name) t2
    on t1.table_schema = t2.table_schema
   and t1.table_name = t2.table_name
 where t1.table_schema = ? `)
	if len(tables) > 0 {
		sql.WriteString(` and t1.table_name in ('`)
		sql.WriteString(strings.Join(tables, `','`))
		sql.WriteString(`') `)
	}
	var result []*model.Table
	if err := q.DB.Raw(sql.String(), q.Database).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表列表失败")
	}
	return result, nil
}

func (q *Mysql) Columns(tables ...string) ([]*model.Column, error) {
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
	sql.WriteString(q.Database)
	sql.WriteString(`' `)
	if len(tables) > 0 {
		sql.WriteString(` and t1.table_name in ('`)
		sql.WriteString(strings.Join(tables, `','`))
		sql.WriteString(`') `)
	}
	sql.WriteString(` order by t1.table_name, t1.ordinal_position`)
	var result []*model.Column
	if err := q.DB.Raw(sql.String()).Scan(&result).Error; err != nil {
		return nil, errorx.Wrap(err, "查询表字段列表失败")
	}
	return result, nil
}
