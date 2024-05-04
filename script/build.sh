mkdir app

#把整个项目复制到script目录下是为了dockerfile中的copy只允许复制当前路径下的文件
cp -R ./../application  app/application
cp -R ./../config  app/config
cp -R ./../tool  app/tool
cp -R ./../go*  app/

#打包成镜像
docker build -t   ccr.ccs.tencentyun.com/zhubaoe/user:latest .

#无需使用进行删除
rm -rf  app

# 上传到hub
#docker push  ccr.ccs.tencentyun.com/zhubaoe/user_a:latest


docker run -it --name user --network my-network ccr.ccs.tencentyun.com/zhubaoe/user:latest
