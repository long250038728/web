package repository

import (
	"context"
	"github.com/long250038728/web/application/user/protoc"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/long250038728/web/tool/tracing/opentelemetry"
	"time"
)

type UserRepository struct {
	util *app.Util
}

func NewUserRepository(util *app.Util) *UserRepository {
	return &UserRepository{
		util: util,
	}
}

func (r *UserRepository) GetName(ctx context.Context, request *user.user) (string, error) {
	type customer struct {
		Name string `json:"name"`
	}
	c := &customer{}
	//orm
	r.util.Db(ctx).Where("id = ?", 1).Find(c)

	////mq
	//_ = r.util.Mq.Send(ctx, "aaa", "", &mq.Message{Data: []byte("hello")})

	////cache
	//_, _ = r.util.Cache.Set(ctx, "hello", "1")
	//_, _ = r.util.Cache.Get(ctx, "hello")

	////lock
	lock, err := r.util.Locker("hello", "123", time.Second*5)
	if err != nil {
		return "", err
	}
	_ = lock.Lock(ctx)
	cancelCtx, cancel := context.WithCancel(ctx)
	go func() {
		err := lock.AutoRefresh(cancelCtx)
		span := opentelemetry.NewSpan(ctx, "AutoRefresh")
		defer span.Close()
		span.AddEvent(err.Error())
	}()
	defer func() {
		cancel()
		lock.UnLock(ctx)
	}()

	////es
	//query := elastic.NewBoolQuery().Must(
	//	elastic.NewTermQuery("merchant_id", 168),
	//	elastic.NewRangeQuery("age").Gt(10).Lte(20),
	//)
	//res, _ := r.util.Es.Search("hello").Query(query).From(0).Size(100).Do(ctx)
	//for _, data := range res.Hits.Hits {
	//	fmt.Println(data.Source)
	//}

	//request http

	client := http.NewClient(http.SetTimeout(time.Second * 5))
	data := map[string]any{
		"merchant_id":      394,
		"merchant_shop_id": 1150,
		"start_date":       "2023-12-01",
		"end_date":         "2023-12-01",
		"field":            "goods_type_id",
		"client_name":      "app",
	}

	_, _, _ = client.Get(ctx, "http://test.zhubaoe.cn:8888/report/sale_report/inventory", data)
	return "hello:" + request.Name + " " + c.Name, nil
}
