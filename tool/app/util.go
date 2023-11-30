package app

import (
	"github.com/long250038728/web/tool/cache"
	"github.com/long250038728/web/tool/locker"
	"github.com/long250038728/web/tool/mq"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

type Util struct {
	//db es 里面涉及库内操作，在没有封装之前暴露第三方的库
	Db *gorm.DB
	Es *elastic.Client

	//cache locker mq 主要是一些通用的东西，可以用接口代替
	Cache  cache.Cache
	Locker locker.Locker
	Mq     mq.Mq
}
