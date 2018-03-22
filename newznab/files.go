package newznab

import (
	"net/url"

	"github.com/smquartz/go-torznab/nzb"
)

// File describes the actual file that an entry corresponds to
type File interface {
	// IsFile is a dummy function that does nothing
	IsFile()
	// Size is a function that returns the size of the file contents in bytes
	Size() uint64
}

// NZBFile is a File implementation that describes a NZB
type NZBFile struct {
	nzb.NZB
}

// IsFile is a dummy function that implements the File interface
func (NZBFile) IsFile() {}

// TorrentFile is a File implementation that describes a torrent file
type TorrentFile struct {
	// size of the torrent contents in bytes
	ContentsSize uint64
	// number of seeders on the torrent
	Seeders uint64
	// number of peers on the torrent
	Peers uint64
	// SHA1 hash of the info part of the torrent
	InfoHash []byte `json:"infohash,omitempty"`
	// URL to download torrent file from
	DownloadURL url.URL
}

// IsFile is a dummy function that implements the File interface
func (TorrentFile) IsFile() {}

// Size returns the size of the torrent contents in bytes
func (t TorrentFile) Size() uint64 {
	return t.ContentsSize
}
