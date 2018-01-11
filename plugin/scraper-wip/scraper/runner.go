package scraper

import (
	"log"
	"time"

	"github.com/matryer/runner"
)

func testRunner() {

	task := runner.Go(func(shouldStop runner.S) error {
		// do setup work
		defer func() {
			// do tear-down work
		}()
		for {

			// do stuff

			// periodically check to see if we should
			// stop or not.
			if shouldStop() {
				break
			}
		}
		return nil // no errors
	})

	task.Stop()
	select {
	case <-task.StopChan():
		// task successfully stopped
	case <-time.After(1 * time.Second):
		// task didn't stop in time
	}

	// execution continues once the code has stopped or has
	// timed out.

	if task.Err() != nil {
		log.Fatalln("task failed:", task.Err())
	}

}
