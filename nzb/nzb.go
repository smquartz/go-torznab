package nzb

import (
	"encoding/xml"
	"html"

	"github.com/smquartz/errors"
)

// NZB represents a usenet NZB file
// the struct has appropriate XML tags and methods to enable unmarshalling
type NZB struct {
	// the name of the XML tag this represents (nzb)
	XMLName xml.Name `xml:"nzb"`
	// XML metadata associated with the NZB file
	Meta Meta `xml:"head>meta"`
	// NZB file entries contained within this NZB file
	Files []File `xml:"file"`
}

// Meta is a map[string]string that implements UnmarshalXML to
// enable appropriate unmarshalling of XML metadata tags
type Meta map[string]string

// UnmarshalXML implements xml.Unmarshaler for the Meta type
func (m *Meta) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tag := struct {
		Type  string `xml:"type,attr"`
		Value string `xml:",innerxml"`
	}{}

	if err := d.DecodeElement(&tag, &start); err != nil {
		return errors.Wrapf(err, "error decoding start element %v", 1, start.Name)
	}

	if *m == nil {
		*m = make(map[string]string)
	}

	(*m)[tag.Type] = html.UnescapeString(tag.Value)

	return nil
}

// MarshalXML implements xml.Marshaler for the Meta type
func (m Meta) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	for k, v := range m {
		tag := struct {
			Type  string `xml:"type,attr"`
			Value string `xml:",innerxml"`
		}{
			Type:  k,
			Value: v,
		}

		if err := e.EncodeElement(&tag, start); err != nil {
			return errors.Wrapf(err, "error encoding start element %v", 1, start.Name)
		}
	}

	return nil
}

// Size returns the sum of the file sizes within the NZB
func (nzb *NZB) Size() (size uint64) {
	// iterate over all Files in the NZB
	for _, file := range nzb.Files {
		// add the file's size to the total
		size += file.Size()
	}

	return size
}
