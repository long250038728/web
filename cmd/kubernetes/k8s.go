package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

//go get kubernetes.io/Client-go@v0.28.0
//go get kubernetes.io/api@v0.28.0
//go get kubernetes.io/apimachinery@v0.28.0

//====================================【 获取资源 】============================================

// getGVK GVK 是 Group、Version 和 Kind 的缩写。 GVK 用于标识 Kubernetes 中的每种资源 —————— 资源的类型信息（用于理解和操作资源模型）
func getGVK() error {
	ctx := context.Background()

	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return err
	}

	// 创建 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// 获取命名空间列表
	namespaceList, err := client.CoreV1().Namespaces().List(ctx, v12.ListOptions{})
	if err != nil {
		return err
	}

	// 输出命名空间列表
	fmt.Println("Namespaces in the cluster:")
	for _, ns := range namespaceList.Items {
		fmt.Printf("- %s\n", ns.Name)
	}

	// 获取命名空间下的所有 Pod
	devPods, err := client.CoreV1().Pods("dev").List(ctx, v12.ListOptions{})
	fmt.Println("pod dev in the cluster:")
	for _, pod := range devPods.Items {
		fmt.Printf("- %s\n", pod.Name)
	}

	// 删除 Pod
	err = client.CoreV1().Pods("dev").Delete(ctx, "aristotle-76c7764bc-ht8kw", v12.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// getGVR GVR 是 Group、Version 和 Resource的缩写。 Resource用于HTTP路径中的资源名称 pods、services。  ——————  资源的 REST 访问路径(用于访问和操作 API 服务器)
func getGVR() error {
	ctx := context.Background()

	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return err
	}

	// 创建 Kubernetes 客户端
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	devPods, err := client.Resource(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}).Namespace("dev").List(ctx, v12.ListOptions{})
	if err != nil {
		return err
	}
	fmt.Println(devPods)
	return nil
}

// getRestMapper  用于解析和确定 Kubernetes 资源的元数据信息   ——————  是基于这些信息构建的映射工具，可以看作是信息的“解析器”，它将高层抽象（GVK）转换为实际的 API 端点（GVR），以支持实际操作
func getRestMapper() error {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return err
	}

	// 创建 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	gr, err := restmapper.GetAPIGroupResources(client)
	if err != nil {
		return err
	}

	m := restmapper.NewDiscoveryRESTMapper(gr)

	fullySpecifiedGVR, groupResource := schema.ParseResourceArg("pods")
	gvk := schema.GroupVersionKind{}

	if fullySpecifiedGVR != nil {
		gvk, _ = m.KindFor(*fullySpecifiedGVR)
	}
	if gvk.Empty() {
		gvk, _ = m.KindFor(groupResource.WithVersion(""))
	}

	res, err := m.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	fmt.Println(string(b)) //{"Resource":{"Group":"","Version":"v1","Resource":"pods"},"GroupVersionKind":{"Group":"","Version":"v1","Kind":"Pod"},"Scope":{}}
	return nil
}

//====================================【 watch 】============================================

// watchGVR cache包 每个cache只能监听一种资源，如pods ，service
func watchGVR() error {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return err
	}

	// 创建 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// 创建 ListWatch
	listWatch := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods", "dev", fields.Everything())

	// 创建 Informer (周期)
	informer := cache.NewSharedInformer(listWatch, &v13.Pod{}, 0)
	res, err := informer.AddEventHandler(&PodHandler{})
	if err != nil {
		return err
	}

	fmt.Println(res.HasSynced())
	informer.Run(wait.NeverStop)
	select {}
}

// watchFactoryGVK informers包 允许同时监听多个资源，如pods ，service（通过工厂方式）
func watchFactoryGVK() error {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return err
	}

	// 创建 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	//informer := informers.NewSharedInformerFactory(Client, 0)
	informer := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("dev"))
	pod := informer.Core().V1().Pods()
	service := informer.Core().V1().Services()

	pod.Informer().AddEventHandler(&PodHandler{})
	service.Informer().AddEventHandler(&ServiceHandler{})

	ch := make(chan struct{})
	informer.Start(ch)
	informer.WaitForCacheSync(ch)

	list, _ := pod.Lister().List(labels.Everything()) //Informer 的 List 操作是在本地缓存中完成，加了WaitForCacheSync 方法，用于等待本地缓存的数据完成同步。否则在调用 List 的时候
	fmt.Println(list)

	return nil
}

// watchFactoryGVR informers包 允许同时监听多个资源，如pods ，service（通过工厂方式）
func watchFactoryGVR() error {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		return err
	}

	// 创建 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	//informer := informers.NewSharedInformerFactory(Client, 0)
	informer := informers.NewSharedInformerFactoryWithOptions(client, 0, informers.WithNamespace("dev"))

	pod, err := informer.ForResource(schema.GroupVersionResource{
		Group: "", Version: "v1", Resource: "pods",
	})
	if err != nil {
		return err
	}

	service, err := informer.ForResource(schema.GroupVersionResource{
		Group: "", Version: "v1", Resource: "service",
	})
	if err != nil {
		return err
	}

	pod.Informer().AddEventHandler(&PodHandler{})
	service.Informer().AddEventHandler(&ServiceHandler{})

	ch := make(chan struct{})
	informer.Start(ch)
	informer.WaitForCacheSync(ch)

	list, _ := pod.Lister().List(labels.Everything()) //Informer 的 List 操作是在本地缓存中完成，加了WaitForCacheSync 方法，用于等待本地缓存的数据完成同步。否则在调用 List 的时候
	fmt.Println(list)

	return nil
}

//====================================【 回调 】=============================================

type PodHandler struct{}

func (h *PodHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("PodHandler OnAdd: ", obj.(*v13.Pod).Name, " isInInitialList", isInInitialList)
}
func (h *PodHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("PodHandler OnUpdate: ", oldObj.(*v13.Pod).Name, "   ", newObj.(*v13.Pod).Name)
}
func (h *PodHandler) OnDelete(obj interface{}) {
	fmt.Println("PodHandler OnDelete: ", obj.(*v13.Pod).Name)
}

type ServiceHandler struct{}

func (h *ServiceHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("ServiceHandler OnAdd: ", obj.(*v13.Service).Name, " isInInitialList", isInInitialList)
}
func (h *ServiceHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("ServiceHandler OnUpdate: ", oldObj.(*v13.Service).Name, "   ", newObj.(*v13.Service).Name)
}
func (h *ServiceHandler) OnDelete(obj interface{}) {
	fmt.Println("ServiceHandler OnDelete: ", obj.(*v13.Service).Name)
}
