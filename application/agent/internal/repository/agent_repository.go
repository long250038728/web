package repository

import (
	"github.com/long250038728/web/tool/app"
)

type AgentRepository struct {
	util *app.Util
}

func NewAgentRepository(util *app.Util) *AgentRepository {
	return &AgentRepository{util: util}
}
