package newznab

import (
	"net/url"
	"strings"
)

// EntryDownloadURL returns the URL to download the entry from
func (c *Client) EntryDownloadURL(entry Entry) *url.URL {
	return c.buildURL(ModePathAPI, url.Values{
		"t":      []string{"get"},
		"id":     []string{strings.Replace(entry.Meta.ID.String(), "-", "", -1)},
		"apikey": []string{c.APIKey},
	})
}

// DownloadEntry returns the bytes of the actual NZB or other file for the given entry
func (c *Client) DownloadEntry(entry Entry) ([]byte, error) {
	return c.getURLResponseBody(c.EntryDownloadURL(entry))
}
