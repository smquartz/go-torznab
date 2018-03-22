package newznab

import (
	"encoding/xml"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/smquartz/errors"
)

// Comments describes comments on an entry, and information relating
// to those comments
type Comments struct {
	// number of comments on the newznab entry
	Number uint64
	// actual comments on the newznab entry
	Comments []Comment
}

// Comment describes an individual comment on an Entry
type Comment struct {
	Title     string
	Content   string
	Published time.Time
}

// rawComments describes the XML response format for multiple comments
type rawComments struct {
	Channel struct {
		Comments []rawComment `xml:"item"`
	} `xml:"channel"`
}

// rawComment describes the XML response format for a single comment
type rawComment struct {
	Title         string `xml:"title"`
	Description   string `xml:"description"`
	PublishedDate string `xml:"pubDate"`
}

// PopulateComments fetches and updates the Comments for the given newznab entry
func (entry *Entry) PopulateComments(c *Client) error {
	data, err := c.getURLResponseBody(c.buildURL(ModePathAPI, url.Values{
		"t":      []string{"comments"},
		"id":     []string{entry.Meta.ID.String()},
		"apikey": []string{c.APIKey},
	}))
	if err != nil {
		return errors.Wrap(err, 1)
	}

	rsp := new(rawComments)
	err = xml.Unmarshal(data, rsp)
	if err != nil {
		return errors.Wrapf(err, "error unmarshalling comments", 1)
	}

	for _, rComment := range rsp.Channel.Comments {
		comment := Comment{
			Title:   rComment.Title,
			Content: rComment.Description,
		}
		if parsedPubDate, err := time.Parse(time.RFC1123Z, rComment.PublishedDate); err != nil {
			log.WithFields(log.Fields{
				"pub_date": rComment.PublishedDate,
				"err":      err,
			}).Error("newznab:Client:PopulateComments: failed to parse date")
		} else {
			comment.Published = parsedPubDate
		}
		entry.Meta.Comments.Comments = append(entry.Meta.Comments.Comments, comment)
	}
	return nil
}
