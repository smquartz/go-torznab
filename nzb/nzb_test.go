package nzb

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestNzb(t *testing.T) {
	n, err := FromString(testNzb)
	if err != nil {
		t.Errorf("Failed to parse test XML; %v", err)
	}

	if n.Meta["category"] != "TV > HD" {
		t.Errorf("Wrong category: %s", n.Meta["category"])
	}

	if len(n.Files) != 41 {
		t.Fatalf("Wrong number of files: %d", len(n.Files))
	}

	f := n.Files[0]

	if len(f.Groups) != 1 {
		t.Errorf("Wrong number of groups for file 1: %d", len(f.Groups))
	}

	if len(f.Segments) != 3 {
		t.Errorf("Wrong number of segments for file 1: %d", len(f.Segments))
	}

	var expectedSize uint64 = 3538339983
	if size := n.Size(); size != expectedSize {
		t.Errorf("Wrong NZB size; got %d expected %d", size, expectedSize)
	}

	name, err := f.ApproximatedName()
	if err != nil {
		t.Errorf("f.ApproximatedName() errored; %v", err)
	}
	if expectedName := "Artificial Logic Reciept 809790909009964.nzb"; name != expectedName {
		t.Errorf("Wrong approximated file name for file 1; got %v expected %v", name, expectedName)
	}

	// corrupt f.Subject and try again
	f.Subject = "(*)*)&)rhrchuc,h.rcrh292[2309p02"
	_, err = f.ApproximatedName()
	if err == nil {
		t.Errorf("f.ApproximatedName() should have errored")
	}

	n, err = FromBytes([]byte(testNzb))
	if err != nil {
		t.Errorf("Failed to parse test XML; %v", err)
	}

	if n.Meta["category"] != "TV > HD" {
		t.Errorf("Wrong category: %s", n.Meta["category"])
	}

	if len(n.Files) != 41 {
		t.Fatalf("Wrong number of files: %d", len(n.Files))
	}

	f = n.Files[0]

	if len(f.Groups) != 1 {
		t.Errorf("Wrong number of groups for file 1: %d", len(f.Groups))
	}

	if len(f.Segments) != 3 {
		t.Errorf("Wrong number of segments for file 1: %d", len(f.Segments))
	}

	expectedSize = 3538339983
	if size := n.Size(); size != expectedSize {
		t.Errorf("Wrong NZB size; got %d expected %d", size, expectedSize)
	}

	name, err = f.ApproximatedName()
	if err != nil {
		t.Errorf("f.ApproximatedName() errored; %v", err)
	}
	if expectedName := "Artificial Logic Reciept 809790909009964.nzb"; name != expectedName {
		t.Errorf("Wrong approximated file name for file 1; got %v expected %v", name, expectedName)
	}

	// corrupt f.Subject and try again
	f.Subject = "(*)*)&)rhrchuc,h.rcrh292[2309p02"
	_, err = f.ApproximatedName()
	if err == nil {
		t.Errorf("f.ApproximatedName() should have errored")
	}

	n, err = FromFile("./nzb_test_data.xml")
	if err != nil {
		t.Errorf("Failed to parse test XML; %v", err)
	}

	if n.Meta["category"] != "TV > HD" {
		t.Errorf("Wrong category: %s", n.Meta["category"])
	}

	if len(n.Files) != 41 {
		t.Fatalf("Wrong number of files: %d", len(n.Files))
	}

	f = n.Files[0]

	if len(f.Groups) != 1 {
		t.Errorf("Wrong number of groups for file 1: %d", len(f.Groups))
	}

	if len(f.Segments) != 3 {
		t.Errorf("Wrong number of segments for file 1: %d", len(f.Segments))
	}

	expectedSize = 3538339983
	if size := n.Size(); size != expectedSize {
		t.Errorf("Wrong NZB size; got %d expected %d", size, expectedSize)
	}

	name, err = f.ApproximatedName()
	if err != nil {
		t.Errorf("f.ApproximatedName() errored; %v", err)
	}
	if expectedName := "Artificial Logic Reciept 809790909009964.nzb"; name != expectedName {
		t.Errorf("Wrong approximated file name for file 1; got %v expected %v", name, expectedName)
	}

	// corrupt f.Subject and try again
	f.Subject = "(*)*)&)rhrchuc,h.rcrh292[2309p02"
	_, err = f.ApproximatedName()
	if err == nil {
		t.Errorf("f.ApproximatedName() should have errored")
	}

	r := strings.NewReader(testNzb)
	n, err = FromReader(r)
	if err != nil {
		t.Errorf("Failed to parse test XML; %v", err)
	}

	if n.Meta["category"] != "TV > HD" {
		t.Errorf("Wrong category: %s", n.Meta["category"])
	}

	if len(n.Files) != 41 {
		t.Fatalf("Wrong number of files: %d", len(n.Files))
	}

	f = n.Files[0]

	if len(f.Groups) != 1 {
		t.Errorf("Wrong number of groups for file 1: %d", len(f.Groups))
	}

	if len(f.Segments) != 3 {
		t.Errorf("Wrong number of segments for file 1: %d", len(f.Segments))
	}

	expectedSize = 3538339983
	if size := n.Size(); size != expectedSize {
		t.Errorf("Wrong NZB size; got %d expected %d", size, expectedSize)
	}

	name, err = f.ApproximatedName()
	if err != nil {
		t.Errorf("f.ApproximatedName() errored; %v", err)
	}
	if expectedName := "Artificial Logic Reciept 809790909009964.nzb"; name != expectedName {
		t.Errorf("Wrong approximated file name for file 1; got %v expected %v", name, expectedName)
	}

	// corrupt f.Subject and try again
	f.Subject = "(*)*)&)rhrchuc,h.rcrh292[2309p02"
	_, err = f.ApproximatedName()
	if err == nil {
		t.Errorf("f.ApproximatedName() should have errored")
	}
}

func TestMalformedNzb(t *testing.T) {
	r := strings.NewReader(testMalformedNzb)
	dec := xml.NewDecoder(r)
	n := &NZB{}

	err := dec.Decode(n)

	if err == nil {
		t.Errorf("Decode should have returned error")
	} else {
		t.Logf("Decode appropriately returned error: %v", err.Error())
	}
}
