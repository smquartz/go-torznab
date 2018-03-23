package newznab

import (
	"net/url"
	"time"

	"github.com/satori/go.uuid"
)

// Source describes information relating to the source of an entry
type Source struct {
	// base endpoint of the indexer this entry was retrieved from
	Endpoint *url.URL
	// api key used to access the indexer that this entry was retrieved from
	APIKey string
}

// EntryDates describes published and usenet dates for an Entry
type EntryDates struct {
	// time the Entry was published
	Published time.Time
	// what is this?
	Usenet time.Time
}

// EntryMeta describes information about an Entry itself, rather than the
// content it describes
type EntryMeta struct {
	// entry's GUID
	ID uuid.UUID
	// entry dates
	Dates EntryDates
	// information relating to the source of the entry
	Source Source
	// comments on the Entry
	Comments Comments
	// number of times the newznab entry has been accessed
	Grabs uint64
}

// EntryGeneral describes general information for an Entry
type EntryGeneral struct {
	// title of the newznab entry
	Title string
	// description of the newznab entry
	Description string
	// information relating to the categorisation of the Entry
	Categorisation EntryCategorisation
}

// EntryCategorisation describes information relating to the categorisation of
// the newznab entry, such as genre, category, etc.
type EntryCategorisation struct {
	// newznab categories that the entry belongs to
	Category []string
	// info string for the newznab entry
	Info string
	// genre that the content of the newznab entry belongs to
	Genre string
}

// Entry describes an individual newznab entry found on an index
type Entry struct {
	// information relating to the newznab entry, rather than the
	// content within
	Meta EntryMeta
	// general information from the newznab entry
	General EntryGeneral
	// information relating to the content of the file the newznab entry corresponds
	// to
	Content Content
	// information relating to the file itself that the newznab entry corresponds to
	File File
}

// Entries is simply a []Entry slice
type Entries []Entry
