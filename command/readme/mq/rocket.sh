./mqadmin clusterList -n localhost:9876                                                                                 # 查询cluster
./mqadmin topicList -n localhost:9876                                                                                   # 查询topic
./mqadmin topicStatus -n localhost:9876 -t [TOPIC_NAME]                                                                 # 查询topic
./mqadmin consumerProgress -n localhost:9876 -g [GROUP_NAME]                                                            # 查询消费组
./mqadmin consumerStatus -n localhost:9876 -g [GROUP_NAME]                                                              # 消费组状态
./mqadmin resetOffsetByTime -n localhost:9876 -g [GROUP_NAME] -s 1743696000000 -t [TOPIC_NAME]                          # 重置消费组
./mqadmin queryMsgByKey -n localhost:9876 -t [TOPIC_NAME] -k [KEY]                                                      # 查询消息（基于topic）
./mqadmin printMsg -n localhost:9876 -t [TOPIC_NAME] | grep                                                             # 打印消息（基于topic）