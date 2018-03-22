package nzb

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"

	"github.com/smquartz/errors"
)

// FromFile parses a NZB from a file
func FromFile(path string) (*NZB, error) {
	fdata, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening file %v", 1, path)
	}
	return FromReader(fdata)
}

// FromString parses a NZB from a string
func FromString(data string) (*NZB, error) {
	return FromReader(bytes.NewBufferString(data))
}

// FromBytes parses a NZB from a byte slice
func FromBytes(data []byte) (*NZB, error) {
	return FromReader(bytes.NewBuffer(data))
}

// FromReader parses a NZB from an io.Reader
func FromReader(buf io.Reader) (*NZB, error) {
	nzb := new(NZB)
	decoder := xml.NewDecoder(buf)
	err := decoder.Decode(nzb)
	if err != nil {
		return nil, errors.Wrapf(err, "error unmarshalling XML into *NZB", 1)
	}
	return nzb, nil
}
