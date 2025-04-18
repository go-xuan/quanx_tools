package internal

import "strings"

// 查询表字段列表
func mysqlTableFieldQuerySql(database, include, exclude string) string {
	sql := strings.Builder{}
	sql.WriteString(`
select t1.column_name as "name",
       t1.table_name as "table_name",
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
	sql.WriteString(`'`)
	if include != "" {
		include = strings.ReplaceAll(include, " ", "")
		include = strings.ReplaceAll(include, ",", "','")
		sql.WriteString(` and t1.table_name in ('`)
		sql.WriteString(include)
		sql.WriteString(`')`)
	} else if exclude != "" {
		exclude = strings.ReplaceAll(exclude, " ", "")
		exclude = strings.ReplaceAll(exclude, ",", "','")
		sql.WriteString(` and t1.table_name not in ('`)
		sql.WriteString(exclude)
		sql.WriteString(`')`)
	}
	sql.WriteString(` order by t1.table_name, t1.ordinal_position`)
	return sql.String()
}

// 查询表字段列表
func pgsqlTableFieldQuerySql(database, schema, include, exclude string) string {
	sql := strings.Builder{}
	sql.WriteString(`
select t1.column_name as name,
       t1.table_name as table_name,
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
	sql.WriteString(`'`)
	if schema != "" {
		sql.WriteString(` and t1.table_schema = '`)
		sql.WriteString(schema)
		sql.WriteString(`'`)
	}
	if include != "" {
		include = strings.ReplaceAll(include, " ", "")
		include = strings.ReplaceAll(include, ",", "','")
		sql.WriteString(` and t1.table_name in ('`)
		sql.WriteString(include)
		sql.WriteString(`')`)
	} else if exclude != "" {
		exclude = strings.ReplaceAll(exclude, " ", "")
		exclude = strings.ReplaceAll(exclude, ",", "','")
		sql.WriteString(` and t1.table_name not in ('`)
		sql.WriteString(exclude)
		sql.WriteString(`')`)
	}
	sql.WriteString(` order by t1.table_name, t1.ordinal_position`)
	return sql.String()
}
