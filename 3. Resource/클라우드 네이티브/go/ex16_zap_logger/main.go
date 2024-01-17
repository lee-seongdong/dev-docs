package main

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger zap.Logger // threadSafe
var Sugar zap.SugaredLogger

func sampling() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = "" // 타임스탬프 미출력

	cfg.Sampling = &zap.SamplingConfig{
		Initial:    3,
		Thereafter: 3,
		Hook: func(e zapcore.Entry, sd zapcore.SamplingDecision) {
			if sd == zapcore.LogDropped {
				fmt.Println("event dropped...")
			}
		},
	}

	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func init() {
	sampling()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can`t initialize zap logger: %v", err)
	}

	Logger = *logger        // zap.L()
	Sugar = *logger.Sugar() // zap.S()
}

func logPrintTest() {
	url := "localhost"

	Logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	Sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)

	Sugar.Infof("failed to fetch URL: %s", url)
}

func logSamplingTest() {
	for i := 1; i <= 10; i++ {
		zap.S().Infow(
			"Testing sampling",
			"index", i,
		)
	}
}

func main() {
	logSamplingTest()
}
