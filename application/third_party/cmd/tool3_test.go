package main

import (
	json2 "encoding/json"
	"testing"
)

//OldMaterialExchangeRelation

func GetGoodsTypeAll() (list []*GoodsType, err error) {
	err = db.Find(&list).Error
	return
}

func TestRecord(t *testing.T) {
	//goodsTypeList, err := GetGoodsTypeAll()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//goodsTypeHash := sliceconv.Map(goodsTypeList, func(item *GoodsType) (key int32, value int32) {
	//	return item.Id, item.SaleChargeType
	//})

	var list []*OldMaterialExchangeRecord
	if err := db.
		//Where("merchant_id = ?", merchantId).
		Find(&list).Error; err != nil {
		t.Error(err)
		return
	}
	for _, r := range list {
		var relations []*OldMaterialExchangeRelation
		if err := json2.Unmarshal([]byte(r.Data), &relations); err != nil {
			t.Error(err)
			return
		}
		if len(relations) == 0 {
			continue
		}

		for index, _ := range relations {
			//if typeId, ok := goodsTypeHash[relations[index].GoodsTypeId]; ok {
			//	relations[index].ChargeType = typeId
			//} else {
			//	relations[index].ChargeType = 99
			//}
			relations[index].IsFreeLabour = 1
			relations[index].PriceType = 1
		}

		b, err := json2.Marshal(relations)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(b))

		db.Model(&OldMaterialExchangeRecord{}).Where("id = ?", r.Id).Update("data", string(b))
	}
	t.Log("ok")
}

func TestRecordPrice(t *testing.T) {
	var list []*OldMaterialExchangeRecord
	if err := db.Find(&list).Error; err != nil {
		t.Error(err)
		return
	}
	for _, r := range list {
		var relations []*OldMaterialExchangeRelation
		if err := json2.Unmarshal([]byte(r.Data), &relations); err != nil {
			t.Error(err)
			return
		}
		if len(relations) == 0 {
			continue
		}

		for index, _ := range relations {
			relations[index].IsFreeLabour = 1
		}
		b, err := json2.Marshal(relations)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(b))

		db.Model(&OldMaterialExchangeRecord{}).Where("id = ?", r.Id).Update("data", string(b))
	}
	t.Log("ok")
}
