package main

import (
	"errors"
	"regexp"
)

func IsHeartbeat(input []byte) bool {
	// SR0001L0001    006969XX    [ID00000000]
	// Note: there might be NULL chars mixed in!
	hbRegex := regexp.MustCompile(`^SR\d{4}L\d{4}\s+\w{8}\s|\x00+\[\w+\]$`)
	match := hbRegex.FindIndex(input)
	return match != nil
}

// ParseSIA retrieves relevant parameters from a SIA encoded message.
// Fields are: sequence, receiver, line, account number, command, zone
func ParseSIA(input []byte) ([]string, error) {
	// 01010053"SIA-DCS"0007R0075L0001[#001465|NRP000*'DECKERS'NM]7C9677F21948CC12|#001465
	siaRegex := regexp.MustCompile(`^\w{8}"SIA-DCS"(\d{4})R(\d{4})L(\d{4})\[#(\d{6})\|\w(\w{2})(\d{3}).*`)
	match := siaRegex.FindSubmatch(input)
	if len(match) > 1 { // remove the first field, which is just the matched string
		match = match[1:]
	}
	// convert output to string - easier to work with
	output := make([]string, len(match))
	for i := 0; i < len(match); i++ {
		output[i] = string(match[i][:])
	}
	if len(output) < 6 {
		err := errors.New("Less than 6 SIA fields found")
		return output, err
	}
	return output, nil
}
