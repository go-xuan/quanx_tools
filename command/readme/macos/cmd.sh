uname -a                                                    # 设备信息
ifconfig en0 | grep "inet " | awk '{print $2}'              # 查看IP
vim ~/.zshrc                                                # 编辑环境变量
brew services list                                          # 查看所有服务
brew services run [服务名]                                   # 单次运行某个服务
brew services start [服务名]                                 # 运行某个服务，并设置开机自动运行。
brew services stop [服务名]                                  # 停止某个服务
brew services restart                                       # 重启某个服务
cd /opt/homebrew/etc                                        # brew配置文件
brew search xxx                                             # brew搜索
brew install xxx                                            # brew安装