package main

import (
	"testing"
)

func TestIsHeartbeat(t *testing.T) {
	heartbeat := "SR0001L0001    006969XX    [ID00000000]"
	match := IsHeartbeat(heartbeat)
	if !match {
		t.Fatal("HB match fail")
	}
	nobeat := "some other text"
	match = IsHeartbeat(nobeat)
	if match {
		t.Fatal("HB matched while it shouldn't")
	}

}

func TestParseSIA(t *testing.T) {
	sia := "01010053\"SIA-DCS\"0007R0075L0001[#001465|NRP000*'DECKERS'NM]7C9677F21948CC12|#001465"
	match := ParseSIA(sia)
	if match == nil {
		t.Fatal("SIA match fail")
	}
	if len(match) != 6 {
		t.Fatalf("Didn't find all fields, found (%d)", len(match))
	}
	if match[0] != "0007" || match[1] != "0075" || match[2] != "0001" ||
		match[3] != "001465" || match[4] != "RP" || match[5] != "000" {
		t.Fatalf("Failed to match sequence %v", match)
	}
}
