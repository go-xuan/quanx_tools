nohup xxxxxx > out.log 2>&1 &                                       # 项目后台启动
nginx -c /usr/local/nginx/conf/nginx.conf                           # nginx
bin/startup.sh -m standalone                                        # nacos单机启动
bin/startup.sh -m cluster                                           # nacos分片启动
bin/pg_ctl -D .../pgsql/data -l .../pgsql/log/logfile restart.      # pgsql启动
bin/zkServer.sh start                                               # zookeeper启动
bin/zkServer.sh restart                                             # zookeeper重启
nohup ./bin/kafka-server-start.sh config/server.properties &        # kafka
bin/redis-server etc/redis.conf                                     # redis
MINIO_ACCESS_KEY=minio MINIO_SECRET_KEY=minio nohup ./minio server --address :9000 --console-address :9001 ./data > ./log/minio.log & # minio