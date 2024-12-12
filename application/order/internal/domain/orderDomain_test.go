package domain

import (
	"context"
	"errors"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/app_error"
	"github.com/long250038728/web/tool/server/rpc"
	"testing"
	"time"
)

func TestOrderDomain_OrderDetail(t *testing.T) {
	rpc.SetCircuits([]string{"/user.User/SayHello"})
	d := NewDomain(repository.NewRepository(app.NewUtil()))

	f := func() {
		resp, err := d.OrderDetail(context.Background(), &order.OrderDetailRequest{Id: 11111})
		if errors.Is(err, app_error.CircuitBreaker) {
			resp = &order.OrderDetailResponse{OrderSn: "熔断"}
			err = nil
		}

		if err != nil {
			t.Log("err: ", err)
			return
		}
		t.Log("resp: ", resp)
	}

	for i := 0; i < 100; i++ {
		go f()
	}
	time.Sleep(time.Second * 3)
}
