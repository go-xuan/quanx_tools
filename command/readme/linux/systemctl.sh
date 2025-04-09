sudo systemctl status xxxx                                          # 查看服务状态
sudo systemctl start xxxx                                           # 启动服务
sudo systemctl stop xxxx                                            # 关闭服务
sudo systemctl restart xxxx                                         # 重启服务
sudo systemctl reload xxxx                                          # 不中断正常功能下重新加载服务
sudo systemctl enable xxxx                                          # 设置服务的开机自启动
sudo systemctl disable xxxx                                         # 关闭服务的开机自启动
sudo systemctl list-units                                           # 查看活跃的单元
sudo systemctl list-unit-files|grep enabled                         # 查看已启动的服务列表
sudo systemctl --failed                                             # 查看启动失败的服务列