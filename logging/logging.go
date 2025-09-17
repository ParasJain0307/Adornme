package logging

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = make(map[string]*zap.SugaredLogger)

var mu sync.Mutex

func InitLogger(module string) *zap.SugaredLogger {
	mu.Lock()
	defer mu.Unlock()

	if l, exists := logger[module]; exists {
		return l
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout", "./logs/" + module + ".log"}
	zapLogger, err := cfg.Build()
	if err != nil {
		zapLogger = zap.NewExample()
	}
	sugar := zapLogger.Sugar()
	logger[module] = sugar
	return sugar

}
