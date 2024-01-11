文档
https://gorm.io/zh_CN/docs/

// go get gorm.io/gorm
// go get gorm.io/driver/mysql


//主要是处理配置封装到dialector中（dns切割，库名，表名前缀，用户密码等）
mysql.open() : dialector 		
    =>  new dialector{  new config }


//生成gorm的db对象 后续操作都是操作gorm.db对象方法
gorm.open()	 : gorm.db			
    1.new db						//创建db对象
    2.dialector.initialize(db)		//把db对象扔到dialector初始化配置
    3.db.coonpool = sql.open()      //sql包     新增连接池从sql包中获取
    4.new statemnet					//新增statement,里面有db对象也存放自己，，查询的时候优先读取statement.DB中的connPool

```
    db =>
        close = 1  =>标识
        Config =》 配置
        Statement =》
            => DB   =》 Statement（同样包含DB,connPool）
            => table,model,Dest,select,joins,omits,sql
            => ReflectValue
            => Setting   			sync.Map对象
            => ConnPool   			sql.db对象
```

---

### db 中的方法
db.DB()    
1.优先获取db.Statement.ConnPool ，否则取 db.ConnPool
2.通过类型断言判断是 *sql.tx 事务,  GetDBConnector, *sql.DB 当前db 类型 ，返回sql.db

db.getInstance()   =》 表示一个事务（可以是多个sql，或同个sql，这就是为什么要类型断言）
1.判断db.clone标识(初始化时是1)   创建一个新的db，Statement中的值是之前db的值，clone = 0，返回新db 。
如果是同条sql的话，getInstance第一次产生的clone = 0，之后各种where都是这个新db处理


db.Find()	
1.把model保存到db.Statement.Dest
2.通过反射把model中的各个信息转换保存到内部对象中
3.把where转成sql
4.执行sql(遍历p.fns中间件中  执行前，执行，执行后),数据转换到模型中
5.打印到终端
