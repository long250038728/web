package orm

import "testing"

func TestQuery(t *testing.T) {
	query := NewBoolQuery().Must(
		Eq("merchant_id", 0, AllowIsZero(true)),
		Neq("merchant_shop_id", 200),
		Gt("goods_id", 2000),
		NewBoolQuery().Should(
			NewBoolQuery().Must(Gte("status", 1), Lt("status", 2)),
			Lte("status", 3),
			Between("status", 4, 5),
		),
		In("order_id", []int32{1, 2, 3, 4, 5}),
		Raw("type = ?", 6),
	)

	if query.IsEmpty() {
		t.Log("is empty")
		return
	}

	sql, args := query.Do()
	t.Log(sql)
	t.Log(args...)
}
