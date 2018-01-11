package scraper

import (
	"time"

	"go.uber.org/zap"
	// "github.com/k0kubun/pp"
	// "github.com/sirupsen/logrus"
	// prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	// logR    = logrus.New()
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
	errInit       error
)

func init() {
	logger, errInit = zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugaredLogger = logger.Sugar()
	// logR.Formatter = new(prefixed.TextFormatter)
	// logR.Level = logrus.DebugLevel
}

func sugaredLoggerTest(url string) {
	sugaredLogger.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugaredLogger.Infof("Failed to fetch URL: %s", url)
}

func fastLoggerTest(url string) {
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}

/*
func logrusTest() {
	logR.WithFields(logrus.Fields{
		"prefix": "main",
		"animal": "walrus",
		"number": 8,
	}).Debug("Started observing beach")

	logR.WithFields(logrus.Fields{
		"prefix":      "sensor",
		"temperature": -4,
	}).Info("Temperature changes")
}
*/
