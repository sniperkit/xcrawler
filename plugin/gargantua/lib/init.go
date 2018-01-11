package lib

import (
	"net/url"
	"sync"
	"time"
)

func StartCrawling(targetURL url.URL, concurrentRequests, timeoutInSeconds int, debugModeIsEnabled bool) error {
	stopTheCrawler := make(chan bool)
	stopTheUI := make(chan bool)
	crawlResult := make(chan error)

	go func() {
		result := crawl(targetURL, CrawlOptions{
			NumberOfConcurrentRequests: int(concurrentRequests),
			Timeout:                    time.Second * time.Duration(timeoutInSeconds),
		}, stopTheCrawler)

		stopTheUI <- true
		crawlResult <- result
	}()

	var uiWaitGroup = &sync.WaitGroup{}
	if debugModeIsEnabled {
		debugf = consoleDebug
	} else {
		debugf = dashboardDebug

		uiWaitGroup.Add(1)
		go func() {
			dashboard(stopTheUI, stopTheCrawler)
			uiWaitGroup.Done()
		}()
	}

	uiWaitGroup.Wait()

	err := <-crawlResult
	if err != nil {
		return err
	}

	return nil
}
