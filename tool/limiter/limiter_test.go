package limiter

import (
	"fmt"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	c := make(chan string, 10)

	go func() {
		i := 0

		for {
			select {
			case c <- fmt.Sprintf("num:%d", i):
				fmt.Println("可以成功写入")
			case <-time.After(time.Second):
				fmt.Println("写入失败")
			}

			i++
		}
	}()

	i := 0
	for {
		a := <-c
		fmt.Println(a)
		if i%2 == 0 {
			time.Sleep(time.Second * 10)
		}
		i++
	}
}
