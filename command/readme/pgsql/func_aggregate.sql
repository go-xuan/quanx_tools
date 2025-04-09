count()                                               -- 计数
sum()                                                 -- 求和
avg()                                                 -- 取平均值
max()                                                 -- 取最大值
min()                                                 -- 取最小值
string_agg(... , separator order by ...)              -- 连接函数，输出字符串，排序可缺省
array_agg(... order by ...)                           -- 连接函数，将多行字段合成数组，排序可缺省
unnest(...)                                           -- 和ARRAY_AGG函数相反，将数组拆成多行

count(...) over(partition by ... order by ...)        -- 分组并给每组计数
sum(col) over(partition by ... order by ...)          -- 分组并取所选列的和
avg(col) over(partition by ... order by ...)          -- 分组并取所选列的平均值
max(col) over(partition by ... order by ...)          -- 分组并取所选列的最大值
min(col) over(partition by ... order by ...)          -- 分组并取所选列的最小值
row_number() over(partition by ... order by ...)      -- 分组并给每组按行编号
rank() over(partition by ... order by ...)            -- 分组并给每组排名,排名支持并列
dense_rank() over(partition by ... order by ...)      -- 分组并给每组排名,排名不支持并列
first_value(col) over(partition by ... order by ...)  -- 分分组并取所选列的第一位值
last_value(col) over(partition by ... order by ...)   -- 分组并取所选列的最后一位值
lag(col,n,m) over(partition by ... order by ...)      -- 分组并取所选列的正序第n行值,n表示偏移位，m表示默认值(可缺省)
lead(col,n,m) over(partition by ... order by ...)     -- 分组并取所选列的倒叙第n行值,n表示偏移位，m表示默认值(可缺省)