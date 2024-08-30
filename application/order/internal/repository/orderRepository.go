package repository

import (
	"context"
	"fmt"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/authorization/session"
	"github.com/long250038728/web/tool/server/http"
)

type OrderRepository struct {
	util *app.Util
}

func NewRepository(util *app.Util) *OrderRepository {
	return &OrderRepository{
		util: util,
	}
}

func (r *OrderRepository) GetName(ctx context.Context, request *user.RequestHello) (string, error) {
	//md, _ := metadata.FromOutgoingContext(ctx)
	//d, _ := json.Marshal(md)
	//fmt.Println(string(d))
	//opentelemetry.NewSpan(ctx, "hello")

	if claims, err := session.GetClaims(ctx); err == nil {
		fmt.Println(claims.Name)
	}

	if session, err := session.GetSession(ctx); err == nil {
		fmt.Println(session.AuthList)
	}

	type customer struct {
		Name string `json:"name"`
	}
	c := &customer{}
	//orm
	db, err := r.util.Db(ctx)
	if err != nil {
		return "", err
	}
	db.Where("id = ?", 1).Find(c)

	////mq
	//_ = r.util.Mq.Send(ctx, "aaa", "", &mq.Message{Data: []byte("hello")})

	////cache
	//_, _ = r.util.Cache.Set(ctx, "hello", "1")
	//_, _ = r.util.Cache.Get(ctx, "hello")

	////lock
	//lock, err := r.util.Locker("hello", "123", time.Second*5)
	//if err != nil {
	//	return "", err
	//}
	//_ = lock.Lock(ctx)
	//_ = lock.UnLock(ctx)

	////es
	//query := elastic.NewBoolQuery().Must(
	//	elastic.NewTermQuery("merchant_id", 168),
	//	elastic.NewRangeQuery("age").Gt(10).Lte(20),
	//)
	//res, _ := r.util.Es.Search("hello").Query(query).From(0).Size(100).Do(ctx)
	//for _, data := range res.Hits.Hits {
	//	fmt.Println(data.Source)
	//}

	data := map[string]any{
		"merchant_id":      394,
		"merchant_shop_id": 1150,
		"start_date":       "2023-12-01",
		"end_date":         "2023-12-01",
		"field":            "goods_type_id",
		"client_name":      "app",
	}

	_, _, _ = http.NewClient().Get(ctx, "http://test.zhubaoe.cn:8888/report/sale_report/inventory", data)

	return "hello:" + request.Name + " " + c.Name, nil
}
