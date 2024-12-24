#!/bin/bash

# 在script路径下执行(根据服务名 构建用哪个dockerfile 打包成什么tag 暴露的端口号是多少)
# bash  build.sh user v1 8001:8001 9001:9001 my-service-network

SERVER=$1
VERSION=$2

HTTP_PORT=$3
GRPC_PORT=$4
NETWORK=$5

# 检查必要参数是否提供
if [ -z "$SERVER" ] || [ -z "$VERSION" ] || [ -z "$HTTP_PORT" ] || [ -z "$GRPC_PORT" ] || [ -z "$NETWORK" ]; then
  echo "Usage: $0 <SERVER> <VERSION> <HTTP_PORT> <GRPC_PORT> <NETWORK>"
  exit 1
fi


# 获取当前脚本所在的目录
SCRIPT_DIR=$(dirname "$0")
DOCKER_FILE="${SCRIPT_DIR}/${SERVER}/dockerfile"


# 检查 Dockerfile 是否存在
if [ ! -f "$DOCKER_FILE" ]; then
  echo "Error: Dockerfile not found at ${DOCKER_FILE}"
  exit 1
fi


# 创建根目录
BASE_DIR="${SCRIPT_DIR}/app"
mkdir -p $BASE_DIR

# 把整个项目复制到script目录下是为了dockerfile中的copy只允许复制当前路径下的文件
cp -R "${SCRIPT_DIR}/../application"  "${BASE_DIR}/application"
cp -R "${SCRIPT_DIR}/../config"       "${BASE_DIR}/config"
cp -R "${SCRIPT_DIR}/../tool"         "${BASE_DIR}/tool"
cp -R "${SCRIPT_DIR}/../protoc"       "${BASE_DIR}/protoc"
cp -R "${SCRIPT_DIR}/../go.mod"       "${BASE_DIR}/"
cp -R "${SCRIPT_DIR}/../go.sum"       "${BASE_DIR}/"

DOCKER_DOMAIN="ccr.ccs.tencentyun.com/linl"
DOCKER_FILE="${SCRIPT_DIR}/${SERVER}/dockerfile"
IMAGE_NAME="${DOCKER_DOMAIN}/${SERVER}:${VERSION}"


# 多端构建时需要先创建mybuilder
#docker buildx create --name mybuilder --driver docker-container --use
#打包成镜像
docker buildx build --load --platform linux/amd64 -f ${DOCKER_FILE} -t ${IMAGE_NAME} "$BASE_DIR"

##无需使用进行删除
rm -rf  $BASE_DIR

# 上传到hub
docker push ${IMAGE_NAME}

# docker 运行
docker run  -itd  --network ${NETWORK}  -p ${HTTP_PORT}  -p ${GRPC_PORT}  --name ${SERVER} ${IMAGE_NAME}


