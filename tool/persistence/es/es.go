package es

import (
	"github.com/long250038728/web/tool/server/http"
	"github.com/olivere/elastic/v7"
)

//github.com/olivere/elastic/v7

type ES struct {
	*elastic.Client
}

func NewEs(config *Config) (*ES, error) {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Address),
		elastic.SetBasicAuth(config.User, config.Password),
		elastic.SetSniff(false),
		elastic.SetHttpClient(http.NewCustomHttpClient(http.Name("ES"))),
	}
	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return &ES{Client: client}, nil
}

//type SearchService struct {
//	*elastic.SearchService
//}
//
//func (c *ES) Search(indices ...string) *SearchService {
//	return &SearchService{SearchService: elastic.NewSearchService(c.Client).Index(indices...)}
//}
//func (s *SearchService) Query(query elastic.Query) *SearchService {
//	s.SearchService.Query(query)
//	return s
//}
//
//func (s *SearchService) From(from int) *SearchService {
//	s.SearchService.From(from)
//	return s
//}
//
//func (s *SearchService) Size(size int) *SearchService {
//	s.SearchService.Size(size)
//	return s
//}
//
//func (s *SearchService) Sort(field string, ascending bool) *SearchService {
//	s.SearchService.Sort(field, ascending)
//	return s
//}
//
//// TrackTotalHits 获取total数量（默认为false，如果数量超过10000则显示10000
//// ES为了避免用户的过大分页请求造成ES服务所在机器内存溢出，默认对深度分页的条数进行了限制，默认的最大条数是10000条
////
////	解决方案 1: index.max_result_window 修改数量
////	        2: track_total_hits  true
//func (s *SearchService) TrackTotalHits(trackTotalHits interface{}) *SearchService {
//	s.SearchService.TrackTotalHits(trackTotalHits)
//	return s
//}
//
//// FetchSourceContext  显示/不显示 哪些字段
//// FetchSourceContext(elastic.NewFetchSourceContext(true).Include("record_id").Exclude("name"))
//func (s *SearchService) FetchSourceContext(fetchSourceContext *elastic.FetchSourceContext) *SearchService {
//	s.SearchService.FetchSourceContext(fetchSourceContext)
//	return s
//}
//
//func (s *SearchService) Do(ctx context.Context) (*elastic.SearchResult, error) {
//	data, err := s.SearchService.Do(ctx)
//	s.addLog(ctx, data, err)
//	return data, err
//}
//
//func (s *SearchService) addLog(ctx context.Context, data *elastic.SearchResult, err error) {
//	span := opentelemetry.NewSpan(ctx, "ES")
//	defer span.Close()
//
//	defer func() {
//		if r := recover(); r != nil {
//			_ = span.AddEvent(map[string]any{
//				"recover": r,
//			})
//		}
//	}()
//
//	// 通过反射获取内部属性searchSource字段。
//	//  field := reflect.ValueOf(s.SearchService).Elem().FieldByName("searchSource")
//	//	searchSource := (*elastic.SearchSource)(field.UnsafePointer())
//	//  info, infoErr := searchSource.Source()
//	info, infoErr := (*elastic.SearchSource)(reflect.Indirect(reflect.ValueOf(s.SearchService)).FieldByName("searchSource").UnsafePointer()).Source()
//	_ = span.AddEvent(map[string]any{
//		"request": info,
//		"error":   infoErr,
//	})
//
//	//_ = span.AddEvent(map[string]any{
//	//	"response": data,
//	//	"error":    err,
//	//})
//}
