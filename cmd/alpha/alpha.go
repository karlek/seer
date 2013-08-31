package main

import (
	"log"
	"os"

	"github.com/karlek/seer"
)

var (
	// Homework filename.
	filename = os.Getenv("GOPATH") + "/src/github.com/karlek/seer/cmd/alpha/alpha.json"
)

// Error wrapper.
func main() {
	err := alphabet()
	if err != nil {
		log.Fatalln(err)
	}
}

// Learn Greek alphabet.
func alphabet() (err error) {
	h, err := srs.Open(filename)
	if err != nil {
		return err
	}

	err = h.Quiz()
	if err != nil {
		return err
	}
	return nil
}
