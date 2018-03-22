package newznab

import (
	"encoding/xml"
	"net/url"

	"github.com/smquartz/errors"
)

// contains functions relating to the processing of newznab responses

// entriesFromURL extracts newznab Entries from the response body returned
// from the given URL.  entriesFromURL performs a GET request against the
// given URL, and parses the response body, ultimately returning Entries.
func (c *Client) entriesFromURL(u url.URL) (entries Entries, err error) {
	rsp, err := c.getURLResponseBody(u)
	if err != nil {
		return nil, errors.Wrap(err, 1)
	}

	feed := new(rawEntries)
	err = xml.Unmarshal(rsp, feed)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshalling XML response into rawEntries", 1)
	}
	if feed.ErrorCode != 0 {
		return nil, errors.Errorf("response body contained error %d: %s", feed.ErrorCode, feed.ErrorDesc)
	}

	entries, err = rawEntriesToEntries(c, *feed)
	if err != nil {
		return nil, errors.Wrapf(err, "error converting rawEntries into Entries", 1)
	}

	return entries, nil
}
