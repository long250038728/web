mkdir app

#把整个项目复制到script目录下是为了dockerfile中的copy只允许复制当前路径下的文件
cp -R ./../application  app/application
cp -R ./../config  app/config
cp -R ./../tool  app/tool
cp -R ./../protoc  app/protoc
cp -R ./../go.mod  app/
cp -R ./../go.sum  app/

# 根据服务名 构建用哪个dockerfile 打包成什么tag 暴露的端口号是多少
DOCKER_FILE="./user_dockerfile"
IMAGE_NAME="ccr.ccs.tencentyun.com/zhubaoe/user:latest"
CONTAINER_NAME="user_docker"
HTTP_PORT="8001:8001"
GRPC_PORT="9001:9001"

#打包成镜像
docker build -f ${DOCKER_FILE} -t ${IMAGE_NAME} .

#无需使用进行删除
rm -rf  app

# 上传到hub
docker push ${IMAGE_NAME}

# docker 运行
docker run  -itd  --network my-service-network  -p ${HTTP_PORT}  -p ${GRPC_PORT}  --name ${CONTAINER_NAME} ${IMAGE_NAME}
