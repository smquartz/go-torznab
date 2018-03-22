package nzb

import (
	"strings"
	"testing"
)

func TestFromFileNonexistentFile(t *testing.T) {
	_, err := FromFile("/dont/exist/file.xml")
	if err == nil {
		t.Errorf("FromFile() should have errored")
	}
}

func TestFromReaderEmptyReader(t *testing.T) {
	r := strings.NewReader("")
	_, err := FromReader(r)
	if err == nil {
		t.Errorf("FromReader() should have errored")
	}
}
