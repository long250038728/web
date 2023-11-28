package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
	"testing"
)

var indexName = "hello_word"
var addr = "http://159.75.1.200:9220"
var user = "elastic"
var password = "zhubaoe2023Es"

func TestNewEsPersistence(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}
	t.Log(persistence)
	t.Log(err)
}

func TestAllIndex(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}
	names, err := persistence.IndexNames()
	if err != nil {
		return
	}
	for _, name := range names {
		if strings.Contains(name, "zby_") || strings.Contains(name, "report") {
			t.Log(name)
		}
	}

	t.Log(err)
}

func TestIndexInfo(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}
	i, err := persistence.IndexGet("sale_order_record_report").Do(context.Background()) //获取index信息
	t.Log(err)
	t.Log(i)
}

func TestIndexCreateIndex(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	ctx := context.Background()
	do, err := persistence.CreateIndex(indexName).Do(ctx)
	if err != nil {
		t.Log(err)
		return
	}
	if !do.Acknowledged {
		t.Log("Failed to create index")
		return
	}

	// 定义字段映射
	//Text（文本类型）：用于全文搜索的字段类型，会对文本进行分词处理。可以通过指定分词器来定义如何分词。
	//Keyword（关键字类型）：不会分词的字段类型，适用于精确匹配和聚合操作。
	//Numeric（数值类型）：
	//		Long（长整型）
	//		Integer（整型）
	//		Short（短整型）
	//		Byte（字节型）
	//		Double（双精度浮点型）
	//		Float（单精度浮点型）
	//		Half_float（半精度浮点型）
	//		Scaled_float（可缩放精度浮点型）
	//Date（日期类型）：用于存储日期和时间信息的字段类型。
	//Boolean（布尔类型）：表示真或假值的字段类型。
	//Binary（二进制类型）：用于存储二进制数据的字段类型。
	//Array（数组类型）：可以包含多个值的字段类型。
	//Object（对象类型）：可以包含其他字段的复杂数据类型。
	//Nested（嵌套类型）：用于处理嵌套对象的字段类型，可以进行独立的查询和索引。
	mapping := map[string]interface{}{
		"mappings": map[string]map[string]map[string]interface{}{
			"properties": {
				"name":   {"type": "text"},
				"age":    {"type": "integer"},
				"gender": {"type": "keyword"},
			},
		},
	}
	doNew, err := persistence.PutMapping().Index(indexName).BodyJson(mapping).Do(ctx)
	if err != nil {
		return
	}
	if !doNew.Acknowledged {
		t.Log("Failed to put index")
		return
	}
	fmt.Println("Data put successfully")
}

func TestIndexDelIndex(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	ctx := context.Background()
	do, err := persistence.DeleteIndex(indexName).Do(ctx)
	if err != nil {
		t.Log(err)
		return
	}
	if !do.Acknowledged {
		t.Log("Failed to delete index")
		return
	}
	fmt.Println("Data delete successfully")
}

func TestIndexInsert(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	// 插入文档
	doc1 := map[string]interface{}{
		"name":   "Document 1",
		"gender": "This is the content of document 1",
		"age":    12,
		"num":    11,
	}

	doc2 := map[string]interface{}{
		"name":   "Document 2",
		"gender": "This is the content of document 2",
		"age":    22,
		"num":    1,
	}

	_, err = persistence.Index().Index(indexName).BodyJson(doc1).Do(context.Background())
	if err != nil {
		return
	}

	_, err = persistence.Index().Index(indexName).BodyJson(doc2).Do(context.Background())
	if err != nil {
		return
	}

	fmt.Println("Data insert successfully")
}

func TestIndexBulkInsert(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	// 插入文档
	doc1 := map[string]interface{}{
		"name":   "Document 19",
		"gender": "This is the content of document 1",
		"age":    12,
		"num":    11,
	}

	doc2 := map[string]interface{}{
		"name":   "Document 20",
		"gender": "This is the content of document 2",
		"age":    22,
		"num":    1,
	}

	bulk := persistence.Bulk()
	bulk.Add(
		elastic.NewBulkIndexRequest().Index(indexName).Doc(doc1),
		elastic.NewBulkIndexRequest().Index(indexName).Doc(doc2),
	)
	do, err := bulk.Do(context.Background())

	if err != nil {
		t.Log(err)
		return
	}
	t.Log(do.Items)
	fmt.Println("Data bulk successfully")
}

