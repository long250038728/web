package domain

import (
	"github.com/google/wire"
	"github.com/long250038728/web/application/agent/internal/repository"
)

var ProviderSet = wire.NewSet(NewAgentDomain)

type Agent struct {
	repository *repository.Agent
}

func NewAgentDomain(repository *repository.Agent) *Agent {
	return &Agent{repository: repository}
}
