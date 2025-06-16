#!/bin/bash

## bash build.sh user /app/config user:v1 web_my-service-network
if [ $# -lt 4 ]; then
    echo "错误：至少需要4个参数"
    echo "用法: $0 服务名 config路径 镜像名 命名空间"
    exit 1
fi

SCRIPT_DIR=$(dirname "$0")
SERVER_NAME=$1
CONFIG_PATH=$2
IMAGE_NAME=$3
NAME_SPACE=$4
BASE_DIR="${SCRIPT_DIR}/app"


# 重启容器
restart_container() {
    local SERVER_NAME="$1"
    local IMAGE_NAME="$2"
    local NAME_SPACE="$3"

    # 获取现有容器 ID 和名称
    local container_ids=$(docker ps -aq --filter "name=^${SERVER_NAME}-")
    local existing_names=$(docker ps -a --filter "name=^${SERVER_NAME}-" --format '{{.Names}}')

    # 启动新的容器
    local new_container_name="${SERVER_NAME}-$(((RANDOM % 1000) + 1))"
    while grep -q "^${new_container_name}$" <<< "$existing_names"; do
        new_container_name="${SERVER_NAME}-$(((RANDOM % 1000) + 1))"
    done

    if ! docker run --rm -itd --name "${new_container_name}" --network="$NAME_SPACE" "$IMAGE_NAME"; then
        echo "Docker 运行失败，终止执行"
        exit 2
    fi

    # 关闭旧的容器
    if [ -n "$container_ids" ]; then
        docker stop $container_ids
    fi
}

## 重启kubernetes
restart_kubernetes() {
    local SERVER_NAME="$1"
    local IMAGE_NAME="$2"
    local NAME_SPACE="$3"

    # 判断 deployment 是否存在，如果不存在则创建
    if ! kubectl get deployment "$SERVER_NAME" -n "$NAME_SPACE" >/dev/null 2>&1; then
        if ! kubectl apply -n "$NAME_SPACE" -f ./kubernetes.yaml; then
            echo "deployment 创建失败"
            exit 2
        fi
    fi

    #
    if ! kubectl set image deployment/"$SERVER_NAME" "container=$IMAGE_NAME" -n "$NAME_SPACE"; then
        echo "deployment 更新失败"
        exit 3
    fi

    return 0
}



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
if ! docker build -f ./app/dockerfile  -t  $IMAGE_NAME . ; then
  echo "Docker 构建失败，终止执行"
  rm -rf "$BASE_DIR"
  exit 2
fi


restart_container "$SERVER_NAME"  "$IMAGE_NAME" "$NAME_SPACE"
# restart_kubernetes "$SERVER_NAME"  "$IMAGE_NAME" "$NAME_SPACE"


# 删除临时文件
rm -rf  $BASE_DIR
exit 0





