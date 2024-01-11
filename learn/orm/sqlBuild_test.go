package orm

import (
	"context"
	"testing"
)

func TestQueryBuild_build(t *testing.T) {
	type Customer struct {
		Id   int32  `json:"id"`
		Name string `json:"name"`
	}

	q, err := (&QueryBuild[Customer]{}).
		Fields("`id`,`telephone`").
		TableName("zby_customer").
		build(context.Background())

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(q.Query)
}
