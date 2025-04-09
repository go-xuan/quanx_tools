### 1 官网下载

ElasticSearch：https://www.elastic.co/cn/downloads/elasticsearch

kibana：https://www.elastic.co/cn/downloads/kibana

### 2 本地解压

解压到本地安装文件夹，ElasticSearch和Kibana都是直接解压就可以直接启动使用

### 3 使用nssm管理服务

进入nssm安装目录，Shift+鼠标右键 >> 打开PowerShell窗口

```shell
# 安装服务
.\nssm install kibana
.\nssm install ElasticSearch
# 安编辑服务
.\nssm edit kibana
.\nssm edit ElasticSearch
```

### 4 编辑启动配置

ElasticSearch

```yaml
cluster.name: quan-elasticsearch
node.name: quan-localhost
path.data: D:\Data\Elastic\data
path.logs: D:\Data\Elastic\logs
http.port: 9200
```

Kibana

```yaml
server.port: 5601
server.name: "quan-localhost-kibana"
elasticsearch.hosts: ["http://localhost:9200"]
i18n.locale: "zh-CN"
```

