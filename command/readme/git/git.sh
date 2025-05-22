# shellcheck disable=SC2102

git remote -v                                                      # 查看代码仓
git remote rename [repo_old_name] [repo_new_name]                  # 重命名代码仓
git remote rm [repo_name]                                          # 删除代码仓
git remote add [repo_name] [url]                                   # 新增代码仓
git remote set-url [repo_name] [url]                               # 重置代码仓url
git remote prune [repo_name]                                       # 清理远程分支
git remote prune [repo_name] --dry                                 # 查看远程可被清理的分支
git remote show [repo_name]                                        # 查看本地分支和追踪情况

git push -u [repo_name] --all                                      # 提交全部
git push -u [repo_name] --tags                                     # 提交标签
git push [repo_name] [branch_name]                                 # 推送分支到远程，-f 强制推送
git push [repo_name] [tag_name]                                    # 推送标签到远程
git branch                                                         # 查看本地分支
git branch | grep ''                                               # 模糊查找分支
git branch -a                                                      # 查看全部分支
git branch -r                                                      # 查看远程分支
git branch [branch_name]                                           # 创建本地分支
git branch -D [branch_name]                                        # 删除本地分支
git branch -d                                                      # 会在删除前检查merge状态（其与上游分支或者与head）
git branch -D                                                      # 是git branch --delete --force的简写，它会直接删除。
git push [repo_name] --delete [branch_name]                        # 删除远程分支
git fetch -p                                                       # 清理无效分支(远程已删除但是本地没删除的分支)， -p即--prune

git tag [tag_name]                                                 # 创建标签
git tag -a [tag_name] -m "message"                                 # 创建带有注释的标签
git tag -l                                                         # 查看本地标签
git ls-remote                                                      # 查看Git仓库的远程标签
git ls-remote --tags [repo_name]                                   # 查看指定仓库的远程标签
git tag -d [tag_name]                                              # 删除本地标签
git push [repo_name] :refs/tags/[tag_name]                         # 删除远程标签
git push [repo_name] [tag_name]                                    # 推送标签到远程仓库
git push --tags                                                    # 推送所有标签到远程仓库
git fetch --tags                                                   # 查看最新的远程标签

git reset [file]                                                   # 移除提交但是未推送的文件