select localtimestamp;                                                  -- 本地时间，不带时区
select now();                                                           -- 本地时间，带时区
select floor(extract(epoch from now()));                                -- 获取时间戳
select date_part('day', now() - '2019-10-18 12:05'::timestamp);         -- 时间差

select to_timestamp('2021-10-08 12:34:56', 'yyyy-mm-dd hh24:mi:ss'),    -- 字符转时间
       to_timestamp(1634838820),                                        -- 时间戳转时间
       to_date('2021-10-08', 'yyyy-mm-dd'),                             -- 字符转日期
       to_date('2021-10-08 12:34:56', 'yyyy-mm-dd hh24:mi:ss'),         -- 字符转日期
       to_date('08 oct 2021', 'dd mon yyyy');                           -- 字符转日期

select to_char(now(), 'yyyy-mm-dd hh24:mi:ss'),                         -- 时间格式化
       to_char(now(), 'yyyy-mm-dd');                                    -- 时间格式化

select date_trunc('year', now()),                                       -- 截取到年
       date_trunc('month', now()),                                      -- 截取到月
       date_trunc('day', now()),                                        -- 截取到日
       date_trunc('hour', now()),                                       -- 截取到时
       date_trunc('min', now()),                                        -- 截取到分
       date_trunc('sec', now())                                         -- 截取到秒


select now() + '1 year',                                                --当前时间加1年
       now() + '1 month',                                               --当前时间加一个月
       now() + '1 day',                                                 --当前时间加一天
       now() + '1 hour',                                                --当前时间加一个小时
       now() + '1 min',                                                 --当前时间加一分钟
       now() + '1 sec',                                                 --加一秒钟
       now() + '1 year 1 month 1 day 1 hour 1 min 1 sec',               --加1年1月1天1时1分1秒

SELECT extract(millennium FROM now()),                                  -- 千年
       extract(epoch FROM now()),                                       -- 时间戳
       extract(century FROM now()),                                     -- 世纪，第几个世纪
       extract(decade FROM now()),                                      -- 十年，第几个十年
       extract(year FROM now()),                                        -- 年份
       extract(quarter FROM now()),                                     -- 季度
       extract(month FROM now()),                                       -- 月份
       extract(week FROM now()),                                        -- 周
       extract(day FROM now()),                                         -- 天
       extract(hour FROM now()),                                        -- 小时
       extract(minute FROM now()),                                      -- 分钟
       extract(second FROM now()),                                      -- 秒,也可以用sec
       extract(dow FROM now()),                                         -- day of week,sunday=0  monday=1
       extract(isodow FROM now()),                                      -- ISO标准，sunday=7  monday=1
       extract(doy FROM localtimestamp);                                -- day of year

show time zone                                                          -- 查看时区
select * from pg_timezone_names;                                        -- 查看支持的时区列表
set time zone 'PRC'                                                     -- 设置成东八区 北京时间  UTC+8