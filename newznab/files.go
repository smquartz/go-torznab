package newznab

import (
	"bytes"
	"io"
	"net/url"

	"github.com/smquartz/errors"
	"github.com/smquartz/go-torznab/nzb"
)

// File describes the actual file that an entry corresponds to
type File interface {
	// Size is a function that returns the size of the file contents in bytes
	Size() uint64
	// URL is a function which returns a url where the raw file being described
	// may be downloaded from
	URL() *url.URL
	// Bytes is a function that returns the bytes of the actual file that File
	// describes; e.g. it might return the bytes of a NZB or Torrent file,
	// depending on the implementation
	Bytes() ([]byte, error)
	// BytesReader is a function that returns an io.Reader for the bytes of the
	// actual file that File describes; e.g. it may return an io.Reader for the
	// bytes of a NZB or Torrent file, depending on the implemenattion
	BytesReader() (io.Reader, error)
	// Populate is a function that updates the information contained within File
	// by downloading the raw file, and parsing it
	Populate(c *Client, e *Entry) error
}

// PopulateFile populates File.  If File is already set to a specific
// implementation, that implementation's Populate() method will be called,
// else, File will be set to a new instance of NZBFile, and its Populate()
// method called.
func (e *Entry) PopulateFile(c *Client) error {
	if e.File == nil {
		e.File = new(NZBFile)
	}
	return e.File.Populate(c, e)
}

// NZBFile is a File implementation that describes a NZB
type NZBFile struct {
	nzb.NZB
	DownloadURL *url.URL
}

// URL returns a URL where the raw NZB file may be downloaded from
func (n NZBFile) URL() *url.URL { return n.DownloadURL }

// Bytes returns the bytes of the usual XML representation of the NZB file
func (n NZBFile) Bytes() ([]byte, error) { return n.Bytes() }

// BytesReader returns an io.Reader for the bytes of the usual XML
// representation of the NZB file
func (n NZBFile) BytesReader() (io.Reader, error) { return n.BytesReader() }

// populateDownloadURL populates the DownloadURL field of an NZBFile with the
// appropriate value
func (n *NZBFile) populateDownloadURL(c *Client, e *Entry) error {
	n.DownloadURL = c.EntryDownloadURL(*e)
	return nil
}

// Populate populates the NZBFile
func (n *NZBFile) Populate(c *Client, e *Entry) error {
	err := n.populateDownloadURL(c, e)
	if err != nil {
		return errors.Wrapf(err, "error populating download URL", 1)
	}

	raw, err := c.getURLResponseBody(n.URL())
	if err != nil {
		return errors.Wrapf(err, "error requesting NZB file", 1)
	}

	parsedNZB, err := nzb.FromBytes(raw)
	if err != nil {
		return errors.Wrapf(err, "error parsing XML response", 1)
	}

	n.NZB = *parsedNZB
	return nil
}

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
	// bytes of the raw torrent file
	Raw []byte
	// URL to download torrent file from
	DownloadURL *url.URL
}

// Size returns the size of the torrent contents in bytes
func (t TorrentFile) Size() uint64 { return t.ContentsSize }

// URL returns a URL where the raw torrent file may be downloaded from
func (t TorrentFile) URL() *url.URL { return t.DownloadURL }

// Bytes returns the bytes of the raw torrent file
func (t TorrentFile) Bytes() ([]byte, error) { return t.Raw, nil }

// BytesReader returns an io.Reader for the bytes of the raw torrent file
func (t TorrentFile) BytesReader() (io.Reader, error) { return bytes.NewBuffer(t.Raw), nil }

// Populate populates the TorrentFile, with
// the information contained within the raw torrent file.
func (t *TorrentFile) Populate(c *Client, e *Entry) (err error) {
	if t.URL() == nil {
		return errors.Errorf("Empty download URL")
	}
	t.Raw, err = c.getURLResponseBody(t.URL())
	if err != nil {
		return errors.Wrap(err, 1)
	}
	return nil
}
