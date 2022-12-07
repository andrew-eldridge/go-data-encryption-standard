package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

// generate a random bit sequence of length nBits
func randBits(nBits int) string {
	base2 := []rune("01")
	bits := make([]rune, nBits)
	for i := range bits {
		bits[i] = base2[rand.Intn(2)]
	}
	return string(bits)
}

// get random 64-bit key
func getKey() string {
	return randBits(64)
}

// convert slice of bytes to bit string
func bytesToBits(bytes []byte) string {
	bitStr := ""
	for _, b := range bytes {
		currBits := strconv.FormatUint(uint64(b), 2)
		for len(currBits) < 8 {
			currBits = "0" + currBits
		}
		bitStr += currBits
	}
	return bitStr
}

// convert bit string to slice of bytes
func bitsToBytes(bits string) []byte {
	var bytes []byte
	nBytes := len(bits) / 8
	for i := 0; i < nBytes; i++ {
		b, err := strconv.ParseUint(bits[i*8:i*8+8], 2, 8)
		if err != nil {
			fmt.Println(err)
		}
		bytes = append(bytes, byte(b))
	}
	return bytes
}

// apply permutation function to bit string
func permutation(bits string, permutationFunc []int) string {
	newBits := ""
	for _, i := range permutationFunc {
		newBits += string(bits[i-1])
	}
	return newBits
}

// apply circular left bit shift to bit string
func circularLeftShift(bits string, shamt int) string {
	return bits[shamt:] + bits[:shamt]
}

// generate keys for each of the 16 rounds of encryption/decryption
func generateKeys(keyBits string) []string {
	var keys []string

	// apply PC1 key permutation to produce 56-bit key
	keyBits = permutation(keyBits, PC1)

	// key generation loop
	for i := 0; i < 16; i++ {
		// apply bit shift for current round to kL and kR, then recombine
		keyBits = circularLeftShift(keyBits[:28], KeyBitShift[i]) + circularLeftShift(keyBits[28:], KeyBitShift[i])

		// apply PC2 key permutation to produce 48-bit round key
		keys = append(keys, permutation(keyBits, PC2))
	}

	return keys
}

// partition bit string into given number of blocks
func partition(bits string, nBlocks int) []string {
	var bitsPartition []string
	blockSize := len(bits) / nBlocks
	for i := 0; i < nBlocks; i++ {
		bitsPartition = append(bitsPartition, bits[blockSize*i:blockSize*(i+1)])
	}
	return bitsPartition
}

// combine bit partition into a single bit string
func concatPartition(bitsPartition []string) string {
	bits := ""
	for _, bitsBlock := range bitsPartition {
		bits += bitsBlock
	}
	return bits
}

// XOR exclusive or operator for bit strings
func bitsXOR(bits1, bits2 string) string {
	xorBits := ""
	for i := 0; i < len(bits1); i++ {
		if int(bits1[i]) != int(bits2[i]) {
			xorBits += "1"
		} else {
			xorBits += "0"
		}
	}
	return xorBits
}

// binary to decimal conversion
func binToDec(bits string) int {
	dec, err := strconv.ParseInt(bits, 2, 64)
	if err != nil {
		fmt.Println(err)
	}
	return int(dec)
}

// decimal to binary conversion
func decToBin(dec int) string {
	bits := strconv.FormatInt(int64(dec), 2)
	for len(bits) < 4 {
		bits = "0" + bits
	}
	return bits
}

// perform single round of DES encryption
func desRound(cipherBits, roundKey string) string {
	// get cL and cR from round input
	cL := cipherBits[:32]
	cR := cipherBits[32:]
	cRorig := cR

	// apply expansion permutation to cR
	cR = permutation(cR, ExpansionPermutation)

	// XOR between expanded cR and round key
	cR = bitsXOR(cR, roundKey)

	// partition cR and apply S-box substitution
	partitionCR := partition(cR, 8)
	for j := 0; j < len(partitionCR); j++ {
		// get S-box col (bits 1-4) and row (bits 0,5)
		col := binToDec(partitionCR[j][1:5])
		row := binToDec(partitionCR[j][0:1] + partitionCR[j][5:])

		// apply S-box substitution corresponding to col/row on current block
		partitionCR[j] = decToBin(SBox[j][col+row*16])
	}

	// recombine cR partition after applying S-box
	cR = concatPartition(partitionCR)

	// apply P-box permutation to output of S-box
	cR = permutation(cR, PBox)

	// XOR between cL and output of P-box
	cR = bitsXOR(cL, cR)

	// concatenate original cR and new cR for round output
	return cRorig + cR
}

// perform DES encryption
func encrypt(message, keyBits string) string {
	// convert plaintext message to bits
	messageBits := bytesToBits([]byte(message))

	// apply initial permutation
	messageBits = permutation(messageBits, InitialPermutation)

	// generate round keys
	keys := generateKeys(keyBits)

	// init cipher
	cipherBits := messageBits

	// perform 16 iterations of encryption process
	for i := 0; i < 16; i++ {
		cipherBits = desRound(cipherBits, keys[i])
	}

	// final swap of cL and cR
	cipherBits = cipherBits[32:] + cipherBits[:32]

	// apply final permutation to get cipher
	return permutation(cipherBits, FinalPermutation)
}

// perform DES decryption
func decrypt(cipherBits, keyBits string) string {
	// apply initial permutation
	cipherBits = permutation(cipherBits, InitialPermutation)

	// generate round keys
	keys := generateKeys(keyBits)

	// perform 16 iterations of reversed encryption process
	for i := 15; i > -1; i-- {
		cipherBits = desRound(cipherBits, keys[i])
	}

	// final swap of cL and cR
	cipherBits = cipherBits[32:] + cipherBits[:32]

	// apply final permutation to get message bits
	messageBits := permutation(cipherBits, FinalPermutation)

	// decode message bits and return plaintext
	return string(bitsToBytes(messageBits))
}
