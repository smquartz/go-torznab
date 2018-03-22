package newznab

import (
	"io/ioutil"
	"net/url"

	"github.com/smquartz/errors"
)

// buildURL produces a url.URL that is made up of the base path specified in
// the Client instance, and the path and query parameters specified as arguments
func (c Client) buildURL(path ModePath, values url.Values) url.URL {
	u := c.BaseURL
	u.Path = string(path)
	u.RawQuery = values.Encode()
	return u
}

// getURLResponseBody is a helper function that performs a GET request on a specified URL,
// and returns the response body as a byte slice
func (c *Client) getURLResponseBody(u url.URL) (data []byte, err error) {
	rsp, err := c.HTTPClient.Get(u.String())
	if err != nil {
		return nil, errors.Wrapf(err, "error performing GET request on %v", 1, u.String())
	}

	data, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading response body", 1)
	}

	return data, nil
}
