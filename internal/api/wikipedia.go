package api

import (
	"fmt"
	"net/url"

	"github.com/lb-selfbot/selfbot/internal/http"
)

// WikipediaSearch searches Wikipedia for a given query and returns the first result.
func WikipediaSearch(query string) (string, error) {
	searchUrl := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=&titles=%s&format=json",
		url.QueryEscape(query),
	)

	resp, err := http.DefaultClient().Get(searchUrl)
	if err != nil {
		return "No results found", err
	}

	data := resp.Json()

	pages := data.Get("query", "pages")
	for pageId := range pages.MustMap() {
		page := pages.Get(pageId)

		if missing := page.Get("missing").Interface(); missing != nil {
			continue
		}

		return page.Get("extract").MustString(), nil
	}

	return "No results found", fmt.Errorf("no results found")
}
