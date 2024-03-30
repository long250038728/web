### consul docker 部署
```
docker pull consul:1.15
docker run --name=consul-test -d -p 8500:8500 consul:1.15 agent -dev -ui -client='0.0.0.0'
```

--name=consul-test 为容器指定一个名称。
-d 表示在后台运行容器。
-p 8500:8500 将容器的8500端口映射到宿主机的8500端口，这样你就可以通过http://localhost:8500/ui访问Consul的UI。
-consul/agent 是启动Consul代理的命令。
-dev 标志启动一个单节点的开发者模式Consul服务器。
-ui 启用Consul的Web UI界面。
-client='0.0.0.0' 允许从任何地址连接到Consul代理。