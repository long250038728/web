mkdir app

# 多端构建时需要先创建mybuilder
#docker buildx create --name mybuilder --driver docker-container --use

# 根据服务名 构建用哪个dockerfile 打包成什么tag 暴露的端口号是多少
SERVER="user"
VERSION="v1"

HTTP_PORT="8001:8001"
GRPC_PORT="9001:9001"

#把整个项目复制到script目录下是为了dockerfile中的copy只允许复制当前路径下的文件
cp -R ./../application  app/application
cp -R ./../config  app/config
cp -R ./../tool  app/tool
cp -R ./../protoc  app/protoc
cp -R ./../go.mod  app/
cp -R ./../go.sum  app/


DOCKER_FILE="./${SERVER}_dockerfile"
DOCKER_DOMAIN="ccr.ccs.tencentyun.com/linl"
IMAGE_NAME="${DOCKER_DOMAIN}/${SERVER}:${VERSION}"


#打包成镜像
docker buildx build --platform linux/amd64 -f ${DOCKER_FILE} -t ${IMAGE_NAME} .

#无需使用进行删除
rm -rf  app

## 上传到hub
#docker push ${IMAGE_NAME}
#
## docker 运行
#docker run  -itd  --network my-service-network  -p ${HTTP_PORT}  -p ${GRPC_PORT}  --name ${SERVER} ${IMAGE_NAME}


