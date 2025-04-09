#!/bin/bash

now_time=$(date '+%Y%m%d%H%M%S')
default_name=''
log_name=''
app_name=''

# 获取第一个参数作为应用名
if [[ $1 && $1 != 0 ]]
then
    app_name=$1
else
    app_name=$default_name
fi

# 查询应用进程ID
# shellcheck disable=SC2009
p_id=$(ps -ef | grep "${app_name}" |grep -v grep | grep -v bash | awk '{print $2}')
if [ "${p_id}"x == ""x ]
then
	echo "应用${app_name}不在运行中"
	# shellcheck disable=SC2162
	read -p "请选择操作编号: 1-运行/0-退出 : " my_choose
	if [ "${my_choose}"x == "1"x ]
	then
		mv "${log_name}".log "${log_name}"_"${now_time}".log
		nohup ./"${app_name}" > "${log_name}".log 2>&1 &
		sleep 2
		# shellcheck disable=SC2009
		new_pid=$(ps -ef | grep "${app_name}" |grep -v grep | grep -v bash | awk '{print $2}')
		echo "${app_name} 已运行，PID为 : ${new_pid}"
	fi
else
	echo "应用${app_name}正在运行中，PID为 : ${p_id}"
	# shellcheck disable=SC2162
	read -p "请选择操作编号: 1-停止/2-重启/0-退出 : " my_choose
	if [ "${my_choose}"x == "1"x ]
	then
		kill -9 "${p_id}"
		echo "${app_name} 已停止"
	elif [ "${my_choose}"x == "2"x ]
	then
		kill -9 "${p_id}"
		echo "${app_name} 已停止，正在重启中..."
		mv "${app_name}".log "${app_name}"_"${now_time}".log
		nohup ./"${app_name}" > "${log_name}".log 2>&1 &
		sleep 2
		# shellcheck disable=SC2009
		new_pid=$(ps -ef | grep "${app_name}" |grep -v grep | grep -v bash | awk '{print $2}')
		echo "${app_name} 已重启，PID为 : ${new_pid}"
	fi
fi
exit