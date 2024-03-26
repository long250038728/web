mkdir app

#把整个项目复制到script目录下是为了dockerfile中的copy只允许复制当前路径下的文件
cp -R ./../../../../web  app

#打包成镜像
docker build -t   web:latest .

#无需使用进行删除
rm -rf  app