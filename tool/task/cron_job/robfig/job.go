package robfig

import (
	"github.com/robfig/cron"
)

type CronJob struct {
	cron *cron.Cron
}

func NewCronJob() *CronJob {
	return &CronJob{
		cron: cron.New(),
	}
}

func (c *CronJob) AddFunc(spec string, cmd func()) error {
	return c.cron.AddFunc(spec, cmd)
}

// Start 开始（内部执行go） 需要阻塞
func (c *CronJob) Start() {
	c.cron.Start()
}

func (c *CronJob) Close() {
	c.cron.Stop()
}
