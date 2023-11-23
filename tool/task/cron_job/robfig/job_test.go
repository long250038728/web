package robfig

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestNewCronJob(t *testing.T) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	cronJob := NewCronJob()
	_ = do1(cronJob)
	_ = do2(cronJob)

	cronJob.Start()
	defer cronJob.Close()

	//阻塞
	select {
	case s := <-quit:
		fmt.Println(s)
		return
	}
}

func do1(cronJob *CronJob) error {
	return cronJob.AddFunc("*/1 * * * *", func() {
		println("1")
	})
}

func do2(cronJob *CronJob) error {
	return cronJob.AddFunc("1 * * * *", func() {
		println("2")
	})
}
