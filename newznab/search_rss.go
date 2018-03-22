package newznab

import (
	"net/url"
	"strconv"

	"github.com/satori/go.uuid"
	"github.com/smquartz/errors"
)

// SearchRSS performs an arbitrary RSS query against the torznab indexer, and
// parses and returns the newznab entries the RSS API responded with
func (c *Client) SearchRSS(values url.Values) (Entries, error) {
	values.Set("r", c.APIKey)
	values.Set("i", strconv.Itoa(c.APIUserID))
	return c.entriesFromURL(c.buildURL(ModePathRSS, values))
}

// SearchRSSUntilEntryID fetches the RSS feed in chunks until it finds the
// entry with the given ID, then stops and returns the newznab entries fetched
// thus far.  If it reaches maxRequests, it will return what was fetched up
// until that point.
func (c *Client) SearchRSSUntilEntryID(categories []Category, num int, id uuid.UUID, maxRequests int) (entries Entries, err error) {
	count := 0
	for {
		partition, err := c.SearchRSS(url.Values{
			"num":    []string{strconv.Itoa(num)},
			"t":      stringifyCategories(categories),
			"dl":     []string{"1"},
			"offset": []string{strconv.Itoa(num * count)},
		})
		count++
		if err != nil {
			return nil, errors.Wrapf(err, "error getting RSS page %d", 1, count+1)
		}
		for k, entry := range partition {
			if entry.Meta.ID == id {
				return append(entries, partition[:k]...), nil
			}
		}
		entries = append(entries, partition...)
		if maxRequests != 0 && count == maxRequests {
			break
		}
	}
	return entries, nil
}

// SearchRecentEntries returns up to <num> of the most recent newznab entries
func (c *Client) SearchRecentEntries(categories []Category, num int) (Entries, error) {
	return c.SearchRSS(url.Values{
		"num": []string{strconv.Itoa(num)},
		"t":   stringifyCategories(categories),
		"dl":  []string{"1"},
	})
}
