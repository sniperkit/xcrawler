/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-05 14:11:14
 * Last Modified: 2017-11-10 22:11:48
 * Modified By: Gaston Siffert
 */
package scraper

/*

//
// Summarize
//

type contentAsyncRequest struct {
	page int
}

type contentAsyncResponse struct {
	data []Content
	err  error
}

// const (
// 	defau
// )

func (c ContentScraper) summarizeWorker(input <-chan contentAsyncRequest,
	output chan<- contentAsyncResponse) {

	for {
		// Receive a task
		request := <-input
		// Handle the close(input)
		if request == (contentAsyncRequest{}) {
			return
		}

		data, err := c.Summarize(request.page)
		output <- contentAsyncResponse{data: data, err: err}
	}
}

refs:
- https://github.com/Vorian-Atreides/allocine/blob/master/movies/scraper.go

*/
