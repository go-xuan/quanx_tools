### 角色
```shell
read          # 具有此角色的用户可以读取指定数据库中的所有非系统集合的数据。例如，当你只希望某个用户能够查看数据库中的数据，但不能进行修改、删除等操作时，可以为其分配该角色。
readWrite     # 该角色的用户不仅可以读取指定数据库中的所有非系统集合的数据，还可以对这些数据进行插入、更新和删除操作。适用于需要对数据进行常规读写操作的场景。
dbAdmin       # 拥有此角色的用户可以执行数据库管理操作，如创建和删除集合、查看统计信息等，但不具备对数据的读写权限。常用于数据库的日常管理维护。
userAdmin     # 该角色允许用户管理指定数据库中的用户和角色，例如创建、修改和删除用户，以及分配角色等。需要谨慎使用，因为拥有该角色的用户可以控制数据库的访问权限。
clusterAdmin  # 这是最高级别的集群管理角色，拥有此角色的用户可以管理整个 MongoDB 集群，包括分片、复制集等的管理操作。只有在需要对整个集群进行全面管理时才分配该角色。
hostManager   # 该角色允许用户管理 MongoDB 服务器的主机，例如查看和修改服务器的配置参数、监控服务器状态等。
backup        # 拥有该角色的用户可以执行备份操作，例如使用 mongodump 工具备份数据库。
restore       # 该角色允许用户执行恢复操作，例如使用 mongorestore 工具恢复数据库。
root          # 这是最高级别的角色，拥有该角色的用户具有所有的权限，包括对所有数据库和集群的完全控制。在生产环境中，应谨慎使用该角色。
```

### 创建用户
```js
db.createUser({
    user: "root",
    pwd: "root",
    roles: [
        {
            role: "root",
            db: "quanx"
        }
    ],
    authenticationRestrictions: [ ],
    mechanisms: [
        "SCRAM-SHA-1",
        "SCRAM-SHA-256"
    ]
})

db.createUser({
    user: "quanchao",
    pwd: "quanchao",
    roles: [
        {
            role: "readWrite",
            db: "quanx"
        }
    ],
    authenticationRestrictions: [ ],
    mechanisms: [
        "SCRAM-SHA-1",
        "SCRAM-SHA-256"
    ]
})
```

