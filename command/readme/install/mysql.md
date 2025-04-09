

## 1 下载

官网下载地址：https://downloads.mysql.com/archives/community/

## 2 安装准备

###  2.1 安装包解压

### 2.2 自定义ini文件

在解压文件下(跟bin目录同级)，新建my.ini文件，配置如下，主要修改以下两个配置，其他配置按需修改

**basedir** ：mysq安装目录
**datadir** ：mysql数据库的数据存放目录

```shell
[mysqld]
port=3306                                             # 设置3306端口
basedir=D:\APP\mysql-8.0.31                           # 设置mysql的安装目录 
datadir=D:\APP\mysql-8.0.31\data                      # 设置mysql数据库的数据的存放目录 
max_connections=100                                   # 允许最大连接数
max_connect_errors=10                                 # 允许连接失败的次数
character-set-server=utf8mb4                          # 服务端使用的字符集默认为utf8mb4
default-storage-engine=INNODB                         # 创建新表时将使用的默认存储引擎
default_authentication_plugin=mysql_native_password   # 默认使用“mysql_native_password”插件认证
[mysql]
default-character-set=utf8mb4                         # 设置mysql客户端默认字符集
[client]
port=3306                                             # 设置mysql客户端连接服务端时默认使用的端口
default-character-set=utf8mb4
```

## 3 安装启动

### 3.1 初始化mysql

在bin目录打开cmd终端，执行初始化命令，==将初始化生成的密码保存，后面修改密码要用==

```sh
mysqld --initialize --console

PS D:\APP\mysql-8.0.31\bin> mysqld --initialize --console
2023-06-03T04:06:44.446327Z 0 [Warning] [MY-010918] [Server] 'default_authentication_plugin' is deprecated and will be removed in a future release. Please use authentication_policy instead.
2023-06-03T04:06:44.446342Z 0 [System] [MY-013169] [Server] D:\APP\mysql-8.0.31\bin\mysqld.exe (mysqld 8.0.31) initializing of server in progress as process 24252
2023-06-03T04:06:44.463987Z 1 [System] [MY-013576] [InnoDB] InnoDB initialization has started.
2023-06-03T04:06:44.748436Z 1 [System] [MY-013577] [InnoDB] InnoDB initialization has ended.
2023-06-03T04:06:45.483555Z 6 [Note] [MY-010454] [Server] A temporary password is generated for root@localhost: :X_CNi=K+5Vr
PS D:\APP\mysql-8.0.31\bin>
```

### 3.2 安装mysql

在bin目录下执行安装命令

==如果提示Install/Remove of the Service Denied! ，关闭当前cmd窗口，重新以管理员身份打开cmd再执行命令即可==

```shell
mysqld --install mysql

D:\APP\mysql-8.0.31\bin>mysqld --install mysql
Service successfully installed.
D:\APP\mysql-8.0.31\bin>
```

### 3.3 启动mysql服务

```sh
net start mysql
```

## 4 连接

### 4.1 连接mysql

需要使用3.1初始化生成的秘密

```sh
mysql -u root -p

D:\APP\mysql-8.0.31\bin>mysql -u root -p
Enter password: ************
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 10
Server version: 8.0.31
Copyright (c) 2000, 2022, Oracle and/or its affiliates.
Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.
Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
mysql>
```

### 4.2 修改密码

```sql
alter user 'root'@'localhost' identified by 'root';
```

### 4.3 创建用户

```sql
-- 创建用户并授权
create user 'quan'@'localhost' identified by 'quan';
grant all on *.* to 'quan'@'localhost';


create user 'nacos'@'localhost' identified by 'nacos';
grant all on *.* to 'nacos'@'localhost';
```