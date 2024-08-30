package repository

import (
	"context"
	"github.com/long250038728/web/application/user/internal/model"
	"github.com/long250038728/web/protoc/user"
	"github.com/long250038728/web/tool/app"
	"github.com/long250038728/web/tool/server/http"
	"github.com/olivere/elastic/v7"
)

type UserRepository struct {
	util *app.Util
}

func NewRepository(util *app.Util) *UserRepository {
	return &UserRepository{
		util: util,
	}
}

func (r *UserRepository) GetName(ctx context.Context, request *user.RequestHello) (string, error) {
	db, err := r.util.Db(ctx)
	if err != nil {
		return "", err
	}

	es, err := r.util.Es()
	if err != nil {
		return "", err
	}

	c := &model.User{}
	//orm
	db.Select("name").Where("id = ?", 1).Find(c)

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

	//es
	query := elastic.NewBoolQuery().Must(
		elastic.NewTermQuery("merchant_id", 240),
		elastic.NewTermQuery("merchant_shop_id", 867),
		elastic.NewRangeQuery("gold_weight").Gte(0).Lte(10000),
		elastic.NewMatchQuery("admin_user_name", "小刘"),
		elastic.NewMatchPhraseQuery("merchant_shop_name", "大"),
	)
	_, _ = es.Search("sale_order_record_report").Query(query).From(0).Size(100).Do(ctx)

	_, _, _ = http.NewClient().Get(ctx, "http://test.zhubaoe.cn:8888/report/sale_report/inventory", map[string]any{
		"merchant_id":      394,
		"merchant_shop_id": 1150,
		"start_date":       "2023-12-01",
		"end_date":         "2023-12-01",
		"field":            "goods_type_id",
		"client_name":      "app",
	})
	_, _, _ = http.NewClient().Post(ctx, "http://test.zhubaoe.cn:9999/", map[string]any{
		"a": "Login",
		"m": "Admin",
		"p": "1",
		"r": "{\"merchant_code\":\"ab190735\",\"user_name\":\"yzt\",\"password\":\"123456\",\"last_admin_id\":\"\",\"last_admin_name\":\"\",\"shift_status\":\"1\"}",
		"t": "00000",
		"v": "2.4.4",
	})

	return "hello:" + request.Name + " " + c.Name, nil
}
