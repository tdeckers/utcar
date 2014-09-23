package main

import (
	"crypto/des"
	"crypto/rand"
	"log"
)

// function to (un)scramble a key
func Scramble(input []byte) []byte {
	// copy input, otherwise we'll be changing the key
	key := make([]byte, len(input))
	copy(key, input)
	key[3] ^= 0x05
	key[4] ^= 0x23
	key[9] ^= 0x29
	key[1] ^= 0x2D
	key[6] ^= 0x39
	key[20] ^= 0x44
	key[8] ^= 0x45
	key[16] ^= 0x45
	key[5] ^= 0x49
	key[18] ^= 0x50
	key[23] ^= 0x54
	key[0] ^= 0x55
	key[22] ^= 0x69
	key[2] ^= 0x6A
	key[15] ^= 0x88
	key[19] ^= 0x8A
	key[12] ^= 0x94
	key[17] ^= 0xA3
	key[7] ^= 0xA8
	key[21] ^= 0xAA
	key[14] ^= 0xB5
	key[13] ^= 0xC2
	key[10] ^= 0xD3
	key[11] ^= 0xE9
	return key
}

// 3DES ECB descryption.
// We need ECB which is not there in go, so implement here.
// https://code.google.com/p/go/issues/detail?id=5597
func Decrypt3DESECB(input []byte, key []byte) []byte {
	if len(key) != 24 {
		log.Panic("Key must be 24 bytes")
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		log.Panic(err)
	}
	bs := block.BlockSize()
	//log.Printf("Block size: %d", bs)
	if len(input)%bs != 0 {
		log.Panic("Input should be a multiple of blocksize")
	}
	//log.Printf("Input length: %d", len(input))
	m := len(input) / bs
	//log.Printf("Going for %d cycles", m)
	data := []byte{}
	buf := make([]byte, bs)
	for i := 0; i < m; i++ {
		//log.Printf("%d: byte: %s", i, hex.EncodeToString(input[0:bs]))
		block.Decrypt(buf, input)
		input = input[bs:]
		data = append(data, buf...)
	}
	return data
}

func Encrypt3DESECB(input []byte, key []byte) []byte {
	if len(key) != 24 {
		log.Panic("Key must be 24 bytes")
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		log.Panic(err)
	}
	bs := block.BlockSize()
	//log.Printf("Block size: %d", bs)
	if len(input)%bs != 0 {
		log.Panic("Input should be a multiple of blocksize")
	}
	//log.Printf("Input length: %d", len(input))
	m := len(input) / bs
	//log.Printf("Going for %d cycles", m)
	data := []byte{}
	buf := make([]byte, bs)
	for i := 0; i < m; i++ {
		//log.Printf("%d: byte: %s", i, hex.EncodeToString(input[0:bs]))
		block.Encrypt(buf, input)
		input = input[bs:]
		data = append(data, buf...)
	}
	return data
}

// Generate a random key (3DES)
func GenerateKey() ([]byte) {
	key := make([]byte, 24)
	_, err := rand.Read(key)
	if err != nil {
		log.Panic(err)
	}
	return key
}
