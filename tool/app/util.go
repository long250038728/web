package app

import (
	"github.com/long250038728/web/tool/cache"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

type Util struct {
	//db es 里面涉及库内操作，在没有封装之前暴露第三方的库
	db *gorm.DB
	es *elastic.Client

	//cache locker 主要是一些通用的东西，可以用接口代替
	cache cache.Cache
	//locker
}
