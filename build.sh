#!/bin/bash

# 获取第一个参数作为应用名
if [[ $1 && $1 != 0 && $1 != "git" ]]
then
    output_binary=$1
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
  echo "当前项目名: $module_name"
  output_binary=$(basename "$module_name")
fi

# 输出文件名
echo "输出二进制文件: $output_binary"

# 检测操作系统和架构
os_type=$(uname -s)
arch_type=$(uname -m)
if [[ "$os_type" == "Darwin" ]]; then
  export CGO_ENABLED=0
  export GOOS=darwin
  export GOARCH=$arch_type
elif [[ "$os_type" == "Linux" ]]; then
  export CGO_ENABLED=0
  export GOOS=linux
  export GOARCH=$arch_type
else
  output_binary="$output_binary".exe
  os_type='Windows'
  SET CGO_ENABLED=0
  SET GOOS=windows
  SET GOARCH="$arch_type"
fi
echo "系统和架构：$os_type/$arch_type"

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
  go_build_cmd="go build -o $output_binary -ldflags \"-X main.commitId=$commit_id\""
else
  go_build_cmd="go build -o $output_binary"
fi

# 执行go build
echo "构建命令：$go_build_cmd"
sh -c "$go_build_cmd"

# 检查构建是否成功
# shellcheck disable=SC2181
if [ $? == 0 ]; then
    echo "构建成功: $output_binary"
else
    echo "构建失败"
fi

# 授权并移动到GOPATH/bin
chmod +x "$output_binary"
bin_path=$GOPATH/bin/
mv "$output_binary" "$bin_path"
echo "移动 $output_binary 到bin目录 $bin_path"