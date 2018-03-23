package newznab

import (
	"net/url"
	"time"
)

// TV is a Content implementation that describes an episode of a TV
// series
type TV struct {
	// air date for the episode according to the entry
	AirDate time.Time
	// ID for the entry for the episode in TheTVDB
	TVDBID int64
	// ID for the entry for the episode in TVRage
	TVRageID int64
	// absolute number of the season
	Season uint
	// number of the episode; entry dependent whether it is absolute or relative
	Episode uint
	// canonical title of the episode
	CanonicalTitle string
	// rating of the episode as recorded in the newznab entry
	Rating float64
}

// IsContent is a dummy function that implements the Content interface
func (TV) IsContent() {}

// Title returns the canonical title of the TV episode
func (t TV) Title() string { return t.CanonicalTitle }

// Aired returns the air date of the episode
func (t TV) Aired() time.Time { return t.AirDate }

// SetAired sets the air date of the episode to the value provided
func (t *TV) SetAired(date time.Time) { t.AirDate = date }

// Movie is a Content implementation that describes a movie
type Movie struct {
	// the air date of the movie according to the newznab entry
	AirDate time.Time
	// ID of the entry for the movie in IMDB
	IMDBID int64
	// title of the movie as recorded by the IMDB entry
	IMDBTitle string
	// year of movie release as recorded by IMDB
	IMDBYear time.Time
	// score of the movie as recorded by IMDB
	IMDBScore float64
	// URL for a cover image for the movie
	Cover *url.URL
}

// IsContent is a dummy function that implements the Content interface
func (Movie) IsContent() {}

// Title returns the IMDB title of the movie
func (m Movie) Title() string { return m.IMDBTitle }

// Aired returns the air date of the movie according to the newznab entry
func (m Movie) Aired() time.Time { return m.AirDate }

// SetAired sets the air date of the movie to the value provided
func (m *Movie) SetAired(date time.Time) { m.AirDate = date }

// Content describes the actual content that an entry corresponds to;
// that is, it describes the movie or episode
type Content interface {
	// IsContent is a dummy function that does nothing
	IsContent()
	// Aired returns the air date of the content
	Aired() time.Time
	// SetAired sets the air date of the content
	SetAired(date time.Time)
	// Title returns the canonical title of the content
	Title() string
}
