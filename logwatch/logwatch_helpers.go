package logwatch

import (
	"fmt"
	"strings"
)

type GrafanaExploreUrl struct {
	baseurl    string
	datasource string
	queries    []GrafanaExploreQuery
	rangeFrom  int64
	rangeTo    int64
}

type GrafanaExploreQuery struct {
	refId     string
	container string
}

func (g GrafanaExploreUrl) getUrl() string {
	url := g.baseurl

	if strings.HasSuffix(url, "/") && len(url) > 0 {
		url = url[:len(url)-1]
	}

	url += "/explore?panes="
	url += "{\"_an\":{\"datasource\":\"" + g.datasource + "\",\"queries\":["
	for i, query := range g.queries {
		url += "{\"refId\":\"" + query.refId + "\",\"expr\":\"{container=\\\"" + query.container + "\\\"}\",\"queryType\":\"range\",\"datasource\":{\"type\":\"loki\",\"uid\":\"" + g.datasource + "\"},\"editorMode\":\"builder\",\"hide\":false}"
		if i < len(g.queries)-1 {
			url += ","
		}
	}

	url += "],\"range\":{\"from\":\"" + fmt.Sprint(g.rangeFrom) + "\",\"to\":\"" + fmt.Sprint(g.rangeTo) + "\"}}}&schemaVersion=1&orgId=1"

	return url
}
