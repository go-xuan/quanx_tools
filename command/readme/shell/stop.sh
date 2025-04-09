#!/bin/bash

default_name='clash'
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
if [ "${p_id}"x != ""x ]
then
	echo "应用${app_name}正在运行中，PID为 : ${p_id}"
	kill -9 "${p_id}"
	echo "${app_name} 已停止"
fi
exit