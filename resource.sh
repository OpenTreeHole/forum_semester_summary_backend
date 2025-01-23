#!/bin/bash

# 定义变量
url="http://example.com"  # 请替换成实际的 URL
user_count=5              # 用户数

# 1. 遍历当前目录，删掉所有文件夹（文件不删除）
for dir in */; do
  if [ -d "$dir" ]; then
    rm -rf "$dir"
  fi
done

# 2. 记下当前目录下所有文件（静态资源）
static_files=()
for file in *; do
  if [ -f "$file" ]; then
    static_files+=("$file")
  fi
done

# 3. 新建 user_count 个文件夹，命名为 1 到 user_count
for ((i=1; i<=user_count; i++)); do
  mkdir "$i"
  
  # 4. 使用 curl 获取 html 并保存
  curl -s "$url" -o "$i/$i.html"
  
  # 5. 为文件夹中的静态资源创建符号链接
  for static_file in "${static_files[@]}"; do
    ln -s "$(pwd)/$static_file" "$i/$static_file"
  done
done

echo "脚本执行完毕！"