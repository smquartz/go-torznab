package newznab

import (
	"encoding/xml"
)

// rawEntry represents an unparsed single newznab item in search results
type rawEntry struct {
	Title    string `xml:"title,omitempty"`
	Link     string `xml:"link,omitempty"`
	Size     int64  `xml:"size,omitempty"`
	Category struct {
		Domain string `xml:"domain,attr"`
		Value  string `xml:",chardata"`
	} `xml:"category,omitempty"`

	GUID struct {
		GUID        string `xml:",chardata"`
		IsPermaLink bool   `xml:"isPermaLink,attr"`
	} `xml:"guid,omitempty"`

	Comments    string `xml:"comments"`
	Description string `xml:"description"`
	Author      string `xml:"author,omitempty"`

	Source struct {
		URL   string `xml:"url,attr"`
		Value string `xml:"url,chardata"`
	} `xml:"source,omitempty"`

	Date xmlTime `xml:"pubDate,omitempty"`

	Enclosure struct {
		URL    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure,omitempty"`

	Attributes []struct {
		XMLName xml.Name
		Name    string `xml:"name,attr"`
		Value   string `xml:"value,attr"`
	} `xml:"attr"`
}

// rawEntries describes responses returned when searching for newznab entries
type rawEntries struct {
	Version   string `xml:"version,attr"`
	ErrorCode int    `xml:"code,attr"`
	ErrorDesc string `xml:"description,attr"`
	Channel   struct {
		Title string `xml:"title"`
		Link  struct {
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"http://www.w3.org/2005/Atom link"`
		Description string `xml:"description"`
		Language    string `xml:"language,omitempty"`
		Webmaster   string `xml:"webmaster,omitempty"`
		Category    string `xml:"category,omitempty"`
		Image       struct {
			URL         string `xml:"url"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description,omitempty"`
			Width       int    `xml:"width,omitempty"`
			Height      int    `xml:"height,omitempty"`
		} `xml:"image"`

		Response struct {
			Offset int `xml:"offset,attr"`
			Total  int `xml:"total,attr"`
		} `xml:"http://www.newznab.com/DTD/2010/feeds/attributes/ response"`

		// All entries that match the search query, up to the response limit.
		Entries []rawEntry `xml:"item"`
	} `xml:"channel"`
}
