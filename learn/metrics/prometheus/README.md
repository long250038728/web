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
