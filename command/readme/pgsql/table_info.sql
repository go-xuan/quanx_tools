-- 查询表
select * from information_schema.tables where table_schema = 'information_schema';
-- 查询字段
select * from information_schema.columns where table_name = 't_user';
-- 查询视图
select * from information_schema.sequences where sequence_schema = 'public';
-- 查询所有函数/存储过程
select * from information_schema.routines where routine_schema  = 'public';


-- 可查询 字段名、默认值、是否为空、字段类型、字段顺序
select t1.column_name as name,
       t1.table_name as table_name,
       t1.table_schema as schema,
       t1.table_catalog as database,
       t1.udt_name as type,
       t1.column_default as "default",
       obj_description(t3.oid) as table_comment,
       t5.description as comment,
       case when t1.numeric_precision is null then t1.character_maximum_length else t1.numeric_precision end as precision,
       t1.numeric_scale as scale,
       t1.is_nullable = 'YES' as nullable
  from information_schema.columns t1
  left join pg_namespace t2
    on t1.table_schema = t2.nspname
  left join pg_class t3
    on t3.relname = t1.table_name and t3.relnamespace = t2.oid
  left join pg_attribute t4
    on t4.attname = t1.column_name and t4.attrelid = t3.oid
  left join pg_description t5
    on t5.objoid = t4.attrelid and t5.objsubid = t4.attnum
 where t1.table_name = 't_user';


-- 包括 字段名、字段类型(带长度)、是否为空、字段注释
select a.attname as name,
       format_type(a.atttypid, a.atttypmod) as type,
       a.attnotnull as notnull,
       col_description(a.attrelid, a.attnum) as comment
  from pg_class c
  left join pg_attribute a
    on a.attrelid = c.oid
where c.relname = 't_user'
  and a.attnum > 0;
