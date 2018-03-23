package newznab

import (
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"github.com/smquartz/errors"
)

func rawEntriesToEntries(c *Client, raw rawEntries) (entries Entries, err error) {
	for _, rawItem := range raw.Channel.Entries {
		entry := new(Entry)
		entry.General.Title = rawItem.Title
		entry.General.Description = rawItem.Description
		entry.Meta.Dates.Published = rawItem.Date.Add(0)
		entry.Meta.Source.APIKey = c.APIKey
		entry.Meta.Source.Endpoint = c.BaseURL

		err = entry.fromRawEntry(rawItem)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing attributes", 1)
		}

		if torrent, ok := entry.File.(*TorrentFile); ok {
			u, err := url.Parse(rawItem.Enclosure.URL)
			if err != nil {
				return nil, errors.Wrapf(err, "error parsing enclosure URL: %v", 1, rawItem.Enclosure.URL)
			}
			torrent.DownloadURL = u
		}

		err = entry.PopulateComments(c)
		if err != nil {
			// return nil, errors.Wrapf(err, "error populating comments", 1)
			log.Println(errors.Wrapf(err, "error populating comments", 1))
		}

		entry.PopulateFile(c)
		/* if err != nil {
			// return nil, errors.Wrapf(err, "error populating File", 1)
			log.Println(errors.Wrapf(err, "error populating File", 1))
		} */

		entries = append(entries, *entry)
	}
	return entries, nil
}

// fromRawEntry accepts a rawEntry and sets the called on Entry's
// fields based on the values of rawEntry
func (e *Entry) fromRawEntry(raw rawEntry) (err error) {
	for _, attr := range raw.Attributes {
		err = e.fromRawAttribute(attr)
		if err != nil {
			return errors.Wrapf(err, "error proceesing attribute", 1)
		}
	}
	return nil
}

