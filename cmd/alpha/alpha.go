package main

import (
	"log"
	"path"

	"github.com/karlek/seer"
	"github.com/mewkiz/pkg/goutil"
)

// Homework filename.
var filename string

func init() {
	dir, err := goutil.SrcDir("github.com/karlek/seer/cmd/alpha")
	if err != nil {
		log.Fatalln(err)
	}
	filename = path.Join(dir, "alpha.json")
}

// Error wrapper.
func main() {
	err := alphabet()
	if err != nil {
		log.Fatalln(err)
	}
}

// Learn Greek alphabet.
func alphabet() (err error) {
	h, err := seer.Open(filename)
	if err != nil {
		return err
	}

	err = h.Quiz()
	if err != nil {
		return err
	}
	return nil
}
