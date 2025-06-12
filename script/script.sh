#!/bin/bash

# 1.复制项目到当前路径下
# 2.进行dockerfile的docker打包
# 3.docker image run / docker push / k8s

## bash script.sh user /app/config user:v1 web_my-service-network
if [ $# -lt 4 ]; then
    echo "错误：至少需要4个参数"
    echo "用法: $0 服务名 config路径 镜像名 命名空间(docker是network、k8s是namespace)"
    exit 1
fi

SCRIPT_DIR=$(dirname "$0")
SERVER_NAME=$1
CONFIG_PATH=$2
IMAGE_NAME=$3
DOCKER_NETWORK=$4
BASE_DIR="${SCRIPT_DIR}/app"

## 创建根目录
mkdir -p $BASE_DIR
mkdir -p "${BASE_DIR}/application"

# 把整个项目复制到script目录下是为了dockerfile中的copy只允许复制当前路径下的文件
cp -R "${SCRIPT_DIR}/../application/${SERVER_NAME}"   "${BASE_DIR}/application/${SERVER_NAME}"
cp -R "${CONFIG_PATH}"                                "${BASE_DIR}/config"
cp -R "${SCRIPT_DIR}/../tool"                         "${BASE_DIR}/tool"
cp -R "${SCRIPT_DIR}/../protoc"                       "${BASE_DIR}/protoc"
cp -R "${SCRIPT_DIR}/../go.mod"                       "${BASE_DIR}/"
cp -R "${SCRIPT_DIR}/../go.sum"                       "${BASE_DIR}/"
cp -R "${SCRIPT_DIR}/${SERVER_NAME}/"                 "${BASE_DIR}/"

# 构建项目
docker build -f ./app/dockerfile  -t  $IMAGE_NAME .

# 删除临时文件
rm -rf  $BASE_DIR

## 运行
container_ids=$(docker ps -aq --filter "name=^${SERVER_NAME}-")
existing_names=$(docker ps -a --filter "name=^${SERVER_NAME}-" --format '{{.Names}}')

#### 启动新的容器
new_container_name="${SERVER_NAME}-$(((RANDOM % 1000) + 1))"
while grep -q "^${new_container_name}$" <<< "$existing_names"; do
    new_container_name="${SERVER_NAME}-$(((RANDOM % 1000) + 1))"
done
docker run --rm -itd --name "${new_container_name}" --network="$DOCKER_NETWORK" "$IMAGE_NAME"


#### 关闭旧的容器
if [ -n "$container_ids" ]; then
    docker stop $container_ids
fi

