docker ps # 正在运行中的容器
docker ps -a # 所有容器
docker ps -q               # 容器ID
docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Image}}"     # 指定输出格式
# 格式化参数
.ID                 # 容器ID
.Names              # 容器名
.Image              # 镜像名
.tag                # 标签
.Command            # 执行命令
.CreatedAt          # 创建时间
.RunningFor         # 运行时长
.Ports              # 端口
.Status             # 状态
.Size               # 大小
.Mounts             # 挂载卷
.Networks           # 容器使用网络

docker rm -f CONTAINER_ID                                       # 强制删除容器
docker exec -it CONTAINER_ID /bin/bash                          # 进入容器
docker cp CONTAINER_ID:container_path local_path                # 从容器复制文件到宿主机
docker images                                                   # 列出所有镜像
docker rmi image_id                                             # 删除镜像
docker pull alpine:3.18                                         # 拉取alpine镜像
docker pull golang:1.18.10                                      # 拉取golang镜像
docker pull java:8u111                                          # 拉取java镜像
docker pull postgres:11.12                                      # 拉取postgres镜像
docker pull mysql:8.0.33                                        # 拉取mysql
docker pull redis:6.2.5                                         # 拉取redis镜像
docker pull nacos/nacos-server:2.0.3                            # 拉取nacos镜像