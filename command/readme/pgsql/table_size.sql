-- 查询表空间
select table_name,
       pg_size_pretty(pg_total_relation_size('"' || table_name || '"')) as total_size,
       pg_size_pretty(pg_relation_size('"' || table_name || '"'))       as table_size,
       pg_size_pretty(pg_indexes_size('"' || table_name || '"'))        as index_size
  from information_schema.tables
 where table_schema = 'public'
   and table_catalog = 'magic'
   and table_type = 'BASE TABLE'
 order by pg_total_relation_size('"' || table_name || '"') desc;

-- 清理表空间
vacuum full analyze t_xxxx;
-- 重建索引
reindex table t_xxxx;

-- 查看某个模式大小，包括索引。不包括索引用pg_relation_size
select schemaname,
       pg_size_pretty(SUM(Pg_total_relation_size(schemaname ||'.' ||tablename))) as TOTAL_SIZE,
       pg_size_pretty(SUM(pg_relation_size(schemaname ||'.' ||tablename))) as TABLE_SIZE
from   pg_tables
where  schemaname = 'public'
group  by 1;


select t.table_name, to_char(c.reltuples,'FM99999999999999999') as count_num, t.total_size, t.table_size, t.index_size
  from (select table_name,
               pg_size_pretty(pg_total_relation_size('"' || table_name || '"')) as total_size,
               pg_size_pretty(pg_relation_size('"' || table_name || '"'))       as table_size,
               pg_size_pretty(pg_indexes_size('"' || table_name || '"'))        as index_size
          from information_schema.tables
         where table_schema = 'public'
           and table_catalog = 'magic'
           and table_type = 'BASE TABLE') t
          left join (select relname, reltuples
                       from pg_class r
                       join pg_namespace n on r.relnamespace = n.oid
                      where r.relkind = 'r' and n.nspname = 'public') c
            on t.table_name = c.relname
 order by c.reltuples desc;