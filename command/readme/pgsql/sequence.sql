-- 序列函数,操作序列
select currval(sequence_name) ;                   -- 返回最近一次用 nextval 获取的指定序列的数值
select nextval(sequence_name);                    -- 递增序列并返回新值
select setval(sequence_name, 10000);              -- 设置当前序列值为10000，调用nextval(seqName)为10001
select setval(sequence_name, 10000, true);        -- 设置当前序列值为10000，调用nextval(seqName)为10001
select setval(sequence_name, 10000, false);       -- 设置当前序列值为10000，调用nextval(seqName)为10000

-- 创建序列关键字含义
-- increment by                                   -- 递增值
-- minvalue                                       -- 序列最小值，no minvalue表示没有最小值
-- maxvalue                                       -- 序列最大值，no maxvalue表示没有最大值
-- start with                                     -- 从几开始
-- cycle                                          -- 是否循环使用，no cycle表示不循环
-- owned by                                       -- 指定到表字段，可以缺省

-- 创建序列
create sequence sequence_name increment by 1 minvalue 0 no maxvalue start with 1 no cycle;
-- 设置表主键为序列
alter sequence sequence_name owned by table_name.id;
-- 更新主键序列值
alter table table_name alter column id set default nextval(sequence_name);

-- 查询序列
select sequence_schema as user_name,
       sequence_name as seq_name,
       currval(sequence_name) as seq_value,
       'drop sequence ' || sequence_name || ';' as drop_sql,
       'select setval(''' || sequence_name || ''', 1, false);' as reset_sql,
       'alter sequence ' || sequence_name || ' owned by ' || sequence_name || ';' as alter_sql,
       'create sequence ' || sequence_name || ' increment 1 minvalue 0 no maxvalue start 1 cache 1 no cycle;' as create_sql
  from information_schema.sequences
 order by sequence_schema, sequence_name;

-- 查询序列
select pu.usename as user_name,
       pc.relname as seq_name,
       nextval(pc.relname ::varchar) as seq_value,
       'drop sequence ' || pc.relname || ';' as drop_sql
  from pg_class pc
  left join pg_user pu
    on pc.relowner = pu.usesysid
 where relkind = 's'
 order by pu.usename, pc.relname;