package newznab

import (
	"net/http"
	"net/url"
	"testing"
)

func TestPopulateComments(t *testing.T) {
	ts := newMockServer()
	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Errorf("Failed to parse mock server URL")
	}
	client := &Client{HTTPClient: &http.Client{}, BaseURL: *u, APIKey: "gibberish"}
	categories := []Category{CategoryTVSD}
	results, err := client.SearchWithTVRage(categories, 2870, 10, 1)
	if err != nil {
		t.Errorf("Failed to search mock indexer")
	}
	if len(results) < 2 {
		t.Errorf("Insufficient number of results returned")
	}
	err = results[1].PopulateComments(client)
	if err != nil {
		t.Errorf("Failed to populate comments for %v", results[1].Meta.ID.String())
	}
}
