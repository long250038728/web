package domain

import (
	"context"
	"github.com/long250038728/web/application/order/internal/repository"
	"github.com/long250038728/web/protoc/order"
	"github.com/long250038728/web/tool/app"
	"testing"
)

func TestOrderDomain_OrderDetail(t *testing.T) {
	app.InitPathInfo("")
	d := NewDomain(repository.NewRepository(app.NewUtil()))
	resp, err := d.OrderDetail(context.Background(), &order.OrderDetailRequest{Id: 11111})
	t.Log(resp, err)
}
