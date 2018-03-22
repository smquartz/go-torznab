package nzb

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/smquartz/errors"
)

// Bytes returns an NZB entry encoded in its original XML format, as a byte
// slice
func (n NZB) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := xml.NewEncoder(buf)
	if err := encoder.Encode(n); err != nil {
		return nil, errors.Wrapf(err, "error marshalling NZB into XML", 1)
	}
	return buf.Bytes(), nil
}

// String returns an NZB entry encoded in its original XML format, as a string
func (n NZB) String() (string, error) {
	bdata, err := n.Bytes()
	return string(bdata), err
}

// BytesReader returns an io.Reader containing an NZB entry encoded in its
// original XML format
func (n NZB) BytesReader() (io.Reader, error) {
	data, err := n.Bytes()
	return bytes.NewReader(data), err
}
