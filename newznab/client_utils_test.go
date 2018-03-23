package newznab

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestBuildURL(t *testing.T) {
	u, _ := url.Parse("https://domain.tld")
	c := &Client{BaseURL: u}
	u2 := c.buildURL(ModePathAPI, url.Values{
		"t": []string{"a", "b", "c"},
	})
	if u2.String() != "https://domain.tld/api?t=a&t=b&t=c" {
		t.Errorf("Build URL produced incorrect URL")
	}

}

func TestGetURLResponseBody(t *testing.T) {
	testURL, err := url.Parse("https://httpbin.org/base64/aGVsbG8gd29ybGQNCg%3D%3D")
	if err != nil {
		t.Fatalf("Could not parse test URL")
	}
	c := &Client{HTTPClient: &http.Client{}}
	data, err := c.getURLResponseBody(testURL)
	if err != nil {
		t.Errorf("getURLResponseBody failed; %v", err.Error())
	}
	if !strings.HasPrefix(string(data), "hello world") {
		t.Errorf("getURLResponseBody failed; wrong response body; expected prefix \"%v\" got \"%v\"", "hello world", string(data))
	}

	// test invalid URL
	testURL, err = url.Parse("http:/,c.hurc,.hurchhrhrhrc,.23090h")
	if err != nil {
		t.Fatalf("Could not parse test URL")
	}
	_, err = c.getURLResponseBody(testURL)
	if err == nil {
		t.Errorf("getURLResponseBody should have errored")
	}
}
