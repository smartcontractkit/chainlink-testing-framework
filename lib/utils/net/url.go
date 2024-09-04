package net

import "net/url"

func IsValidURL(testURL string) bool {
	parsedURL, err := url.Parse(testURL)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}
