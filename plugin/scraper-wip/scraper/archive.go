package scraper

import (
	"fmt"

	"github.com/d4l3k/go-internetarchive"
)

func getSnapshotInternetArchive(link string) (string, error) {
	if link == "" {
		link = "http://fn.lc"
	}
	resp, err := archive.Snapshot(link)
	if err != nil {
		return "", err
	}
	fmt.Printf("Resp: %+v\n", resp)
	return fmt.Sprintf("%+v", resp), nil
}
