select coalesce(null,'a',null,'b','c');                             -- 多参数函数，实际上是返回首个非空值，可用于判空
select position('cd' in 'abcdefg1234567a');                         -- 下标/位置
select md5('123456');                                               -- md5
select repeat('abc', 4);                                            -- 重复


select char_length('中文'),                                          -- 文字长度 2
       char_length('abcdefg1234567a'),                              -- 文字长度 15
       bit_length('中文'),                                           -- 字节长度，一个汉字=24bit, 2*24=48
       bit_length('abcdefg1234567a');                               -- 字节长度，一个字符=24bit，15*8=120

select overlay('123456789' placing 'xxxxxx' from 2 for 4),          -- 将下标区间的字符进行替换
       replace('aab-abb-ab-aabb', 'ab', 'xy'),                      -- 字符替换，整体替换
       translate('aab-abb-ab-aabb', 'ab', 'xy');                    -- 翻译，一对一替换

select substring('aaa123456bbb', 4, 6),                             -- substring(x,y,z)：截取x的第y位之后z个字符
       substring('aaa123456bbb', 4),                                -- substring(x,y)：截取x的第y位之后的所有字符
       substr('aaa123456bbb', 4, 6),                                -- 等同于substring
       substr('aaa123456bbb', 4);                                   -- 等同于substring

select left('aaa123456bbb', 3),                                     -- 前截取
       left('aaa123456bbb', -3),                                    -- 后截取
       right('aaa123456bbb', 3),                                    -- 后截取
       right('aaa123456bbb', -3);                                   -- 前截取

select ltrim('aaa123456bbb', 'abc'),                                -- 去除前缀
       rtrim('aaa123456bbb', 'abc'),                                -- 去除后缀
       btrim('aaa123456bbb', 'abc'),                                -- 去除前后缀
       trim(leading 'abc' from 'aaa123456bbb'),                     -- 等同于ltrim
       trim(trailing 'abc' from 'aaa123456bbb'),                    -- 等同于rtrim
       trim(both 'abc' from 'aaa123456bbb'),                        -- 等同于btrim

select lpad('abc', 10, '01'),                                       -- 字符左填充
       rpad('abc', 10, '01');                                       -- 字符右填充

select upper('abc1234567a'),                                        -- 转大写
       lower('abc1234567a');                                        -- 转小写

select regexp_split_to_array('abc_def_ghi_jkl_mn_opq', '_'),        -- 拆分成数组
       string_to_array('abc_def_ghi_jkl_mn_opq', '_'),              -- 拆分成数组
       split_part('abc_def_ghi_jkl_mn_opq', '_', 2);                -- 拆分后取第n个值


