#!/bin/bash

RED='\033[31m'
YELLOW='\033[33m'
MAGENTA='\033[35m'
CYAN='\033[36m'
RESET='\033[0m'

# 判断是否包含git参数
if_add_git_info() {
  for arg in "$@"; do # 此处$@表示传递进方法内的全部参数
    if [ "$arg" == "git" ]; then
      return 0 # 成功，参数存在
    fi
  done
  return 1 # 失败，参数不存在
}

# 获取参数值
get_param_value() {
  param_name=$1
  param_value=$2
  # 解析命令行参数
  while [ $# -gt 0 ]; do
    if [ "$1" == "-$param_name" ]; then
      if [ -n "$2" ]; then
        param_value="$2"
        break
      fi
    fi
    shift
  done
  echo "$param_value"
  return 0
}

# 获取当前工作目录
go_mod_path=$(pwd)/go.mod
# 检查当前目录是否存在go.mod文件
if [ ! -f "$go_mod_path" ]; then
  echo -e "构建出错：${RED}当前目录下未找到go.mod文件${RESET}"
  exit 1
fi
# 读取go.mod文件的第一行来获取module名
module_name=$(head -n 1 "$go_mod_path" | cut -d ' ' -f 2)
# 检查是否成功获取module名
if [ -z "$module_name" ]; then
  echo -e "构建出错：${RED}获取go.mod文件的module名失败${RESET}"
  exit 1
fi

build_name=$(get_param_value "name" "$(basename "$module_name")" "$@")
echo -e "构建包名：${CYAN}${build_name}${RESET}"

build_os=$(get_param_value "os" "$(uname -s)" "$@")
build_os=$(echo "$build_os" | tr '[:upper:]' '[:lower:]')
echo -e "构建系统：${CYAN}${build_os}${RESET}"

build_arch=$(get_param_value "arch" "$(uname -m)" "$@")
build_arch=$(echo "$build_arch" | tr '[:upper:]' '[:lower:]')
echo -e "构建架构：${CYAN}${build_arch}${RESET}"

# 判断操作系统
case $build_os in
  linux)
    build_path=$GOPATH/bin/$build_name
    export GOOS=linux
    export GOARCH=$build_arch
    ;;
  darwin)
    build_path=$GOPATH/bin/$build_name
    export GOOS=darwin
    export GOARCH=$build_arch
    ;;
  windows|MINGW*|MSYS*|CYGWIN*)
    build_name="$build_name".exe
    build_path=$GOPATH\\bin\\$build_name
    export GOOS=windows
    if [[ "$build_arch" == "x86_64" ]]; then
      export GOARCH=amd64
    fi
    ;;
  *)
    echo -e "操作系统未知: ${RED}${build_os}${RESET}"
    exit 1
    ;;
esac

# 添加git提交信息
if if_add_git_info "$@"; then
  commit_id=$(git rev-parse --short HEAD)
  echo -e "git提交：${YELLOW}${commit_id}${RESET}"
  go_build_cmd="go build -o $build_name -ldflags \"-X main.commitId=$commit_id\""
else
  go_build_cmd="go build -o $build_name"
fi

# 执行go build
echo -e "构建命令：${YELLOW}${go_build_cmd}${RESET}"
sh -c "$go_build_cmd"

# 检查构建是否成功
# shellcheck disable=SC2181
if [ $? == 0 ]; then
    echo -e "构建成功：${MAGENTA}$(pwd)/${build_name}${RESET}"
else
    echo -e "${MAGENTA}构建失败！！！${RESET}"
fi

# 授权并移动到GOPATH/bin
chmod +x "$build_name"
mv "$build_name" "$build_path"
echo "文件移动：$build_path"