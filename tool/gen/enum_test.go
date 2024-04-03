package gen

import (
	"os"
	"testing"
)

func TestEnum_Gen(t *testing.T) {
	var b []byte
	var err error

	str := `
[
    {
        "key": "order_type",
        "comment": "订单类型",
        "items": [
            {
                "key": "sale",
                "value": 1,
                "comment": "订单销售"
            },
            {
                "key": "exchange",
                "value": 2,
                "comment": "订单换购"
            }
        ]
    },
    {
        "key": "order_status",
        "comment": "订单状态",
        "items": [
            {
                "key": "doing",
                "value": 1,
                "comment": "进行中"
            },
            {
                "key": "finish",
                "value": 2,
                "comment": "已完成"
            }
        ]
    }
]
`

	//gen
	if b, err = NewEnumGen().GenStr(str); err != nil {
		t.Error(err)
		return
	}

	//write file
	if err := os.WriteFile("./demo2.go", b, os.ModePerm); err != nil {
		t.Error(err)
		return
	}

	t.Log("ok")
}
