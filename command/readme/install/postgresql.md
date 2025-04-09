## 1 下载

官网下载地址：https://www.postgresql.org/ftp/source/v11.12/

## 2 解压

```shell
cd /usr/local/install
tar -zxvf postgresql-11.12.tar.gz
```

## 3 编译安装

### 3.1 编译检测(./configure)

./configure用来检测你的安装平台的目标特征的。比如它会检测你是不是有CC（C Complie）或GCC（GNU Compile Collections），并不是需要CC或GCC，它是个shell脚本。使用它之后，会生成一个makefile文件，为下一步的编译做准备。

**prefix=/usr/local/install/pgsql**：表示该软件安装在 /usr/local/install/pgsql 下面，执行文件就会安装在 /usr/local/install/pgsql/bin而不是默认的 /usr/local/bin，资源文件就会安装在 /usr/local/install/pgsql/share，而不是默认的/usr/local/share

```shell
mkdir /usr/local/install/pgsql
cd /usr/local/install/postgresql-11.12
./configure --prefix=/usr/local/install/pgsql
```

### 3.2 编译(make)

make是用来编译的，它从Makefile中读取指令，然后编译。它需要用到Makefile文件；

输出信息中出现：All of PostgreSQL successfully made. Ready to install 即编译成功

```shell
cd /usr/local/install/postgresql-11.12
make
....
All of PostgreSQL successfully made. Ready to install.
```

### 3.3 安装(make intsall)

make intsall是用来安装的，它也从Makefile中读取指令，安装到指定的位置。

输出信息中出现：PostgreSQL installation complete. 即安装成功

```shell
cd /usr/local/install/postgresql-11.12
make install
....
PostgreSQL installation complete.
```

### 3.4 安装目录

--prefix=/usr/local/install/pgsql配置的安装目录是否新增bin、include、lib、share文件

```shell
[root@localhost /]# cd /usr/local/install/pgsql
[root@localhost pgsql]# ll
drwxr-xr-x. 2 root root 4096 2021-09-14 16:05:16 bin
drwxr-xr-x. 6 root root 4096 2021-09-14 16:05:15 include
drwxr-xr-x. 4 root root 4096 2021-09-14 16:05:16 lib
drwxr-xr-x. 6 root root 4096 2021-09-14 16:05:16 share
```

## 4 配置环境变量

```shell
#编辑
vim /etc/profile
 
export PGSQL_HOME=/usr/local/install/pgsql
export PGSQL_DATA=/usr/local/install/pgsql/data
export PGSQL_LIB=/usr/local/install/pgsql/lib
export PATH=$PATH:$PGSQL_HOME/bin

#生效
source /etc/profile
```

## 5 启动

###  5.1 创建用户

```shell
#创建用户组
groupadd postgres
#创建用户
useradd -g postgres postgres
#修改密码
passwd postgres
```

### 5.2 创建目录并授权

使用chown -R postgres:postgres 命令给postgres账号授予目录权限

```shell
#创建数据和日志文件保存目录
cd /usr/local/install/pgsql
mkdir data log

#授权给postgres用户
chown -R postgres:postgres data log
```

### 5.3 初始化数据库

**先切换用户，再初始化数据库**

```shell
# 切换用户
su postgres
#初始化，指定data目录
cd /usr/local/install/pgsql/bin
initdb -D /usr/local/install/pgsql/data
```

输出以下信息即可：

Success. You can now start the database server using:

pg_ctl -D /data/postgres/data/ -l logfile start

```shell
su postgres
initdb -D /usr/local/install/pgsql/data
...
Success. You can now start the database server using:
    pg_ctl -D /usr/local/install/pgsql/data -l logfile start
```

### 5.4 启动数据库

```shell
#启动数据库
pg_ctl -D /usr/local/install/pgsql/data -l /usr/local/install/pgsql/log/logfile start
waiting for server to start.... done
server started

#查看日志文件
cat /usr/local/install/pgsql/log/logfile
2022-11-03 11:31:41.812 CST [4904] LOG:  listening on IPv6 address "::1", port 5432
2022-11-03 11:31:41.812 CST [4904] LOG:  listening on IPv4 address "127.0.0.1", port 5432
2022-11-03 11:31:41.818 CST [4904] LOG:  listening on Unix socket "/tmp/.s.PGSQL.5432"
2022-11-03 11:31:41.839 CST [4905] LOG:  database system was shut down at 2022-11-03 11:30:09 CST
2022-11-03 11:31:41.843 CST [4904] LOG:  database system is ready to accept connections
```

### 5.5 设置监听

修改/usr/local/install/pgsql/data下的pg_hba.conf配置文件，将IPv4 local connections下面的一行改为如下所示：

```
host    all             all             127.0.0.1/32            trust
```

修改/usr/local/install/pgsql/data下的postgresql.conf配置文件，修改为如下所示：

![image-20221103113725911](C:\Users\quanchao\AppData\Roaming\Typora\typora-user-images\image-20221103113725911.png)

### 5.6 重启数据库

```shell
#重启
pg_ctl -D /usr/local/install/pgsql/data -l /usr/local/install/pgsql/log/logfile restart
waiting for server to shut down.... done
server stopped
waiting for server to start.... done
server started
```

## 6 连接数据库

#### 6.1 登录pgsql并重置用户密码

```shell
[postgres@VM-4-3-centos data]$ psql
psql (11.12)
Type "help" for help.

postgres=# alter user postgres with password 'postgres';
ALTER ROLE
postgres=# \q
[postgres@VM-4-3-centos data]$ psql -U postgres -W postgres
Password: 
psql (11.12)
Type "help" for help.

postgres=# 
```

#### 6.1 创建pg用户和数据库并授权

```shell
postgres=# create database quanchao;
CREATE DATABASE                                 ^
postgres=# create user chao_demo with password 'chao_demo';
CREATE ROLE
postgres=# grant all on database chao_demo to chao_demo;
GRANT
```

## 7 安装依赖

```shell
#缺少readline依赖
configure: error: readline library not found
#缺少zlib依赖
configure: error: zlib library not found

# CentOS 安装 readline
yum install readline-devel
# CentOS 安装 zlib
yum install zlib-devel

# Ubantu 安装 readline
sudo apt-get install libreadline-dev
# Ubantu 安装 zlib
sudo apt-get install zlib1g
sudo apt-get install zlib1g.dev
#或
sudo apt-get install zlib1g-dev
sudo apt-get install libx32z1-dev
sudo apt-get install lib32z1
# Ubantu 安装 make
sudo apt-get install make
```

