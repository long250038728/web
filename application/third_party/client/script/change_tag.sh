#!/bin/bash

# 检查是否提供了项目名称作为参数
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <project_name>"
    exit 1
fi

# 从参数中获取项目名称
project_name="$1"
project_dir="${project_name}"  # 使用项目名称构建目录的相对路径
tag_file="./tag"



# 切换到目录
cd "$project_dir"

# 检查是否成功切换目录
if [ $? -ne 0 ]; then
    echo "Error: Failed to change directory to $project_dir."
    exit 1
fi

# 执行git pull以更新本地仓库
git pull
git checkout master
git pull

# 检查tag文件是否存在
if [ ! -f "$tag_file" ]; then
    echo "Error: tag file does not exist in $project_dir."
    exit 1
fi

# 读取tag文件内容
tag_content=$(cat "$tag_file")

# 使用.分割tag内容
IFS='.' read -ra tag_parts <<< "$tag_content"

# 判断tag_parts数组的长度
if [ "${#tag_parts[@]}" -eq 3 ]; then
    # 三位的情况，添加第四位为1
    new_tag="${tag_content}.1"
else
    # 四位的情况，最后一位加一
    last_part=$(( ${tag_parts[3]} + 1 ))
    new_tag="${tag_parts[0]}.${tag_parts[1]}.${tag_parts[2]}.$last_part"
fi

# 输出新的tag内容
echo "New tag: $new_tag"

# 将新的tag内容写回tag文件
echo "$new_tag" > "$tag_file"

# 将新版本的tag文件添加到git的staging area
git add "$tag_file"
git commit -m "tag update"
git push
