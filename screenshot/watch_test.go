package screenshot_test

import (
	"imgSearcher/screenshot"
	"log"
	"os"
	"testing"
)

func TestCrop(t *testing.T) {
	err := os.Chdir("D:/golang/Golang-LocalImageMatcher")
	if err != nil {
		log.Fatal(err)
	}
	screenshot.WatchClipboard()
}
