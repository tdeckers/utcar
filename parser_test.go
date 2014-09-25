package main

import (
	"bytes"
	"testing"
)

func TestIsHeartbeat(t *testing.T) {
	heartbeat := "SR0001L0001    006969XX    [ID00000000]"
	match := IsHeartbeat([]byte(heartbeat))
	if !match {
		t.Fatal("HB match fail")
	}
	nobeat := "some other text"
	match = IsHeartbeat([]byte(nobeat))
	if match {
		t.Fatal("HB matched while it shouldn't")
	}

}

func TestParseSIA(t *testing.T) {
	sia := "01010053\"SIA-DCS\"0007R0075L0001[#001465|NRP000*'DECKERS'NM]7C9677F21948CC12|#001465"
	match := ParseSIA([]byte(sia))
	if match == nil {
		t.Fatal("SIA match fail")
	}
	if len(match) != 6 {
		t.Fatalf("Didn't find all fields, found (%d)", len(match))
	}
	if !bytes.Equal(match[0], []byte("0007")) ||
		!bytes.Equal(match[1], []byte("0075")) ||
		!bytes.Equal(match[2], []byte("0001")) ||
		!bytes.Equal(match[3], []byte("001465")) ||
		!bytes.Equal(match[4], []byte("RP")) ||
		!bytes.Equal(match[5], []byte("000")) {
		t.Fatalf("Failed to match sequence %v", match)
	}
}
