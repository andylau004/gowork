package main

import (
	"fmt"
)

func dummy_reader() {
	fmt.Println("")
}

type alphaReader struct {
	// alphaReader 里组合了标准库的 io.Reader
	reader io.Reader
}

func newAlphaReader(reader io.Reader) *alphaReader {
	return &alphaReader{reader: reader}
}

func alpha(r byte) byte {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

func TstReaderEntry() {

}
