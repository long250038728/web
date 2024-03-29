package job

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type SqlJob struct {
	db *gorm.DB
}

func NewSqlJob(db *gorm.DB) *SqlJob {
	return &SqlJob{db}
}

func (j *SqlJob) run(ctx context.Context, t string, sql string) error {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}

	executionTime, err := time.ParseInLocation("2006-01-02 15:04:05", t, location)
	if err != nil {
		return err
	}

	subTime := executionTime.Sub(time.Now())
	if subTime < 0 {
		return errors.New("SubTime Is Error")
	}

	// 等待定时器触发
	timer := time.NewTimer(subTime)
	<-timer.C

	//执行
	res := j.db.Session(&gorm.Session{Context: ctx}).Exec(sql)
	return res.Error
}
