package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/karlek/seer"
	"github.com/mewkiz/pkg/goutil"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: alpha PATH.json")
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	var rand bool
	flag.BoolVar(&rand, "r", false, "randomize order")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	jsonRelPath := flag.Arg(0)

	dir, err := goutil.SrcDir("github.com/karlek/seer/cmd/alpha")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	jsonPath := path.Join(dir, jsonRelPath)

	if err := alphabet(jsonPath, rand); err != nil {
		log.Fatalf("%+v", err)
	}
}

// Learn Greek alphabet.
func alphabet(jsonPath string, rand bool) error {
	h, err := seer.Open(jsonPath)
	if err != nil {
		return err
	}
	if rand {
		h.RandomizeOrder()
	}
	if err := h.Quiz(); err != nil {
		return err
	}
	return nil
}
