package urlcleaner

import (
	"strings"

	"mvdan.cc/xurls"
)

type UrlCleaner struct{}

func (UrlCleaner) Clean(text string) string {
	toRemove := xurls.Relaxed.FindAllString(text, -1)
	for _, remStr := range toRemove {
		text = strings.Replace(text, remStr, "", -1)
	}
	return text
}

func New() UrlCleaner {
	return UrlCleaner{}
}
