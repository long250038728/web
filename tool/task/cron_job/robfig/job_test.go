package robfig

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestNewCronJob(t *testing.T) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cronJob := NewCronJob()
	cronJob.Start()
	defer cronJob.Close()

	do1(cronJob)
	do2(cronJob)

	//阻塞
	select {
	case s := <-quit:
		fmt.Println(s)
		return
	}
}

// 秒 分 时 天 月 周几

func do1(cronJob *CronJob) (int, error) {
	return cronJob.AddFunc("*/1 * * * *", func() {
		println("1")
	})
}

func do2(cronJob *CronJob) (int, error) {
	return cronJob.AddFunc("41 * * * *", func() {
		println("2")
		cronJob.Remove(1)
	})
}
