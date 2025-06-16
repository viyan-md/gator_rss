package rss

import (
	"context"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, client *http.Client, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, errors.New("failed creating request")
	}

	req.Header.Add("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("request failed")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed reading response")
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, errors.New("failed to unmarshal xml")
	}

	unescapeString(&feed)

	return &feed, nil
}

func unescapeString(f *RSSFeed) {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)

	for i := range f.Channel.Item {
		f.Channel.Item[i].Title = html.UnescapeString(f.Channel.Item[i].Title)
		f.Channel.Item[i].Description = html.UnescapeString(f.Channel.Item[i].Description)
	}
}
