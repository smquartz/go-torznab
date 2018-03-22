package newznab

import (
	"net/url"
	"strconv"
)

// Search performs an arbitrary API query against the torznab indexer, and
// parses and returns the newznab entries the API responded with
func (c *Client) Search(values url.Values) (Entries, error) {
	values.Set("apikey", c.APIKey)
	return c.entriesFromURL(c.buildURL(ModePathAPI, values))
}

// SearchWithTVRage returns NZBs for the given parameters
func (c *Client) SearchWithTVRage(categories []Category, tvRageID int, season int, episode int) (Entries, error) {
	return c.Search(url.Values{
		"rid":     []string{strconv.Itoa(tvRageID)},
		"cat":     stringifyCategories(categories),
		"season":  []string{strconv.Itoa(season)},
		"episode": []string{strconv.Itoa(episode)},
		"t":       []string{"tvsearch"},
	})
}

// SearchWithTVDB returns NZBs for the given parameters
func (c *Client) SearchWithTVDB(categories []Category, tvDBID int, season int, episode int) (Entries, error) {
	return c.Search(url.Values{
		"tvdbid":  []string{strconv.Itoa(tvDBID)},
		"cat":     stringifyCategories(categories),
		"season":  []string{strconv.Itoa(season)},
		"episode": []string{strconv.Itoa(episode)},
		"t":       []string{"tvsearch"},
	})
}

// SearchWithIMDB returns NZBs for the given parameters
func (c *Client) SearchWithIMDB(categories []Category, imdbID string) (Entries, error) {
	return c.Search(url.Values{
		"imdbid": []string{imdbID},
		"cat":    stringifyCategories(categories),
		"t":      []string{"movie"},
	})
}

// SearchWithQuery returns NZBs for the given parameters
func (c *Client) SearchWithQuery(categories []Category, query string, searchType string) (Entries, error) {
	return c.Search(url.Values{
		"q":   []string{query},
		"cat": stringifyCategories(categories),
		"t":   []string{searchType},
	})
}
