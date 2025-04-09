## 一、HugeGraph-Server

### 1 下载安装	

新版本下载地址：https://hugegraph.apache.org/cn/docs/download/download/

旧版本下载地址：https://github.com/apache/incubator-hugegraph-doc/wiki/Apache-HugeGraph-(Incubating)-Old-Versions-Download

| **工具**           | **描述**              | **版本**                                                                                                           |
|------------------|---------------------|------------------------------------------------------------------------------------------------------------------|
| HugeGraph-Server | HugeGraph-Server主程序 | [0.12.0](https://github.com/hugegraph/hugegraph/releases/download/v0.12.0/hugegraph-0.12.0.tar.gz)               |
| HugeGraph-Hubble | 基于Web的可视化图形界面       | [1.6.0](https://github.com/hugegraph/hugegraph-hubble/releases/download/v1.6.0/hugegraph-hubble-1.6.0.tar.gz)    |
| HugeGraph-Loader | 数据导入工具              | [0.12.0](https://github.com/hugegraph/hugegraph-loader/releases/download/v0.12.0/hugegraph-loader-0.12.0.tar.gz) |
| HugeGraph-Tools  | 命令行工具集              | [1.6.0](https://github.com/hugegraph/hugegraph-tools/releases/download/v1.6.0/hugegraph-tools-1.6.0.tar.gz)      |

#### 2.1 安装jdk

HugeGraphServer基于jdk-8开发，代码用到了较多jdk-1.8中的类和方法，请自行安装配置

```shell
[root@VM-4-3-centos /]# java -version
java version "1.8.0_333"
Java(TM) SE Runtime Environment (build 1.8.0_333-b02)
Java HotSpot(TM) 64-Bit Server VM (build 25.333-b02, mixed mode)
```

#### 2.2 安装 GCC

- 如果使用的是RocksDB补充，执行gcc --version查看gcc版本
- 需要GLIBCXX_3.4.10及以上版本

```shell
yum install gcc-c++
cd co	
```

### 2 解压安装

```shell
cd /usr/local/install/hugegraph
tar -zxvf hugegraph-0.12.0.tar.gz
```

### 3 配置HugegraphServer

HugeGraphServer 内部集成了 GremlinServer 和 RestServer，而 gremlin-server.yaml 和 rest-server.properties 就是用来配置这两个 Server 的。

- [GremlinServer](http://tinkerpop.apache.org/docs/3.2.3/reference/#gremlin-server)：GremlinServer 接受用户的 gremlin 语句，解析后转而调用 Core 的代码。
- RestServer：提供 restful API，根据不同的 HTTP 请求，调用对应的 Core API，如果用户请求体是 gremlin 语句，则会转发给 GremlinServer，实现对图数据的操作。

#### 3.1 修改gremlin-server.yaml文件

**host**： 部署 GremlinServer 机器的机器名或 IP

**port**： 部署 GremlinServer 机器的端口，端口自己设定，可以使用默认的8182

```shell
vim conf/gremlin-server.yaml

# host and port of gremlin server, need to be consistent with host and port in rest-server.properties
host: 192.168.0.1
port: 8182
```

#### 3.2 修改rest-server.properties文件

```shell
vim conf/rest-server.properties

# 只能本机访问
restserver.url=http://192.168.0.1:8881
# 所有机器可访问
restserver.url=http://0.0.0.0:8881

# gremlin server url, need to be consistent with host and port in gremlin-server.yaml
# 配置项 gremlinserver.url 是 GremlinServer 为 RestServer 提供服务的 url，该配置项默认为 http://localhost:8182，如需修改，需要和 gremlin-server.yaml 中的 host 和 port 相匹配
gremlinserver.url=http://192.168.0.1:8182

graphs=./conf/graphs
```

### 4 启动服务

- 启动分为"首次启动"和"非首次启动"，这么区分是因为在第一次启动前需要初始化后端数据库，然后启动服务。 当人为停掉服务或者其他原因需要再次启动服务时，因为后端数据库是持久化存在的，直接启动服务即可。

- HugeGraphServer 启动时会连接补充备份并试用补充更新版本号，如果未初始化完整版本或者已完成初始化但版本不匹配时（旧版本数据），HugeGraphServer 会启动失败，并提供错误信息。

- 如果需要指定端口号访问HugeGraphServer，请修改./conf/rest-server.properties的restserver.url配置项，修改成机器名或IP地址。

  这里服务器使用的是8881端口：http://192.168.0.1:8881，可以访问该地址

#### 4.1 Memory启动方式(默认)

Memory后端的数据是保存在内存中无法持久化的，不需要初始化后端，这也是唯一一个不需要初始化的后端。

##### 4.1.1 修改 hugegraph.properties

```shell
backend=memory
serializer=text
```

##### 4.1.2 启动服务

```shell
bin/start-hugegraph.sh 

Starting HugeGraphServer...
Connecting to HugeGraphServer (http://192.168.0.1:8881/graphs)........OK
Started

# 提示的 url 与 rest-server.properties 中配置的 restserver.url 一致
```

#### 4.2 RocksDB启动方式(推荐)

RocksDB是一个嵌入式的数据库，不需要手动安装部署，但是依赖GCC，且要求 GCC 版本 >= 4.3.0（GLIBCXX_3.4.10），如不满足，需要提前升级 GCC

##### 4.2.1 创建rocksdb数据文件存储目录

```shell
cd /usr/local/install/hugegraph/hugegraph-0.12.0
mkdir my_database
mkdir my_database/data my_database/wal
```

##### 4.2.2 修改hugegraph.properties

```shell
vim conf/graphs/hugegraph.properties

....
backend=rocksdb
serializer=binary

# 配置上一步事先创建的rocksdb数据文件存储路径
rocksdb.data_path=/usr/local/install/hugegraph/hugegraph-0.12.0/my_database/data
rocksdb.wal_path=/usr/local/install/hugegraph/hugegraph-0.12.0/my_database/wal
```

![img](https://cdn.nlark.com/yuque/0/2022/png/21913226/1644285681870-0fc43394-1b67-4de8-9c74-599f899aee26.png)

##### 4.2.3 首次启动初始化

仅首次启动需要，后续不需要

```shell
# 执行初始化脚本
bin/init-store.sh
# 当看到 Initialization finished.输出信息，即完成了初始化
```

##### 4.2.4 启动服务

```sh
# 执行hugegraph启动脚本
bin/start-hugegraph.sh

# 正常启动打印信息如下：
.....
Starting HugeGraphServer...
Connecting to HugeGraphServer (http://192.168.0.1:8881/graphs)........OK
Started

# 提示的 url 与 rest-server.properties 中配置的 restserver.url 一致即可
```

#### 4.3 其他启动方式(Cassandra、ScyllaDB、HBase)

其他外置数据库，需要自行安装，可参考:https://hugegraph.github.io/hugegraph-doc/quickstart/hugegraph-server.html

### 5 停止服务

```sh
# 执行hugegraph服务关停脚本
bin/stop-hugegraph.sh

# 正常关停打印信息如下：
no crontab for deploy
The HugeGraphServer monitor has been closed
Killing HugeGraphServer(pid 10025)...OK
```

## 二 、HugeGraph-Hubble

### 1 解压安装

```sh
cd /usr/local/install/hugegraph
tar -zxvf hugegraph-hubble-1.6.0.tar.gz
```

### 2 修改配置

如果需要从另外的机器访问，需要修改文件conf/hugegraph-hubble.properties，将server.host改为0.0.0.0，post端口号可自定义

```sh
cd /usr/local/install/hugegraph-hubble-1.6.0/
vim conf/hugegraph-hubble.properties 

# 0.0.0.0标识当前web服务可任意用户端进行访问
server.host=0.0.0.0
server.port=8882
```

### 3 启动Hubble

```shell
# hugegraph-hubble启动脚本
bin/start-hubble.sh 

# 正常启动打印信息如下：
starting HugeGraphHubble.....OK
logging to /usr/local/install/hugegraph-hubble-1.6.0/logs/hugegraph-hubble.log

```

### 4 访问Hubble

使用浏览器访问http:xx.xx.xx.xx:8882，即可访问web页面
这时候需要创建图

| 图ID | 随意取个值，在当前web服务保证唯一即可                                                                              |
|-----|---------------------------------------------------------------------------------------------------|
| 图名称 | hugegraph(首次启动HugeGraph-Server 前，使用 init-store.sh 初始化的图)                                          |
| 主机名 | HugeGraph-Server所在的主机的域名或IP，参考HugeGraph-Server-home}/conf/rest-server.properties中的rest.server-url |
| 端口号 | HugeGraph-Server所配置的端口号，参考HugeGraph-Server-home}/conf/rest-server.properties中的rest.server-url     |
| 用户名 | 空                                                                                                 |
| 密码  | 空                                                                                                 |
