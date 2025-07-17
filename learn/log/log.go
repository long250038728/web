package log

import (
	"fmt"
	"golang.org/x/exp/slog"
	"os"
	"plugin"
)

// Plus window版本不支持
// 构建跟使用都需要使用相同go版本
// 无法热卸载
// interface{}需要类型断言
func Plus() {
	// go build -buildmode=plugin -o xxx.so main.go
	p, err := plugin.Open("./xxx.so")
	if err != nil {
		panic(err)
	}
	nameSym, err := p.Lookup("Name")
	if err != nil {
		panic(err)
	}
	funcSym, err := p.Lookup("H")
	if err != nil {
		panic(err)
	}

	name := nameSym.(*string)
	helloFunc := funcSym.(func(string) string)
	if err != nil {
		panic(err)
	}
	fmt.Print(name)
	fmt.Println(helloFunc("hello"))
}

// Log 日志
func Log() {
	// sirupsen/logrus
	// logrus.SetFormatter(&logrus.JSONFormatter{}) //全局设置格式
	// logrus.WithFields(logrus.fields{				//这里通过反射，性能可能会较差
	//	   "animal": "walrus",
	//	   "size": 10
	// }).Info("Hello")

	// uber/zap
	// logger,_ := zap.NewProduction()
	// defer logger.Sync()
	// logger.Info("Hello",						//这里使用避免反射（已经确定了类型）
	//		zap.String("key","value"),
	//		zap.Int("size",10)
	// )

	//官方
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			//将info级别的level 字段改为 "severity"
			if a.Key == slog.LevelKey && a.Value.Any().(slog.Level) == slog.LevelInfo {
				a.Key = "severity"
			}
			return a
		},
	}

	//fileWriter := &lumberjack.Logger{
	//	Filename:   "./log/a.log",
	//	MaxSize:    1,    //最大大小1M
	//	MaxBackups: 3,    //最多保留3个备份数量
	//	MaxAge:     7,    //备份最多保留7天
	//	Compress:   true, //是否压缩
	//	LocalTime:  true, //使用本地时间命名备份
	//}
	//handler := slog.NewTextHandler(os.Stdout, opts)
	//handler := slog.NewJSONHandler(fileWriter, opts)
	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)

	//全局 设置及使用
	slog.SetDefault(logger)
	slog.Error("world", slog.StringValue("1"), slog.IntValue(1))

	logger.Info("hello", slog.StringValue("1"), slog.IntValue(1))
	logger.Error("world", slog.StringValue("1"), slog.IntValue(1))

	//派生实例
	userLogger := logger.With(slog.Group("user_info", slog.String("hello", "world"))) //属性
	userLogger.Error("world", slog.StringValue("1"), slog.IntValue(1))
}
