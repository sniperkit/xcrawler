package ajii

import (
	"time"
)

func BackgroundWorker(conf Config, quit chan bool) bool {
	logger := GetLogger("worker")
	for {
		select {
		case msg := <-quit:
			logger.Info("Worker closed")

			return msg
		default:
			logger.Debug("Polling")
			time.Sleep(3000 * time.Millisecond)
		}
	}
}
