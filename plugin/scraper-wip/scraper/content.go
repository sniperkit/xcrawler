/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-04 20:11:56
 * Last Modified: 2017-11-11 11:11:50
 * Modified By: Gaston Siffert
 */

package scraper

/*
import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
)

const (
	dateFormat     = "2 January 2006"
	durationFormat = ""
	hour           = 60
)

var (
	location, _               = time.LoadLocation("Europe/Paris")
	language    monday.Locale = monday.LocaleFrFR
)

// Movie hold all the information related to a movie
type Content struct {
	// Title of the content/page
	Title string `json:"title"`

	// Fetch duration in minutes
	Duration int `json:"duration"`

	// Synopsis summarize the content
	Summary string `json:"synopsis"`

	// Released is the date when the movie was released in France
	Released time.Time `json:"released"`

	// Categories is an array of content kind
	Categories []string `json:"categories"`

	// Image link to an remotly hosted image
	Image string `json:"image"`
}

// MUST be called with a selection targeting a movie division,
// from the summary page.
// Page: http://www.allocine.fr/films/
// Division: "div.card.card-entity.card-entity-list.cf"
//
// TODO: it's an horrible way to extract the data,
// it must exist a way to simplify this extraction:
//
// <span class="ACrL2ZACrpbG0vYWdlbmRhL3NlbS0yMDE3LTEwLTExLw== date blue-link">11 octobre 2017</span>
// <span class="spacer">/</span>
// 1h 56min
// <span class="spacer">/</span>
// <span class="ACrL2ZACrpbG1zL2dlbnJlLTEzMDAyLw==">Comédie dramatique</span>
func (c *Content) summaryMetaData(s *goquery.Selection) (time.Time, int, []string) {
	detailedSelection := s.Find("div.meta-body-item.meta-body-info")

	r, _ := regexp.Compile("[0-9]h [0-6][0-9]min")
	releaseDate := time.Time{}
	duration := 0
	categories := []string{}
	detailedSelection.Contents().Each(func(i int, s *goquery.Selection) {
		class, exist := s.Attr("class")
		if exist {
			if strings.Contains(class, "date") {
				release := s.Text()
				// We ignore the parsing time error, it will be time.Time{}
				// otherwise
				releaseDate, _ = monday.ParseInLocation(dateFormat,
					release, location, language)
			} else if class == "spacer" {
			} else {
				categories = append(categories, s.Text())
			}
		} else {
			text := s.Text()
			if r.MatchString(text) {
				durationString := r.FindString(s.Text())
				t, _ := time.Parse("15h 04min", durationString)
				duration = t.Hour()*hour + t.Minute()
			}
		}
	})

	return releaseDate, duration, categories
}

// MUST be called with a selection targeting a movie division,
// from the summary page.
// Page: http://www.allocine.fr/films/
// Division: "div.card.card-entity.card-entity-list.cf"
func (c *Movie) fromSummary(s *goquery.Selection) {
	title := s.Find("h2.meta-title").Find("a").Text()
	// We ignore if the image isn't found, in any case it will be an empty link
	image, _ := s.Find("img.thumbnail-img").Attr("data-src")
	summary := s.Find("div.summary").Text()

	// Retrieve the meta data
	releaseDate, duration, categories := c.summaryMetaData(s)
	// Assign the value
	*c = Content{
		Title:      strings.TrimSpace(title),
		Image:      image,
		Summary:    strings.TrimSpace(summary),
		Released:   releaseDate,
		Categories: categories,
		Duration:   duration,
	}
}

refs:
- https://github.com/Vorian-Atreides/allocine/blob/master/movies/scraper.go

*/
