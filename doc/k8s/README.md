## istio服务网格
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
```


## yaml文件及验证
### 基础配置
yaml文件
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

验证
```
[root@k8s]# kubectl get pod -n normal
NAME                                  READY   STATUS    RESTARTS       AGE
order-deployment-6754b6dcfd-z9dnh     2/2     Running   1 (101m ago)   101m
user-v1-deployment-564f67c956-zccvl   2/2     Running   1 (101m ago)   101m
user-v2-deployment-69f88c998-tfs5t    2/2     Running   2 (101m ago)   101m
[root@k8s-worker-02 yaml]# kubectl get svc -n normal
NAME    TYPE        CLUSTER-IP        EXTERNAL-IP   PORT(S)             AGE
order   ClusterIP   192.168.254.32    <none>        8002/TCP,9002/TCP   101m
user    ClusterIP   192.168.255.124   <none>        8001/TCP,9001/TCP   101m


[root@k8s]# curl 192.168.253.198:8002/order/order/detail
```

### isto + Kubernetes Gateway
安装Kubernetes Gateway
```
## 对namespace为normal 创建一个Kubernetes Gateway(Kubernetes Gateway API CRD 在大多数 Kubernetes 集群上不会默认安装)
kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || { kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v1.2.0" | kubectl apply -f -; }
kubectl annotate gateway normal-gateway networking.istio.io/service-type=ClusterIP --namespace=normal
```

yaml文件
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

验证(LoadBalancer 在normal)
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

## 通过k8s gateway LoadBalancer 进行访问
curl http://123.207.102.51:10000/order/order/detail
```

### isto + istio gateway
安装istio-ingress Gateway
```
## 安装 istio-ingressgateway
istioctl install --set profile=default -y
```

yaml文件
```
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: normal-gateway
  namespace: normal
spec:
  selector:
    istio: ingressgateway # 使用默认的 Istio Ingress Gateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*" # 允许所有主机，或者指定特定域名

---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: gateway-vs
  namespace: normal
spec:
  hosts:
  - "*"
  gateways:
  - normal-gateway
  http:
  - match:
    - uri:
        prefix: /order
    route:
    - destination:
        host: order.normal.svc.cluster.local
        port:
          number: 8002
  - match:
    - uri:
        prefix: /user
    route:
    - destination:
        host: user
        subset: v1
        port:
          number: 8001  # 明确指定端口号
      weight: 20
    - destination:
        host: user
        subset: v2
        port:
          number: 8001  # 明确指定端口号
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
```

验证(LoadBalancer 在istio-system)
```
[root@k8s]# kubectl get all -n istio-system
NAME                                        READY   STATUS    RESTARTS   AGE
pod/istio-ingressgateway-854dd6765d-zdtn8   1/1     Running   0          3m26s
pod/istiod-557d49b9fd-scfhc                 1/1     Running   0          3m30s

NAME                           TYPE           CLUSTER-IP       EXTERNAL-IP     PORT(S)                                      AGE
service/istio-ingressgateway   LoadBalancer   192.168.253.72   81.71.145.142   15021:30431/TCP,80:30891/TCP,443:31597/TCP   3m26s
service/istiod                 ClusterIP      192.168.252.94   <none>          15010/TCP,15012/TCP,443/TCP,15014/TCP        99m

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/istio-ingressgateway   1/1     1            1           3m26s
deployment.apps/istiod                 1/1     1            1           99m

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/istio-ingressgateway-854dd6765d   1         1         1       3m26s
replicaset.apps/istiod-557d49b9fd                 1         1         1       3m31s
replicaset.apps/istiod-58b6458b69                 0         0         0       99m

NAME                                                       REFERENCE                         TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/istio-ingressgateway   Deployment/istio-ingressgateway   3%/80%    1         5         1          3m26s
horizontalpodautoscaler.autoscaling/istiod                 Deployment/istiod                 0%/80%    1         5         1          3m31s

## 通过istio gateway LoadBalancer 进行访问
[root@k8s]# curl 81.71.145.142:80/order/order/detail
```