// fromRawAttribute accepts a raw XML attribute and sets the corresponding
// field in Entry
func (e *Entry) fromRawAttribute(raw rawAttribute) (err error) {
	switch {
	case strings.Contains("category,genre,info", raw.Name):
		return e.fromRawGeneralAttribute(raw)
	case strings.Contains("guid,comments,grabs,usenetdate", raw.Name):
		return e.fromRawMetaAttribute(raw)
	case strings.Contains("rating,tvtitle,episode,season,rageid,tvdbid,tvairdate,imdb,imdbtitle,imdbyear,imdbscore,coverurl", raw.Name):
		return e.fromRawContentAttribute(raw)
	case strings.Contains("size,seeders,peers,infohash", raw.Name):
		return e.fromRawFileAttribute(raw)
	default:
		// return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
	return nil
}

// fromRawGeneralAttribute accepts a raw XML attribute that corresponds to a
// field in Entry.General, and sets the corresponding field
func (e *Entry) fromRawGeneralAttribute(raw rawAttribute) error {
	switch raw.Name {
	case "category":
		e.General.Categorisation.Category = append(e.General.Categorisation.Category, raw.Value)
	case "genre":
		e.General.Categorisation.Genre = raw.Value
	case "info":
		e.General.Categorisation.Info = raw.Value
	default:
		return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
	return nil
}

// fromRawMetaAttribute accepts a raw XML attribute that corresponds to a
// field in Entry.Meta, and sets the corresponding field
func (e *Entry) fromRawMetaAttribute(raw rawAttribute) (err error) {
	switch raw.Name {
	case "guid":
		e.Meta.ID, err = uuid.FromString(raw.Value)
		if err != nil {
			return errors.Wrapf(err, "error parsing rawEntry's ID: %s", 1, raw.Value)
		}
	case "grabs":
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing number of grabs: %v", 1, raw.Value)
		}
		e.Meta.Grabs = parsedUint
	case "comments":
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing number of comments: %v", 1, raw.Value)
		}
		e.Meta.Comments.Number = parsedUint
	case "usenetdate":
		if parsedUsetnetDate, err := parseDate(raw.Value); err != nil {
			return errors.Wrapf(err, "failed to parse usenet date: %v", 1, raw.Value)
		} else {
			e.Meta.Dates.Usenet = parsedUsetnetDate
		}
	default:
		return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
	return nil
}

// fromRawContentAttribute accepts a raw XML attribute that corresponds to a
// field in Entry.Content, and sets the corresponding field
func (e *Entry) fromRawContentAttribute(raw rawAttribute) error {
	switch {
	case strings.Contains("rating,tvtitle,episode,season,rageid,tvdbid,tvairdate", raw.Name):
		return e.fromRawTVAttribute(raw)
	case strings.Contains("imdb,imdbtitle,imdbyear,imdbscore,coverurl", raw.Name):
		return e.fromRawMovieAttribute(raw)
	default:
		return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
}

// fromRawTVAttribute accepts a raw XML attribute that corresponds to a field
// in the TV implementation of Entry.Content, and sets the corresponding field.
// If Content is not already set, it will be set to TV.  If it is set to
// another implementation, an error will be returned.
func (e *Entry) fromRawTVAttribute(raw rawAttribute) error {
	tv, ok := e.Content.(*TV)
	if !ok && e.Content != nil {
		return errors.Errorf("encountered TV specific attribute but Content implementation is not set to TV")
	} else if !ok {
		e.Content = new(TV)
		tv = e.Content.(*TV)
	}

	switch raw.Name {
	case "tvairdate":
		if parsedAirDate, err := parseDate(raw.Value); err != nil {
			return errors.Errorf("newznab:Client:Search: failed to parse tvairdate: %v", err)
		} else {
			e.Content.SetAired(parsedAirDate)
		}
	case "tvdbid":
		parsedInt, err := strconv.ParseInt(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing TVDB ID: %v", 1, raw.Value)
		}
		tv.TVDBID = parsedInt
	case "rageid":
		parsedInt, err := strconv.ParseInt(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing TVRage ID: %v", 1, raw.Value)
		}
		tv.TVRageID = parsedInt
	case "season":
		raw.Value = strings.Trim(strings.ToUpper(raw.Value), "S")
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing season number: %v", 1, raw.Value)
		}
		tv.Season = uint(parsedUint)
	case "episode":
		raw.Value = strings.Trim(strings.ToUpper(raw.Value), "E")
		if strings.Contains(raw.Value, "/") {
			raw.Value = strings.Split(raw.Value, "/")[1]
		}
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing episode number: %v", 1, raw.Value)
		}
		tv.Episode = uint(parsedUint)
	case "tvtitle":
		tv.CanonicalTitle = raw.Value
	case "rating":
		parsedFloat, err := strconv.ParseFloat(raw.Value, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing rating: %v", 1, raw.Value)
		}
		tv.Rating = parsedFloat
	default:
		return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
	return nil
}

// fromRawMovieAttribute accepts a raw XML attribute that corresponds to a field
// in the Movie implementation of Entry.Content, and sets the corresponding field.
// If Content is not already set, it will be set to Movie.  If it is set to
// another implementation, an error will be returned.
func (e *Entry) fromRawMovieAttribute(raw rawAttribute) error {
	movie, ok := e.Content.(*Movie)
	if !ok && e.Content != nil {
		return errors.Errorf("encountered Movie specific attribute but Content implementation is not set to Movie")
	} else if !ok {
		e.Content = new(Movie)
		movie = e.Content.(*Movie)
	}

	switch raw.Name {
	case "imdb":
		parsedInt, err := strconv.ParseInt(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing IMDB ID: %v", 1, raw.Value)
		}
		movie.IMDBID = parsedInt
	case "imdbtitle":
		movie.IMDBTitle = raw.Value
	case "imdbyear":
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing IMDB year: %v", 1, raw.Value)
		}
		t := time.Date(int(parsedUint), time.January, 1, 0, 0, 0, 0, time.UTC)
		movie.IMDBYear = t
	case "imdbscore":
		parsedFloat, err := strconv.ParseFloat(raw.Value, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing IMDB score: %v", 1, raw.Value)
		}
		movie.IMDBScore = parsedFloat
	case "coverurl":
		u, err := url.Parse(raw.Value)
		if err != nil {
			return errors.Wrapf(err, "error parsing cover URL: %v", 1, raw.Value)
		}
		movie.Cover = u
	default:
		return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
	return nil
}

// fromRawFileAttribute accepts a raw XML attribute that corresponds to a
// field in Entry.File, and sets the corresponding field
func (e *Entry) fromRawFileAttribute(raw rawAttribute) error {
	switch {
	case strings.Contains("size,seeders,peers,infohash", raw.Name):
		return e.fromRawTorrentAttribute(raw)
	/* case strings.Contains("", raw.Name):
	return e.fromRawNZBAttribute(raw) */
	default:
		return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
}

// fromRawTorrentAttribute accepts a raw XML attribute that corresponds to a field
// in the TorrentFile implementation of Entry.File, and sets the corresponding field.
// If File is not already set, it will be set to TorrentFile.  If it is set to
// another implementation, an error will be returned.
func (e *Entry) fromRawTorrentAttribute(raw rawAttribute) error {
	torrent, ok := e.File.(*TorrentFile)
	if !ok && e.File != nil {
		return errors.Errorf("encountered Torrent specific attribute but File implementation is not set to Torrent")
	} else if !ok {
		e.File = new(TorrentFile)
		torrent = e.File.(*TorrentFile)
	}

	switch raw.Name {
	case "size":
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing torrent contents size: %v", 1, raw.Value)
		}
		torrent.ContentsSize = parsedUint
	case "seeders":
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing number of seeders: %v", 1, raw.Value)
		}
		torrent.Seeders = parsedUint
	case "peers":
		parsedUint, err := strconv.ParseUint(raw.Value, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing number of peers: %v", 1, raw.Value)
		}
		torrent.Peers = parsedUint
	case "infohash":
		parsedHex, err := hex.DecodeString(raw.Value)
		if err != nil {
			return errors.Wrapf(err, "error parsing infohash: %v", 1, raw.Value)
		}
		torrent.InfoHash = parsedHex
	default:
		return errors.Errorf("encountered unknown attribute %v: %v", raw.Name, raw.Value)
	}
	return nil
}
