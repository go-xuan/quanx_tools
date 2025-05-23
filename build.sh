#!/bin/bash

# 获取第一个参数作为应用名
if [[ $1 && $1 != 0 && $1 != "git" ]]
then
    build_name=$1
else
  # 获取当前工作目录
  current_dir=$(pwd)
  # 检查当前目录是否存在go.mod文件
  if [ ! -f "$current_dir/go.mod" ]; then
    echo "Error: 当前目录下未找到go.mod文件"
  fi
  # 读取go.mod文件的第一行来获取module名
  module_name=$(head -n 1 "$current_dir/go.mod" | cut -d ' ' -f 2)
  # 检查是否成功获取module名
  if [ -z "$module_name" ]; then
    echo "Error: 获取go.mod文件的module名失败"
  fi
  # 只取项目名
  echo "当前项目：$module_name"
  build_name=$(basename "$module_name")
fi

# 检测操作系统和架构
os_type=$(uname -s)
arch_type=$(uname -m)
echo "系统架构：$os_type/$arch_type"

# 判断操作系统
case $os_type in
  Linux)
    bin_path=$GOPATH/bin/$build_name
    export GOOS=linux
    export GOARCH=$arch_type
    ;;
  Darwin)
    bin_path=$GOPATH/bin/$build_name
    export GOOS=darwin
    export GOARCH=$arch_type
    ;;
  MINGW*|MSYS*|CYGWIN*)
    build_name="$build_name".exe
    bin_path=$GOPATH\\bin\\$build_name
    export GOOS=windows
    if [[ "$arch_type" == "x86_64" ]]; then
      export GOARCH=amd64
    fi
    ;;
  *)
    echo "未知的操作系统: $os_type"
    exit 1
    ;;
esac

# 判断是否包含特定参数
add_git_info() {
    for arg in "$@"; do
        if [[ "$arg" == "git" ]]; then
            return 0 # 成功，参数存在
        fi
    done
    return 1 # 失败，参数不存在
}

# 添加git提交信息
if add_git_info "$@"; then
  commit_id=$(git rev-parse --short HEAD)
  echo "当前git提交ID为 $commit_id"
  go_build_cmd="go build -o $build_name -ldflags \"-X main.commitId=$commit_id\""
else
  go_build_cmd="go build -o $build_name"
fi

# 执行go build
echo "构建命令：$go_build_cmd"
sh -c "$go_build_cmd"

# 检查构建是否成功
# shellcheck disable=SC2181
if [ $? == 0 ]; then
    echo "构建成功：$build_name"
else
    echo "构建失败！！！"
fi

# 授权并移动到GOPATH/bin
chmod +x "$build_name"
bin_path=$GOPATH/bin/$build_name
mv "$build_name" "$bin_path"
echo "文件更新：$bin_path"