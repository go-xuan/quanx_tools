kubectl get all                                                               # 查看所有资源信息xxx
kubectl get namespace                                                         # 查看所有名称空间
kubectl get node                                                              # 查看节点信息
kubectl get pod                                                               # 查看集群中所有pod信息
kubectl get pod -n namespace                                                  # 查看名称空间pod
kubectl get pod [pod_name]                                                    # 查看集群中某个pod的信息
kubectl get pod -o wide                                                       # 查看pod信息与调度情况
kubectl get pod -o yaml                                                       # 查看资源的yaml格式信息
kubectl logs [pod_name]                                                       # 查看Pod所有日志
kubectl logs [pod_name] -f --tail=100                                         # 查看pod最后100条日志
kubectl logs -l app=[server_name] -f --tail=300                               # 查看标签为[server_name]的pod最后300条日志
kubectl logs [pod_name] -f                                                    # 通过pod名称查看pod日志,-f表示实时查看
kubectl logs [pod_name] -c [container_name]                                   # 查看pod中容器的日志
kubectl patch deploy [server_name] -p '{"spec":{"replicas":1}}'               # 设置服务的副本数量为1
kubectl describe pod [pod_name]                                               # 查看pod状态信息
kubectl exec -it [pod_name] -- /bin/sh                                        # 进入pod内部
kubectl exec -it [pod_name] -n [namespace] -- sh                              # 进入pod内部
kubectl set image deploy [server_name] [container_name] [image_name]          # 修改某个服务的镜像
kubectl delete pod [pod_name]                                                 # 删除某个pod
kubectl rollout restart deploy [server_name]                                  # 重启服务
kubectl edit deploy [server_name]                                             # 编辑某个服务的清单
kubectl edit cm [server_name]                                                 # 编辑某个服务的config map
kubectl get svc                                                               # 查看svc
kubectl edit svc [server_name] -n [namespace]                                 # 编辑svc