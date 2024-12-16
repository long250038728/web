## istio
### istio安装
```
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.24.1/
export PATH=$PWD/bin:$PATH
```
### istio安装
```
## istio gateways 安装 (namespace为istio-system)
istioctl install -f samples/bookinfo/demo-profile-no-gateways.yaml

## 对namespace为normal进行istio注入
kubectl label namespace normal istio-injection=enabled

## 对namespace为normal 创建一个gateway
kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || { kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v1.2.0" | kubectl apply -f -; }
kubectl annotate gateway normal-gateway networking.istio.io/service-type=ClusterIP --namespace=normal
```

## yaml文件编写
### 基础服务配置
```
apiVersion: v1
kind: Namespace
metadata:
  name: normal

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-v1-deployment
  namespace: normal
spec:
  replicas: 1
  selector:
    matchLabels:
      app: user
      version: v1
  template:
    metadata:
      labels:
        app: user
        version: v1
    spec:
      containers:
      - name: user-container
        image: ccr.ccs.tencentyun.com/linl/user:v1
        ports:
        - containerPort: 8001
        - containerPort: 9001
        imagePullPolicy: Always

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-v2-deployment
  namespace: normal
spec:
  replicas: 1
  selector:
    matchLabels:
      app: user
      version: v2
  template:
    metadata:
      labels:
        app: user
        version: v2
    spec:
      containers:
      - name: user-container
        image: ccr.ccs.tencentyun.com/linl/user:v2
        ports:
        - containerPort: 8001
        - containerPort: 9001
        imagePullPolicy: Always

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-deployment
  namespace: normal
spec:
  replicas: 1
  selector:
    matchLabels:
      app: order
  template:
    metadata:
      labels:
        app: order
    spec:
      containers:
      - name: order-container
        image: ccr.ccs.tencentyun.com/linl/order:v1
        ports:
        - containerPort: 8002
        - containerPort: 9002
        imagePullPolicy: Always

---
apiVersion: v1
kind: Service
metadata:
  name: user
  namespace: normal
spec:
  selector:
    app: user
  ports:
    - protocol: TCP
      port: 8001
      targetPort: 8001
      name: user-port-8001
    - protocol: TCP
      port: 9001
      targetPort: 9001
      name: user-port-9001
  type: ClusterIP

---
apiVersion: v1
kind: Service
metadata:
  name: order
  namespace: normal
spec:
  selector:
    app: order
  ports:
    - protocol: TCP
      port: 8002
      targetPort: 8002
      name: order-port-8002
    - protocol: TCP
      port: 9002
      targetPort: 9002
      name: order-port-9002
  type: ClusterIP
```

### isto
```
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user
  namespace: normal
spec:
  hosts:
  - user
  http:
  - route:
    - destination:
        host: user
        subset: v1  # 指向DestinationRule的subsets中定义的name
      weight: 20
    - destination:
        host: user
        subset: v2  # 指向DestinationRule的subsets中定义的name
      weight: 80

---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: user
  namespace: normal
spec:
  host: user
  subsets:    # 创建两个subsets, v1/v2:通过labels中的version匹配
  - name: v1        
    labels:
      version: v1
  - name: v2
    labels:
      version: v2

---      
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: normal-gateway
  namespace: normal
spec:
  gatewayClassName: istio
  listeners:
  - name: http
    port: 10000  # 为了避免80端口冲突及权限问题，定义为10000
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: Same 

---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: normal
  namespace: normal
spec:
  parentRefs:
  - name: normal-gateway
  rules:
  - matches:
    - path:
        type: PathPrefix   # 改为 PathPrefix
        value: /order       # 允许匹配以 /order 开头的路径
    backendRefs:
    - name: order
      port: 8002
  - matches:
    - path:
        type: PathPrefix
        value: /user
    backendRefs:
    - name: user
      port: 8001
```

## 验证
### kubectl查看
```
[root@k8s]# kubectl get all -n normal
NAME                                        READY   STATUS    RESTARTS       AGE
pod/normal-gateway-istio-84bdd969cf-m7qb7   1/1     Running   0              115m
pod/order-deployment-6754b6dcfd-6vff2       2/2     Running   1 (165m ago)   165m
pod/user-v1-deployment-564f67c956-7wptc     2/2     Running   1 (165m ago)   165m
pod/user-v2-deployment-69f88c998-zvp8j      2/2     Running   2 (165m ago)   165m

NAME                           TYPE           CLUSTER-IP        EXTERNAL-IP      PORT(S)                           AGE
service/normal-gateway-istio   LoadBalancer   192.168.253.250   123.207.102.51   15021:31911/TCP,10000:32266/TCP   129m
service/order                  ClusterIP      192.168.253.198   <none>           8002/TCP,9002/TCP                 166m
service/user                   ClusterIP      192.168.255.24    <none>           8001/TCP,9001/TCP                 166m

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/normal-gateway-istio   1/1     1            1           129m
deployment.apps/order-deployment       1/1     1            1           166m
deployment.apps/user-v1-deployment     1/1     1            1           166m
deployment.apps/user-v2-deployment     1/1     1            1           166m

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/normal-gateway-istio-84bdd969cf   1         1         1       129m
replicaset.apps/order-deployment-6754b6dcfd       1         1         1       166m
replicaset.apps/user-v1-deployment-564f67c956     1         1         1       166m
replicaset.apps/user-v2-deployment-69f88c998      1         1         1       166m
```

### 访问
```
## 通过server进行访问
curl 192.168.253.198:8002/order/order/detail
## 通过istio gateway进行访问
curl http://123.207.102.51:10000/order/order/detail
```