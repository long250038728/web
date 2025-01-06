go get github.com/prometheus/client_golang

counter 计数器
    用于增加或重置的度量值（请求数，错误数量）
    在promQL中 rate函数可以获取一段时间内的度量历史  rate(指标名[时间])  =》 rate(err_count[5m])

Gauge 
    可以上升或下降的数字（当前的cpu值，pod数量等）
    在promQL中 max_over_time,min_over_time,avg_over_time可用于测量指标

histogram直方图
    基于桶值计算的任何值，bucket边界可以由开发人配置（延迟）
    假设我们想观测api请求花费的时间，通过直方图方式存储在桶中，而不是存储每个请求的请求时间。我们可以定义
        小于等于0.3，小于等于0.5，小于等于0.7，小于等于1，小于等于1.2 这几个桶
    histogram_quantile()函数可用于从直方图计算分位数histogram_quantile(0.9,prometheus_http_request_duration_seconds_bucket{handler="/graph"})

summary
    他是直方图的替代，ta更高效，但是会丢失更多数据，ta是在应用程序级别上计算的，因此不能从同一个流程的多个实例中聚合度量，当度量桶事先不知道时，就使用ta



## prometheus
1. Counter/CounterVec 只有增加指标 (Inc加1   Add加N)
2. Gauge/GaugeVec 增加或减少的指标 (有Set赋值    Inc Dec加减1   Add  Sub加减N)
3. Histogram/HistogramVec 统计数据分布的直方图 (定义桶Buckets)
   * prometheus.LinearBuckets(0, 10, 5) =》(初始值，每个桶加N值, 有多少个桶)
   * []float64{0.1, 0.2, 0.5, 1.0, 2.0, 5.0}  => (统计值 ≤ 0.1 的数量。统计值 ≤ 0.2 的数量（包含了 ≤ 0.1 的值）以此类推)
   * Observe(15)方法用于记录值
4. Summary/SummaryVec   提供更精确的百分位统计 (定义Objectives)
   * Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}
     1. Objectives定义了map中3个key 0.5,0.9,0.99这三个key 代表百分之50，百分之90，百分之99有多少值/数量
     2. Objectives定义了map中key对应value  0.05, 0.01, 0.001 即允许的误差是 5% 1% 0.1%
   * Observe(15)方法用于记录值,会根据Observe插入的值进行排序，
     1. 此时0.5代表中位数（位置中间）的值是多少 
     2. 此时0.9代表百分之90的位置的值是多少 
     3. 此时0.9代表百分之99的位置的值是多少


### Vec补充
定义时会指定标签,在使用时需要对标签中进行赋值
```
countVec := prometheus.NewCounterVec(opts2, []string{"title", "name", "time"}) //指定标签
countVec.WithLabelValues("hello", "john", "1s").Add(10) //符合这个标签值进行Add(number) 有该行记录赋值，无则插入一条记录
```