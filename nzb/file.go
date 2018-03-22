package nzb

import (
	"strings"

	"github.com/smquartz/errors"
)

// File describes the file elements within an NZB.  It has appropriate tags and
// methods to enable deserialisation from XML.
type File struct {
	// string describing who posted the NZB
	Poster string `xml:"poster,attr"`
	// unix formatted date that the NZB was posted
	Date int `xml:"date,attr"`
	// describes what the contents of the NZB file are
	Subject string `xml:"subject,attr"`
	// usenet groups associated with the NZB
	Groups []string `xml:"groups>group,internalxml"`
	// actual file segments to be downloaded
	Segments []Segment `xml:"segments>segment"`
}

// Segment describes a piece of an NZB file, to be downloaded separately.  It
// has appropriate tags and methods to enable deserialisation from XML.
type Segment struct {
	// number of the segment relative to the file
	Number int `xml:"number,attr"`
	// size of the segment in bytes
	Size int `xml:"bytes,attr"`
	// identifier of the segment
	ID string `xml:",innerxml"`
}

// ApproximatedName returns the approximated name of the file that an NZB file
// entry represents.
func (f *File) ApproximatedName() (string, error) {
	parts := strings.Split(f.Subject, `"`)

	n := ""
	if len(parts) > 1 {
		n = strings.Replace(parts[1], "/", "-", -1)
	} else {
		return "", errors.Errorf("could not parse subject")
	}
	return n, nil
}

// Size returns the size in bytes of the file.  It returns the sum of the sizes
// of its segments.
func (f *File) Size() (size uint64) {
	// iterate over all file segments
	for _, segment := range f.Segments {
		// add the segment's size to the file's size
		size += uint64(segment.Size)
	}
	return size
}
