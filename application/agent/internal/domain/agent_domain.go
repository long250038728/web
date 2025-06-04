package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/wire"
	"github.com/long250038728/web/application/agent/internal/repository"
	"github.com/long250038728/web/protoc/agent"
	"github.com/long250038728/web/tool/kubernetes"
)

var ProviderSet = wire.NewSet(NewAgentDomain)

type AgentDomain struct {
	repository           *repository.AgentRepository
	kubernetesConfigPath string
}

func NewAgentDomain(repository *repository.AgentRepository) *AgentDomain {
	return &AgentDomain{
		kubernetesConfigPath: "/Users/linlong/Downloads/cls-09eyrddg-config",
		repository:           repository,
	}
}

//http://192.168.1.136:8011/agent/info/logs?ns=dev&name=aristotle-6fc64f487b-c76js&container=aristotle
//http://192.168.1.136:8011/agent/info/events?ns=dev&resource=pod
//http://192.168.1.136:8011/agent/info/resources?ns=dev&resource=Pod

// Logs 获取日志列表(namespace下 pod的log日志)
func (s *AgentDomain) Logs(ctx context.Context, req *agent.LogsRequest) (*agent.LogsResponse, error) {
	if req.Ns == "" || req.Name == "" || req.Container == "" {
		return nil, fmt.Errorf("ns,name,container is required")
	}
	client, err := kubernetes.NewAgent(s.kubernetesConfigPath)
	if err != nil {
		return nil, err
	}
	logs, err := client.GetLogs(ctx, req.Ns, req.Name, req.Container)
	if err != nil {
		return nil, err
	}
	return &agent.LogsResponse{Log: logs}, nil
}

// Events 获取事件列表(namespace下 resource的资事件列表)
func (s *AgentDomain) Events(ctx context.Context, req *agent.EventsRequest) (*agent.EventsResponse, error) {
	if req.Ns == "" || req.Resource == "" {
		return nil, fmt.Errorf("ns,resource is required")
	}
	client, err := kubernetes.NewAgent(s.kubernetesConfigPath)
	if err != nil {
		return nil, err
	}
	events, err := client.GetPodEvents(ctx, req.Resource, req.Ns)
	if err != nil {
		return nil, err
	}

	event := make([]string, 0, len(events))
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return nil, err
		}
		event = append(event, string(b))
	}
	return &agent.EventsResponse{Event: event}, nil
}

// Resources 获取资源列表(namespace下 resource的资源列表)
func (s *AgentDomain) Resources(ctx context.Context, req *agent.ResourcesRequest) (*agent.ResourcesResponse, error) {
	if req.Ns == "" || req.Resource == "" {
		return nil, fmt.Errorf("ns,resource is required")
	}
	client, err := kubernetes.NewAgent(s.kubernetesConfigPath)
	if err != nil {
		return nil, err
	}

	list, err := client.ListResource(ctx, req.Resource, req.Ns)
	if err != nil {
		return nil, err
	}
	event := make([]string, 0, len(list))
	for _, o := range list {
		b, err := json.Marshal(o)
		if err != nil {
			return nil, err
		}
		event = append(event, string(b))
	}

	return &agent.ResourcesResponse{Resource: event}, nil
}
