package kubernetes

import (
	"context"
	"fmt"
	"io"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metaV2 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"strings"
)

type Agent struct {
	agentClient *Client
}

func NewAgent(configPath string) (*Agent, error) {
	c, err := NewClient(configPath)
	if err != nil {
		return nil, err
	}
	return &Agent{agentClient: c}, nil
}

// ====================================【 操作资源 】============================================

func (a *Agent) CreateResource(ctx context.Context, resource, ns, yaml string) error {
	resourceInterface, err := a.getResourceInterface(resource, ns)
	if err != nil {
		return err
	}
	obj := &unstructured.Unstructured{}
	_, _, err = scheme.Codecs.UniversalDeserializer().Decode([]byte(yaml), nil, obj)
	if err != nil {
		return err
	}
	_, err = resourceInterface.Create(ctx, obj, metaV2.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) DeleteResource(ctx context.Context, resource, ns, name string) error {
	resourceInterface, err := a.getResourceInterface(resource, ns)
	if err != nil {
		return err
	}
	return resourceInterface.Delete(ctx, name, metaV2.DeleteOptions{})
}

func (a *Agent) ListResource(ctx context.Context, resource, ns string) ([]runtime.Object, error) {
	mapping, err := a.getResourceMapping(a.agentClient.restMapper, resource)
	if err != nil {
		return nil, err
	}
	informer, _ := a.agentClient.fact.ForResource(mapping.Resource)
	a.agentClient.fact.Start(ctx.Done())
	a.agentClient.fact.WaitForCacheSync(ctx.Done())

	list, _ := informer.Lister().ByNamespace(ns).List(labels.Everything())
	return list, nil
}

//====================================【 log ,events 】============================================

// GetLogs 获取的应用的的log日志
func (a *Agent) GetLogs(ctx context.Context, ns, name, container string) ([]string, error) {
	req := a.agentClient.client.CoreV1().Pods(ns).GetLogs(name, &coreV1.PodLogOptions{Container: container})
	rc, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rc.Close()
	}()

	logData, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	// 使用 strings.Split 将 logData 分割成字符串数组
	logLines := strings.Split(string(logData), "\n")
	return logLines, nil
}

// GetPodEvents 获取的是服务pod的变动日志
func (a *Agent) GetPodEvents(ctx context.Context, resource, ns string) ([]coreV1.Event, error) {
	// 获取events 变更列表
	events, err := a.agentClient.client.CoreV1().Events(ns).List(ctx, metaV2.ListOptions{})
	if err != nil {
		return nil, err
	}

	list := make([]coreV1.Event, 0, len(events.Items))
	for _, event := range events.Items {
		if event.InvolvedObject.Kind != resource {
			continue
		}
		list = append(list, event)
	}
	return list, nil
}

// ====================================【 私有方法 】============================================

// GVK 唯一标识 Kubernetes 资源的`逻辑类型`   应用场景: 编码时通过 Scheme 将 GVK 映射到 Go 语言中的具体结构体（YAML ↔ Go 结构体映射）
// GVR 唯一标识 Kubernetes 资源的`HTTP路径`  应用场景： GVR 向 API Server 发起 RESTful 请求 （通过 RESTMapper 实现 GVK 到 GVR 的映射）（GVK → HTTP 路径）

// getResourceMapping
// 像 kubernetes 这样的工具需要根据用户输入动态适配各种资源，RESTMapper 是关键组件。
// 1. resource 生成GVR 如果生成成功转 GVK ,然后创建meta.RESTMapping对象
// 2. resource 生成GVK ,然后创建meta.RESTMapping对象
func (a *Agent) getResourceMapping(restMapper meta.RESTMapper, resource string) (*meta.RESTMapping, error) {
	//fullySpecifiedGVR  完整的Group-Version-Resource，如果用户输入了精确的资源 REST 表示（如 apps/v1.deployments）。
	//groupResource 如果用户输入的是资源名（如 deployments），只返回 GroupResource
	fullySpecifiedGVR, groupResource := schema.ParseResourceArg(resource)
	gvk := schema.GroupVersionKind{}

	// 通过GVR获取GVK
	if fullySpecifiedGVR != nil {
		gvk, _ = restMapper.KindFor(*fullySpecifiedGVR)
	}
	if gvk.Empty() {
		gvk, _ = restMapper.KindFor(groupResource.WithVersion(""))
	}
	if !gvk.Empty() {
		return restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	}

	//fullySpecifiedGVR  完整的Group-Version-Kind，如果用户输入了精确的资源 REST 表示（如 apps/v1.Deployment）。
	//groupKind 如果用户输入的是资源名（如 deployments），只返回 GroupKind
	fullySpecifiedGVK, groupKind := schema.ParseKindArg(resource)
	if fullySpecifiedGVK == nil {
		gvk := groupKind.WithVersion("")
		fullySpecifiedGVK = &gvk
	}

	if !fullySpecifiedGVK.Empty() {
		if mapping, err := restMapper.RESTMapping(fullySpecifiedGVK.GroupKind(), fullySpecifiedGVK.Version); err == nil {
			return mapping, nil
		}
	}

	mapping, err := restMapper.RESTMapping(groupKind, gvk.Version)
	if err != nil {
		if meta.IsNoMatchError(err) {
			return nil, fmt.Errorf("the server doesn't have a resource type %q", groupResource.Resource)
		}
		return nil, err
	}

	return mapping, nil
}

func (a *Agent) getResourceInterface(resource, ns string) (dynamic.ResourceInterface, error) {
	mapping, err := a.getResourceMapping(a.agentClient.restMapper, resource)
	if err != nil {
		return nil, err
	}

	// 提供 Get, Create, Update, Delete 等方法
	var resourceInterface dynamic.ResourceInterface = a.agentClient.dynamicClient.Resource(mapping.Resource)
	if mapping.Scope.Name() == "namespace" {
		// 判断资源是命名空间级别的还是集群级别的
		resourceInterface = a.agentClient.dynamicClient.Resource(mapping.Resource).Namespace(ns)
	}
	return resourceInterface, nil
}
