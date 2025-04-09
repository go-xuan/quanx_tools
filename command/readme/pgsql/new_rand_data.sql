create or replace function new_random_data(table_name_ varchar, num bigint) returns varchar
as

$body$
declare inser_sql varchar(1000);
        column_sql varchar(1000);
        error_position varchar(1000);
        currIndex bigint = 1;

begin
  inserSqlStr = 'insert into '|| table_name_;
  error_position = 'select start';

  -- 查询内置表，获取insert语句主体
select '(' || string_agg(tc.column_name, ',' order by tc.ordinal_position) ||
       ') values (' ||
       string_agg(tc.column_value, ',' order by tc.ordinal_position) || ');'
into column_sql
from (select column_name,
             ordinal_position,
             (case when udt_name = 'varchar' then '''测试造数据varchar'''
                   when udt_name = 'text' then '''测试造数据text'''
                   when udt_name = 'timestamp' then 'new_random_timestamp(''2021-01-01'',''2021-12-31'')'
                   when udt_name = 'date' then 'new_random_date(''2021-01-01'',''2021-12-31'')'
                   when udt_name = 'numeric' then  'random()*100'
                   when udt_name like 'int%' then 'ceil(random()*100)'
                   else  'null'
               end) as column_value
        from information_schema.columns
       where table_name = table_name_
         and column_default is null) tc;

error_position = 'select end';

-- 拼接insert语句
inser_sql = inser_sql || column_sql;

error_position = 'loop start'; -- 开始循环
while curr_index <= num loop
  execute inser_sql; -- 执行sql
  curr_index = curr_index + 1;
end loop;
error_position = 'loop end'; -- 循环结束

return curr_index;

exception when others then return inser_sql;

end;

$body$

language plpgsql;