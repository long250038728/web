package main

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//go get k8s.io/client-go@v0.28.0
//go get k8s.io/api@v0.28.0
//go get k8s.io/apimachinery@v0.28.0

func main() {
	_ = get()
}

func get() error {
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
	namespaceList, err := client.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	// 输出命名空间列表
	fmt.Println("Namespaces in the cluster:")
	for _, ns := range namespaceList.Items {
		fmt.Printf("- %s\n", ns.Name)
	}

	// 获取命名空间下的所有 Pod
	devPods, err := client.CoreV1().Pods("dev").List(ctx, v1.ListOptions{})
	fmt.Println("pod dev in the cluster:")
	for _, pod := range devPods.Items {
		fmt.Printf("- %s\n", pod.Name)
	}

	// 删除 Pod
	err = client.CoreV1().Pods("dev").Delete(ctx, "aristotle-76c7764bc-ht8kw", v1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
