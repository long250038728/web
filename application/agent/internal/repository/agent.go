package repository

import (
	"github.com/google/wire"
	"github.com/long250038728/web/tool/app"
)

var ProviderSet = wire.NewSet(NewAgentRepository)

type Agent struct {
	util *app.Util
}

func NewAgentRepository(util *app.Util) *Agent {
	return &Agent{util: util}
}
