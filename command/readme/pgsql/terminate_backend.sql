-- 查询进程
select query,
       'SELECT pg_terminate_backend(' || pid || ');' as kill_sql,
       usename,
       application_name,
       client_addr,
       query_start,
       state
  from pg_stat_activity
 where datname = 'magic'
   and usename = 'postgres'
   and query not like '%PostgreSQL JDBC Driver%'
 order by query_start desc;

-- 杀进程
select pg_terminate_backend(pid);

-- 查询连接数
select count(*), usename from pg_stat_activity group by usename;