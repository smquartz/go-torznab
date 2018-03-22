package newznab

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/smquartz/errors"
)

type xmlTime struct {
	time.Time
}

func (t *xmlTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	e.EncodeToken(xml.CharData([]byte(t.UTC().Format(time.RFC822))))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

func (t *xmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string

	err := d.DecodeElement(&raw, &start)
	if err != nil {
		return err
	}
	date, err := time.Parse(time.RFC1123Z, raw)

	if err != nil {
		return err
	}

	*t = xmlTime{date}
	return nil

}

func parseDate(date string) (time.Time, error) {
	formats := []string{time.RFC3339, time.RFC1123Z}
	var parsedTime time.Time
	var err error
	for _, format := range formats {
		if parsedTime, err = time.Parse(format, date); err == nil {
			return parsedTime, nil
		}
	}
	return parsedTime, errors.Errorf("failed to parse date %s as one of %s", date, strings.Join(formats, ", "))
}
