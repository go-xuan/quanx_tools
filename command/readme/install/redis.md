### 1 下载安装

官网下载地址：https://redis.io/download

#### 1.1 GCC环境

```shell
gcc --version
#安装GCC
yum install gcc-c++
```

#### 1.2 解压编译安装

```shell
# cd /usr/local/install
# tar zxvf redis-6.2.5.tar.gz

#make编译
cd /usr/local/install/redis-6.2.5
make

# 测试编译是否成功(这一步时间会比较长,测试耗时5分钟左右)
make test

#make install编译
make install
```

### 2 配置

#### 2.1 创建文件夹

```shell
cd /usr/local/install/redis-6.2.5
mkdir bin etc data
```

#### 2.2 将redis.conf 移动到etc文件夹

```shell
cd /usr/local/install/redis-6.2.5
mv redis.conf ./etc/
```

#### 2.3 将src内脚本移动到bin文件夹

```shell
cd /usr/local/install/redis-6.2.5/src
mv mkreleasehdr.sh redis-benchmark redis-check-aof redis-cli redis-server ../bin
```

#### 2.4 编辑 redis.conf配置文件

```
vim etc/redis.conf
```

| 备注          | 行数   | 配置          | 参数                                                 |
|-------------|------|-------------|----------------------------------------------------|
| 远程访问Ip      | :75  | bind        | #bind 127.0.0.1                                    |
| redis服务端口   | :98  | port        | port 6379                                          |
| redis后台启用开关 | :257 | daemonize   | daemonize yes                                      |
| pid文件       | :289 | pidfile     | pidfile /var/run/redis_6379.pid                    |
| 日志路径        | :302 | logfile     | logfile "/usr/local/install/redis-6.2.5/redis.log" |
| 改数据路径       | :454 | dir         | dir /usr/local/install/redis-6.2.5/data            |
| redis服务密码   | :901 | requirepass | requirepass Init1234                               |

### 3 启动

#### 3.1 设置Redis开机启动

```shell
vi /etc/rc.d/rc.local
/usr/local/install/redis-6.2.5/bin/redis-server  /usr/local/install/redis-6.2.5/etc/redis.conf
```

#### 3.2 手动启动redis

```shell
cd /usr/local/install/redis-6.2.5
bin/redis-server etc/redis.conf
```

#### 3.3 redis-cli登录

```shell
cd /usr/local/install/redis-6.2.5
# 指定端口登录客户端
bin/redis-cli -p 6379
# 鉴权
127.0.0.1:9601> auth Init1234
OK
# 退出客户端
127.0.0.1:9601> exit
```
