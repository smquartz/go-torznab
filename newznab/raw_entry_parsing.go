package newznab

import (
	"encoding/hex"
	"net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"github.com/smquartz/errors"
)

func rawEntriesToEntries(raw rawEntries, baseURL url.URL, key string, userID int) (entries Entries, err error) {
	for _, rawItem := range raw.Channel.Entries {
		entry := new(Entry)
		entry.General.Title = rawItem.Title
		entry.General.Description = rawItem.Description
		entry.Meta.Dates.Published = rawItem.Date.Add(0)
		entry.Meta.Source.APIKey = key
		entry.Meta.Source.Endpoint = baseURL

		err = entry.AttributesFromRawEntry(rawItem)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing attributes", 1)
		}

		entries = append(entries, *entry)
	}
	return entries, nil
}

// AttributesFromRawEntry accepts a rawEntry and sets the called on Entry's
// attributes based on the values of rawEntry
func (e *Entry) AttributesFromRawEntry(raw rawEntry) (err error) {
	for _, attr := range raw.Attributes {
		switch attr.Name {
		case "tvairdate":
			if parsedAirDate, err := parseDate(attr.Value); err != nil {
				return errors.Errorf("newznab:Client:Search: failed to parse tvairdate: %v", err)
			} else {
				e.Content.SetAired(parsedAirDate)
			}
		case "guid":
			e.Meta.ID, err = uuid.FromString(attr.Value)
			if err != nil {
				return errors.Wrapf(err, "error parsing rawEntry's ID: %s", 1, attr.Value)
			}
		/* case "size":
		parsedInt, err := strconv.ParseInt(attr.Value, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "error parsing rawEntry's size: %v", 1, attr.Value)
		}
		// e.File.SetSize(uint64(parsedInt)) */
		case "grabs":
			parsedUint, err := strconv.ParseUint(attr.Value, 0, 32)
			if err != nil {
				return errors.Wrapf(err, "error parsing number of grabs: %v", 1, attr.Value)
			}
			e.Meta.Grabs = parsedUint
		case "comments":
			parsedUint, err := strconv.ParseUint(attr.Value, 0, 32)
			if err != nil {
				return errors.Wrapf(err, "error parsing number of comments: %v", 1, attr.Value)
			}
			e.Meta.Comments.Number = parsedUint
		case "seeders":
			if e.File == nil {
				e.File = new(TorrentFile)
			}
			parsedUint, err := strconv.ParseUint(attr.Value, 0, 32)
			if err != nil {
				return errors.Wrapf(err, "error parsing number of seeders: %v", 1, attr.Value)
			}
			f, ok := e.File.(*TorrentFile)
			if !ok {
				return errors.Errorf("encountered torrent only attribute but entry.File is not of type TorrentFile")
			}
			f.Seeders = parsedUint
		case "peers":
			if e.File == nil {
				e.File = new(TorrentFile)
			}
			parsedUint, err := strconv.ParseUint(attr.Value, 0, 32)
			if err != nil {
				return errors.Wrapf(err, "error parsing number of peers: %v", 1, attr.Value)
			}
			f, ok := e.File.(*TorrentFile)
			if !ok {
				return errors.Errorf("encountered torrent only attribute but entry.File is not of type TorrentFile")
			}
			f.Peers = parsedUint
		case "infohash":
			if e.File == nil {
				e.File = new(TorrentFile)
			}
			parsedHex, err := hex.DecodeString(attr.Value)
			if err != nil {
				return errors.Wrapf(err, "error parsing infohash: %v", 1, attr.Value)
			}
			f, ok := e.File.(*TorrentFile)
			if !ok {
				return errors.Errorf("encountered torrent only attribute but entry.File is not of type TorrentFile")
			}
			f.InfoHash = parsedHex
		case "category":
			e.General.Categorisation.Category = append(e.General.Categorisation.Category, attr.Value)
		case "genre":
			e.General.Categorisation.Genre = attr.Value
		case "tvdbid":
			if e.Content == nil {
				e.Content = new(TV)
			}
			parsedInt, err := strconv.ParseInt(attr.Value, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing TVDB ID: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*TV)
			if !ok {
				return errors.Errorf("encountered TV only attribute but entry.Content is not of type TV")
			}
			c.TVDBID = parsedInt
		case "rageid":
			if e.Content == nil {
				e.Content = new(TV)
			}
			parsedInt, err := strconv.ParseInt(attr.Value, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing TVRage ID: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*TV)
			if !ok {
				return errors.Errorf("encountered TV only attribute but entry.Content is not of type TV")
			}
			c.TVRageID = parsedInt
		case "info":
			e.General.Categorisation.Info = attr.Value
		case "season":
			if e.Content == nil {
				e.Content = new(TV)
			}
			parsedUint, err := strconv.ParseUint(attr.Value, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing season number: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*TV)
			if !ok {
				return errors.Errorf("encountered TV only attribute but entry.Content is not of type TV")
			}
			c.Season = uint(parsedUint)
		case "episode":
			if e.Content == nil {
				e.Content = new(TV)
			}
			parsedUint, err := strconv.ParseUint(attr.Value, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing episode number: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*TV)
			if !ok {
				return errors.Errorf("encountered TV only attribute but entry.Content is not of type TV")
			}
			c.Episode = uint(parsedUint)
		case "tvtitle":
			if e.Content == nil {
				e.Content = new(TV)
			}
			c, ok := e.Content.(*TV)
			if !ok {
				return errors.Errorf("encountered TV only attribute but entry.Content is not of type TV")
			}
			c.CanonicalTitle = attr.Value
		case "rating":
			if e.Content == nil {
				e.Content = new(TV)
			}
			parsedFloat, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing rating: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*TV)
			if !ok {
				return errors.Errorf("encountered TV only attribute but entry.Content is not of type TV")
			}
			c.Rating = parsedFloat
		case "imdb":
			if e.Content == nil {
				e.Content = new(Movie)
			}
			parsedInt, err := strconv.ParseInt(attr.Value, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing IMDB ID: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*Movie)
			if !ok {
				return errors.Errorf("encountered Movie only attribute but entry.Content is not of type Movie")
			}
			c.IMDBID = parsedInt
		case "imdbtitle":
			if e.Content == nil {
				e.Content = new(Movie)
			}
			c, ok := e.Content.(*Movie)
			if !ok {
				return errors.Errorf("encountered Movie only attribute but entry.Content is not of type Movie")
			}
			c.IMDBTitle = attr.Value
		case "imdbyear":
			if e.Content == nil {
				e.Content = new(Movie)
			}
			parsedUint, err := strconv.ParseUint(attr.Value, 0, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing IMDB year: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*Movie)
			if !ok {
				return errors.Errorf("encountered Movie only attribute but entry.Content is not of type Movie")
			}
			t := time.Date(int(parsedUint), time.January, 1, 0, 0, 0, 0, time.UTC)
			c.IMDBYear = t
		case "imdbscore":
			if e.Content == nil {
				e.Content = new(Movie)
			}
			parsedFloat, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return errors.Wrapf(err, "error parsing IMDB score: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*Movie)
			if !ok {
				return errors.Errorf("encountered Movie only attribute but entry.Content is not of type Movie")
			}
			c.IMDBScore = parsedFloat
		case "coverurl":
			if e.Content == nil {
				e.Content = new(Movie)
			}
			u, err := url.Parse(attr.Value)
			if err != nil {
				return errors.Wrapf(err, "error parsing cover URL: %v", 1, attr.Value)
			}
			c, ok := e.Content.(*Movie)
			if !ok {
				return errors.Errorf("encountered Movie only attribute but entry.Content is not of type Movie")
			}
			c.Cover = *u
		case "usenetdate":
			if parsedUsetnetDate, err := parseDate(attr.Value); err != nil {
				return errors.Wrapf(err, "failed to parse usenet date: %v", 1, attr.Value)
			} else {
				e.Meta.Dates.Usenet = parsedUsetnetDate
			}
		default:
			log.WithFields(log.Fields{
				"name":  attr.Name,
				"value": attr.Value,
			}).Debug("encounted unknown attribute")
		}
	}
	return nil
}
