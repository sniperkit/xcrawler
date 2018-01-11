package scraper

import (
	"testing"

	"github.com/d4l3k/go-internetarchive"
)

func TestSnapshot(t *testing.T) {
	resp, err := archive.Snapshot("http://fn.lc")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Resp: %+v", resp)
}
