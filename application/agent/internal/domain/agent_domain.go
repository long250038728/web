package domain

import (
	"github.com/google/wire"
	"github.com/long250038728/web/application/agent/internal/repository"
)

var ProviderSet = wire.NewSet(NewAgentDomain)

type AgentDomain struct {
	repository *repository.AgentRepository
}

func NewAgentDomain(repository *repository.AgentRepository) *AgentDomain {
	return &AgentDomain{
		repository: repository,
	}
}
