package content_test

import (
	"testing"

	"github.com/ITLab-Projects/pkg/models/content"
)

func TestFunc_GetUrl(t *testing.T) {
	c := &content.Content{
		HTMLURL: "https://github.com/RTUITLab/ITLab/blob/master/LANDING.md",
	}

	download := c.GetURLForContent("landing/1.png")
	t.Log(download)
	if download != "https://github.com/RTUITLab/ITLab/raw/master/landing/1.png" {
		t.Log("Assert error")
		t.FailNow()
	}
}
