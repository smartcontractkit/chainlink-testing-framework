package parrot

import "net/url"

type Recorder struct {
	URL *url.URL
}

func NewRecorder(url *url.URL) *Recorder {
	return &Recorder{
		URL: url,
	}
}
