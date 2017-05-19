// The extract_monsters tool extracts monsters assets from the Diablo 1 game.
//
// Note, this tool requires an original copy of diablo.exe. None of the Diablo 1
// game assets are provided by this project.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/sanctuary/exp/d1"
)

func usage() {
	const use = `
Extract monsters assets from the Diablo 1 game.

Usage:

	extract_monster [OPTION]... diablo.exe

Flags:
`
	fmt.Fprint(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	exePath := flag.Arg(0)

	// Extract monster assets from diablo.exe.
	if err := extract(exePath); err != nil {
		log.Fatalf("%+v", err)
	}
}

// extract extracts monster assets from the diablo.exe executable.
func extract(exePath string) error {
	exe, err := d1.ParseFile(exePath)
	if err != nil {
		return errors.WithStack(err)
	}
	pretty.Println("monsters:", exe.Monsters)
	return nil
}
