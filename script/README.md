## 通过Jenkins构建

#### 1.安装Jenkins
```bash
docker run -d -p 8080:8080 -p 50000:50000  --name jenkins   -v /Users/linlong/Desktop/jenkin:/var/jenkins_home jenkins/jenkins:lts

# 初次登录需要解锁密码：
docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword
```
* -p 8080:8080: 将容器的8080端口映射到宿主机的8080端口，这是Jenkins的Web界面端口。
* -p 50000:50000: 将容器的50000端口映射到宿主机的50000端口，这是Jenkins用于构建的代理端口。
* -v jenkins_home:/var/jenkins_home: 将宿主机的jenkins_home目录挂载到容器的/var/jenkins_home目录，这是Jenkins的持久化存储位置。

#### 2.选择插件(Dashboard ->  系统管理  ->  插件管理 )
安装gitee插件、other
- Stage View Plugin (Stage可视化展示 pipeline 阶段执行状态)
- Git Parameter Plugin (git选择分支号参数插件)

#### 3.新建Credentials（Dashboard ->  系统管理  ->  凭证）
生成公钥及密钥
密码：对应生成的密钥  （公钥放在gitee项目中）
```bash
ssh-keygen -t rsa -b 2048 -C "250038728@qq.com"
cat ~/.ssh/id_rsa   # 私钥
cat ~/.ssh/id_rsa.pub  # 公钥
```

#### 4.添加节点(Dashboard ->  系统管理  ->  节点列表) ———— 注意节点上需要安装java
##### Launch agents via SSH
* 远程工作目录  /home/zhubaoe/jenkins_workspace  (在执行构建任务时，用于存放所有相关文件和工作空间的目录)
* 启动方式   Launch agents via SSH
* 主机    服务器ip
* Credentials  选择已经创建的Credentials
*  Host Key Verification Strategy    Non verifying Verification Starategy
* 可用性  尽量保持代理在线

##### SSH Username with private key
* 远程工作目录  /home/zhubaoe/jenkins_workspace  (在执行构建任务时，用于存放所有相关文件和工作空间的目录)
* 启动方式   SSH Username with private key
* 主机     服务器ip
* private 通过下面生成的 `cat ~/.ssh/id_rsa` 粘贴（把 `cat ~/.ssh/id_rsa.pub` 粘贴到服务器上的 `~/.ssh/authorized_keys`）
```bash
ssh-keygen -t rsa -b 2048 -C "250038728@qq.com"
```
* 可用性  尽量保持代理在线

```
启动方式
    Launch agent via execution of command on the master （在Jenkins的主节点（master）上执行）
    Launch agents via SSH （SSH协议连接到远程机器执行）
    Let Jenkins control this Windows agent as a Windows service （对Windows环境的，Jenkins会将代理节点作为一个Windows服务来管理和控制）
    通过Java Web启动代理 （网络启动一个Java应用程序）

Host Key Verification Strategy:
    Known hosts file Verification Strategy
    Manually provided key Verification Strategy  
    Manually trusted key Verification Strategy
    Non verifying Verification Strategy
```

#### 5.新增流水线（流水线模式）
##### General
* 构建触发器 （触发构建）
	- gitee webhook触发构建  （gitee项目中的管理webhooks。新建 url填写Jenkins提供的，密码xxxxxxx）
	- 推送代码
	- 触发分支 master
	- Gitee WebHook xxxxxxx密码（点击生成）
* 参数设置：（参数化构建）
	- 参数化构建过程
		- BRANCH
			- 参数类型： 分支
			- 默认值： master
* 流水线
	- 流水线定义： Pipeline script  from SCM
		- SCM: GIT
		- Repository URL:     git address
		- Credentials  选择之前创建的Credentials
		- Branches to build   选择分支或或页面选择如      /master   或  ${BRANCH}
		- 脚本路径 /xxx/xxx    在git address中的那个文件

```
流水线定义：
Pipeline script  from SCM  ：在git项目上面编写
Pipeline script ：在Jenkins上面编写
```

#### 6.Pipeline流水线
服务器是通过参数化构建时选择参数
```groovy
pipeline {
  agent any 
    stages {
        stage('build') {  //1.文件打包
            steps {
                sh '''
                    /var/jenkins_home/go/bin/go env -w GO111MODULE=on
                    /var/jenkins_home/go/bin/go env -w GOPROXY=https://goproxy.io,direct
                    /var/jenkins_home/go/bin/go env -w GOPRIVATE=gitee.com/zhubaoe-go/jordan
                '''
                
                sh 'export GIT_TERMINAL_PROMPT=1'
                sh '/var/jenkins_home/go/bin/go build -o lmcrm services/lmcrm/lmcrm-service/cmd/lmcrm/main.go'
                sh 'chmod +x lmcrm'
            }
        }
        stage('image') { //2.执行对应的shell
            steps {
                sh "bash /var/jenkins_home/workspace/kobe-lmcrm/build-work/lmcrm/build.sh"
            }
        }
        stage('deploy') { //3.在对应的服务器上面运行（$SYSTEM 中是通过ui参数化构建时选择参数）
            steps {
                sh '''ssh $SYSTEM "docker pull ccr.ccs.tencentyun.com/zhubaoe/kobe:lmcrm_$(cat tag) && docker stop lmcrm && docker rm lmcrm && docker run -p 20015:20015 -p 21015:21015 -itd -v /biz-code/configs:/biz-code/configs -v /biz-code/logs:/biz-code/logs -v /etc/localtime:/etc/localtime --name lmcrm --net kong-net $(cat build-work/lmcrm/image)"'''
            }
        }
    }
}
```

在pipeline流水线指定节点
```groovy
pipeline {
    agent none
    stages {
        stage('web-184') {
            agent{
                label 'web-184'   //选择web-128标签的节点
            }
            steps {
                dir("/biz-code/socrates"){   //切换工作目录，在节点拉取最新的代码
                    checkout([$class: 'GitSCM', branches: [[name: 'master']], doGenerateSubmoduleConfigurations: false, extensions: [],  submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'zhubaoe-pwd', url: 'https://gitee.com/zhubaoe/socrates.git']]])
                }
            }
        }
    }
}
```

## 其他
用curl调用的是
```bash
//测试服
curl -X POST https://jenkins.zhubaoe.cn/job/kobe-service-common/buildWithParameters \
--user admin:11739a99e314641a8f7c039db95458f6e1 \
--data-urlencode "BRANCH=check" 
  
 
//获取jenkins队列 
curl -X GET https://jenkins.zhubaoe.cn/queue/api/json \
--user admin:xxxxx 

// gitee branch 创建分支
curl -X POST --header 'Content-Type: application/json;charset=UTF-8' \
 'https://gitee.com/api/v5/repos/zhubaoe/socrates/branches' \
  -d '{"access_token":"xxxxx","refs":"master","branch_name":"hotfix/reshape_20240410"}'

// gitee pr 提交
curl -X POST --header 'Content-Type: application/json;charset=UTF-8' \
'https://gitee.com/api/v5/repos/zhubaoe/socrates/pulls' \
 -d '{"access_token":"xxxxx","title":"feature/sm0407","head":"feature/sm0407","base":"release/v3.5.40"}'

// gitee pr 合并
curl -X PUT --header 'Content-Type: application/json;charset=UTF-8' \
'https://gitee.com/api/v5/repos/zhubaoe/socrates/pulls/1498/merge' \
-d '{"access_token":"xxxxx","merge_method":"merge"}'

```