package main

import (
	"encoding/hex"
	"testing"
	"bytes"
)

func TestScramble(t *testing.T) {
	// Test known scramble
	input, err := hex.DecodeString("cb10238cca30e305dd4bc511d474aabf6cea262d1cd008dd")
	if err != nil {
		t.Fatal(err)
	}
	output := Scramble(input)
	outstr := hex.EncodeToString(output)
	const validation = "9e3d4989e979daad986216f840b61f37294976a7587a6189"
	if outstr != validation {
		t.Errorf("Output is %v, expected %v", outstr, validation)
	}

	input, err = hex.DecodeString("dfedcba86c657d87b9f9f54e57f811ab8e9c8ceffba8ef39")
	if err != nil {
		t.Fatal(err)
	}
	output = Scramble(Scramble(input))
	if ! bytes.Equal(input, output) {
		t.Error("Input and output for Scramble-Scramble are not equal")
	}
}

func TestDecrypt3DESECB(t *testing.T) {
	input, err := hex.DecodeString("5590d1ecb21cc8fc2dbcb241fb777a41")
	if err != nil {
		t.Fatal(err)
	}
	key, err := hex.DecodeString("4bcee48bdd7e6b08496328a76c14becec761a8dac4ca078b")
	if err != nil {
		t.Fatal(err)
	}
	const validation = "0a3031303130303542225349412d4443"
	output := Decrypt3DESECB(input, key)
	outstr := hex.EncodeToString(output)
	if outstr != validation {
		t.Errorf("Output is %v, expected %v", outstr, validation)
	}
}

func TestEncrypt3DESECB(t *testing.T) {
	input, err := hex.DecodeString("0a3031303130303542225349412d4443")
	if err != nil {
		t.Fatal(err)
	}
	key, err := hex.DecodeString("4bcee48bdd7e6b08496328a76c14becec761a8dac4ca078b")
	if err != nil {
		t.Fatal(err)
	}
	const validation = "5590d1ecb21cc8fc2dbcb241fb777a41"
	output := Encrypt3DESECB(input, key)
	outstr := hex.EncodeToString(output)
	if outstr != validation {
		t.Errorf("Output is %v, expected %v", outstr, validation)
	}
}
