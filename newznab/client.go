package newznab

import (
	"net/http"
	"net/url"
)

// ModePath is a string type that describes the path to append to a base URL
// to use either API or RSS mode
type ModePath string

// API type path constants
const (
	ModePathAPI ModePath = "/api"
	ModePathRSS ModePath = "/rss"
)

// Client is a type for interacting with the newznab API
type Client struct {
	// an optional key to authenticate to the API with
	APIKey string
	// an optional user ID to authenticate to the API with
	APIUserID int
	// the base URL to use for interactions with the API
	BaseURL *url.URL
	// http client to use for interactions with the API
	HTTPClient *http.Client
	// stores capability information retrieved from the API;
	// this describes things like details on what is indexed, supported functions
	// , etc
	capabilities Capabilities
}
