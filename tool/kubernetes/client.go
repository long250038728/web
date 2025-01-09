package kubernetes

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	config        *rest.Config
	client        *kubernetes.Clientset
	dynamicClient *dynamic.DynamicClient
	restMapper    meta.RESTMapper
	fact          informers.SharedInformerFactory
}

func NewClient(configPath string) (c *Client, err error) {
	c = &Client{}
	c.config, err = c.getConfig(configPath)
	if err != nil {
		return nil, err
	}
	// 创建 Kubernetes 客户端
	c.client, err = kubernetes.NewForConfig(c.config)
	if err != nil {
		return nil, err
	}
	// 创建 dynamicClient 动态客户端
	c.dynamicClient, err = c.getDynamicClient(c.config)
	if err != nil {
		return nil, err
	}
	// 创建 源映射关系restMapper工具
	c.restMapper, err = c.getRestMapping(c.config)
	if err != nil {
		return nil, err
	}

	c.fact = c.getInformerFactory()
	return c, nil
}

func (c *Client) getConfig(path string) (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", path)
}

// getDynamicClient 获取动态操作GVR客户端，这里的动态是指GVR. Resources。
// Resources  某些自定义资源但在编译时不知道具体的资源类型时，DynamicClient 提供了一种灵活的方式。(如果很明确的话可以直接使用GVK)
func (c *Client) getDynamicClient(config *rest.Config) (*dynamic.DynamicClient, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return dynamicClient, nil
}

// getRestMapping生成资源映射关系restMapper工具
// RESTMapper 的用途 :
//  1. GVK 转 GVR
//  2. 配合dynamic.DynamicClient，动态资源操作
func (c *Client) getRestMapping(config *rest.Config) (meta.RESTMapper, error) {
	// 创建 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Client.Discovery()  用于与 Kubernetes API 服务器的 discovery 子系统交互，检索集群中所有可用资源及其元信息。
	// restmapper.GetAPIGroupResources(cl discovery.DiscoveryInterface) 获取集群中所有 API 组及其资源信息。返回值是一个 []*restmapper.APIGroupResources，描述了每个 API 组的资源列表。
	groupResources, err := restmapper.GetAPIGroupResources(client.Discovery())
	if err != nil {
		return nil, err
	}

	// mapper 动态地将 GVK（Group-Version-Kind）映射到 GVR（Group-Version-Resource），并支持从 GVR 反查到 GVK。
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)
	return mapper, nil
}

func (c *Client) getInformerFactory() informers.SharedInformerFactory {
	fact := informers.NewSharedInformerFactory(c.client, 0) //创建通用informer工厂
	return fact
}
