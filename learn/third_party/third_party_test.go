package third_party

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
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

	handler := slog.NewJSONHandler(os.Stderr, opts)
	logger := slog.New(handler)

	//全局 设置及使用
	slog.SetDefault(logger)
	slog.Error("world", slog.StringValue("1"), slog.IntValue(1))

	logger.Info("hello", slog.StringValue("1"), slog.IntValue(1))
	logger.Error("world", slog.StringValue("1"), slog.IntValue(1))

	//派生实例
	userLogger := logger.With(slog.Group("user_info", slog.String("hello", "world"))) //所有都派生的都会带上with后的信息
	userLogger.Error("world", slog.StringValue("1"), slog.IntValue(1))

	//把logger放入ctx中
	ctx := context.WithValue(context.Background(), slog.Logger{}, userLogger)
	loggerInterface := ctx.Value(slog.Logger{})
	if log, ok := loggerInterface.(*slog.Logger); ok {
		log.Info("ok", "ok")
	}
}
