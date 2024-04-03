文档
https://gorm.io/zh_CN/docs/

https://mp.weixin.qq.com/s/Lj7W6ijlTp3el1WxGiQZng

// go get gorm.io/gorm
// go get gorm.io/driver/mysql


//主要是处理配置封装到dialector中（dns切割，库名，表名前缀，用户密码等）
mysql.open() : dialector 		
    =>  new dialector{  new config }


//生成gorm的db对象 后续操作都是操作gorm.db对象方法
gorm.open(dialector,opts)	 : gorm.db		
    1.创建grom.Config配置文件
        1.1 如果dialector有Apply接口则对config进行值的修改，遍历opts对对config进行值的修改
        1.2 如果config中的某些值为nil，则设置默认值
    2.创建一个gorm.DB对象
        2.1 把上面的配置赋值到db.config中，db.clone等于1
        2.2 设置db.callbacks = &callbacks{  //initializeCallbacks方法，主要是用于执行这些语句用这个当前这个db
            processors: map[string]*processor{
                "create": {db: db},
                "query":  {db: db},
                "update": {db: db},
                "delete": {db: db},
                "row":    {db: db},
                "raw":    {db: db},
            },
        }
        2.3 把db传入dialector初始化，生成sql链接对象赋值到ConnPool上，后续执行sql可以用到这个     
            db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)  //创建db的ConnPool池这个是用sql包
    3.创建gorm.Statement对象并赋值到db上
        3.1 db.Statement = &Statement{ DB:db, ConnPool: db.ConnPool}

```
    db =>
        close = 1  =>标识
        Config =》 配置
        Statement =》
            => DB   =》 Statement（同样包含DB,connPool）
            => ConnPool   			sql.db对象
            => table,model,Dest,select,joins,omits,sql
            => ReflectValue
            => Setting   			sync.Map对象
           
```
---

### db 中的方法
db.Session(session)  or  db.WithContext(ctx)
    1.新建一个db出来。config是之前的db.config Merge 传进来的session


db.DB()    
    1.获取connPool优先db.Statement.ConnPool ，否则取 db.ConnPool
    2.通过类型断言判断是 *sql.tx 事务,  GetDBConnector, *sql.DB 当前db 类型 ，返回sql.db
        *sql.tx 通过反射获取 tx.sql (因为不暴露只能反射)
        *sql.DB 直接返回
        GetDBConnector 通过GetDBConn返回


db.getInstance()  根据clone的值返回db
    //默认值1，是最外层的db，一般用于返回一个新的db对象
    //0值代表的是当前sql，由于where/limit/find 都在同一个db上面处理
    //2值事务内需要用到同个ConnPool
    0：返回本身db     //当前sql
    1: 创建一个新的db(tx) , tx.Statement.db = tx , tx.Statement.ConnPool = db 
    2: 创建一个新的db(tx) , tx.Statement.db = tx , tx.Statement.ConnPool = db.Statement.ConnPool


db.Find()	
1.把model保存到db.Statement.Dest
2.通过反射把model中的各个信息转换保存到内部对象中
3.把where转成sql
4.执行sql(遍历p.fns中间件中  执行前，执行，执行后),数据转换到模型中
5.打印到终端