func TestUpdateDoc(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	ctx := context.Background()
	do, err := persistence.Update().Index(indexName).Id("RljKT4sBRD4pu07fMDim").Doc(map[string]interface{}{"gender": "update"}).Do(ctx)
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Println(do)
	fmt.Println("Data update successfully")
}

func TestIndexSearch(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	////filter 不计算相关性
	////must no_must should   计算相关性
	query := elastic.NewBoolQuery()

	//// text: Term精确查询  match模糊匹配单词  match_phrase模糊匹配短语
	query.Must(
		elastic.NewTermQuery("merchant_id", 168),
		elastic.NewTermQuery("merchant_shop_id", 628),
		elastic.NewRangeQuery("gold_weight").Gte(0).Lte(1000),
		elastic.NewMatchQuery("admin_user_name", "吴亦凡"),
		//elastic.NewMatchPhraseQuery("merchant_shop_name", "店"),
	)
	source, _ := query.Source() //es对应的查询语句
	j, _ := json.Marshal(source)
	t.Log(string(j))

	data, err := persistence.Search("sale_order_record_report").
		Query(query).
		Sort("update_date", true).
		From(40).
		Size(10).
		//FetchSourceContext(elastic.NewFetchSourceContext(true).Include("record_id").Exclude("name")). //显示/不显示 哪些字段
		TrackTotalHits(true). //获取total数量（默认为false，如果数量超过10000则显示10000）
		Do(context.Background())
	if data.Hits.TotalHits.Value <= 0 {
		t.Log("找不到")
		return
	}

	for _, s := range data.Hits.Hits {
		fmt.Println(string(s.Source))
	}
}

func TestIndexSearchAge(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	q := elastic.NewBoolQuery()
	q.Must(
		elastic.NewTermQuery("merchant_id", 168),
	)

	//Field分组： mysql中group by的意思
	agr := elastic.NewTermsAggregation().Field("classify_id").
		SubAggregation("sum", elastic.NewSumAggregation().Field("label_price")).
		SubAggregation("min", elastic.NewMinAggregation().Field("label_price"))

	// 加条件用NewFilterAggregation
	// 分组用NewTermsAggregation
	// 具体函数NewSum/MinAggregation
	agr2 := elastic.NewFilterAggregation().Filter(q).SubAggregation(
		"merchant_id_168", elastic.NewTermsAggregation().Field("classify_id").
			SubAggregation("sum", elastic.NewSumAggregation().Field("label_price")).
			SubAggregation("min", elastic.NewMinAggregation().Field("label_price")),
	)

	data, err := persistence.Search("sale_order_record_report").
		Size(0).              //Aggregation无需返回hits数据
		TrackTotalHits(true). //获取total数量（默认为false，如果数量超过10000则显示10000）
		Aggregation("sum_label_price", elastic.NewSumAggregation().Field("label_price")).
		Aggregation("min_label_price", elastic.NewMinAggregation().Field("label_price")).
		Aggregation("sub_label_price", agr).
		Aggregation("sub_label_price2", agr2).
		Do(context.Background())

	aggregation := string(data.Aggregations["sub_label_price2"])
	fmt.Println(aggregation)

}

func TestIndexSearchMerchantGoodsType(t *testing.T) {
	persistence, err := NewEs(addr, user, password)
	if err != nil {
		return
	}

	terms := []elastic.Query{
		elastic.NewTermQuery("platform", 1),
	}
	q := elastic.NewBoolQuery().Filter(terms...)

	goodsTypeAgr := elastic.NewTermsAggregation().Field("goods_type_id").
		SubAggregation("label_price", elastic.NewSumAggregation().Field("label_price"))

	classifyAgr := elastic.NewTermsAggregation().Field("classify_id").
		SubAggregation("label_price", elastic.NewSumAggregation().Field("label_price"))

	//NewTermsAggregation Field分组： mysql中group by的意思
	agr := elastic.NewTermsAggregation().Field("merchant_id").
		SubAggregation("goods_type_id", goodsTypeAgr).
		SubAggregation("classify_id", classifyAgr)

	data, err := persistence.Search("sale_order_record_report").
		Size(0). //Aggregation无需返回hits数据
		Query(q).
		TrackTotalHits(true). //获取total数量（默认为false，如果数量超过10000则显示10000）
		Aggregation("merchant_id", agr).
		Do(context.Background())

	aggregation := string(data.Aggregations["merchant_id"])
	fmt.Println(aggregation)

}
