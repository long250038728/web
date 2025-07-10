package robfig

import "github.com/robfig/cron/v3"

//go get github.com/robfig/cron

type CronJob struct {
	cron *cron.Cron
}

func NewCronJob() *CronJob {
	return &CronJob{
		cron: cron.New(),
	}
}

//# ┌───────────── minute (0–59)
//# │ ┌───────────── hour (0–23)
//# │ │ ┌───────────── day of the month (1–31)
//# │ │ │ ┌───────────── month (1–12)
//# │ │ │ │ ┌───────────── day of the week (0–6) (Sunday to Saturday;
//# │ │ │ │ │                                   7 is also Sunday on some systems)
//# │ │ │ │ │
//# │ │ │ │ │
//# * * * * * <command to execute>

// */n  表示每隔n个时间

func (c *CronJob) AddFunc(spec string, cmd func()) (int, error) {
	id, err := c.cron.AddFunc(spec, cmd)
	return int(id), err
}

func (c *CronJob) Remove(id int) {
	c.cron.Remove(cron.EntryID(id))
}

// Start 开始（内部执行go） 需要阻塞
func (c *CronJob) Start() {
	c.cron.Start()
}

func (c *CronJob) Close() {
	c.cron.Stop()
}
