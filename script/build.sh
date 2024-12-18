mkdir app

# 在script路径下执行
# bash  build.sh user v1 8001:8001 9001:9001 my-service-network

# 多端构建时需要先创建mybuilder
#docker buildx create --name mybuilder --driver docker-container --use

DOCKER_DOMAIN="ccr.ccs.tencentyun.com/linl"

# 根据服务名 构建用哪个dockerfile 打包成什么tag 暴露的端口号是多少
SERVER=$1
VERSION=$2

HTTP_PORT=$3
GRPC_PORT=$4
NETWORK=$5


#把整个项目复制到script目录下是为了dockerfile中的copy只允许复制当前路径下的文件
cp -R ./../application  app/application
cp -R ./../config  app/config
cp -R ./../tool  app/tool
cp -R ./../protoc  app/protoc
cp -R ./../go.mod  app/
cp -R ./../go.sum  app/


DOCKER_FILE="./${SERVER}/dockerfile"
IMAGE_NAME="${DOCKER_DOMAIN}/${SERVER}:${VERSION}"


#打包成镜像
docker buildx build --load --platform linux/amd64 -f ${DOCKER_FILE} -t ${IMAGE_NAME} .

#无需使用进行删除
rm -rf  app

# 上传到hub
docker push ${IMAGE_NAME}

# docker 运行
docker run  -itd  --network ${NETWORK}  -p ${HTTP_PORT}  -p ${GRPC_PORT}  --name ${SERVER} ${IMAGE_NAME}


