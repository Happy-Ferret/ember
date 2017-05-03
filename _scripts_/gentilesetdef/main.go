// The gentilesetdef tool generates tileset definitions based on dungeon type.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mewkiz/pkg/osutil"
	"github.com/pkg/errors"
)

func usage() {
	const use = `
Generate tileset definitions based on dungeon type.

Usage:

	gentilesetdef [OPTION]...

Flags:
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line flags.
	var (
		// dtype specifies the dungeon type (town, l1, l2, l3 or l4).
		dtype string
		// mpqDir specifies the path to an extracted "diabdat.mpq".
		mpqDir string
	)
	flag.StringVar(&dtype, "dtype", "l1", "dungeon type (town, l1, l2, l3 or l4)")
	flag.StringVar(&mpqDir, "mpqdir", "diabdat", `path to extracted "diabdat.mpq"`)
	flag.Usage = usage
	flag.Parse()
	if !osutil.Exists(mpqDir) {
		log.Fatalf("unable to locate %q directory", mpqDir)
	}

	// Determine dungeon type specific metrics.
	var (
		// Tile height in pixels of each tile within the tileset.
		tileHeight int
		// Name of tileset.
		tileset string
		// Number of tiles per row in tileset.
		ntilesPerRow int
	)
	switch dtype {
	case "town":
		tileHeight = 256
		tileset = "tileset_tristram"
		ntilesPerRow = 64
	case "l1":
		tileHeight = 160
		tileset = "tileset_cathedral_theme_1"
		ntilesPerRow = 32
	case "l2":
		tileHeight = 160
		tileset = "tileset_catacombs_theme_1"
		ntilesPerRow = 32
	case "l3":
		tileHeight = 160
		tileset = "tileset_caves_theme_1"
		ntilesPerRow = 32
	case "l4":
		tileHeight = 256
		tileset = "tileset_hell_theme_1"
		ntilesPerRow = 32
	default:
		panic(fmt.Errorf("support for dungeon type %q not yet implemented", dtype))
	}

	// Parse SOL file.
	dtypeDataDir := fmt.Sprintf("%sdata", dtype)
	solName := fmt.Sprintf("%s.sol", dtype)
	solPath := filepath.Join(mpqDir, "levels", dtypeDataDir, solName)
	sol, err := ioutil.ReadFile(solPath)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}

	// Number of dungeon pieces contained within <dtype>.MIN
	ndpieces := len(sol)

	fmt.Printf("img=images/tilesets/%s.png\n\n", tileset)
	x, y := 0, 0
	const (
		firstID   = 41
		tileWidth = 64
	)
	for dpieceID := 1; dpieceID <= ndpieces; dpieceID++ {
		id := firstID - 1 + dpieceID
		fmt.Printf("tile=%d,%d,%d,64,%d,32,%d\n", id, x*tileWidth, y*tileHeight, tileHeight, tileHeight-16)
		x++
		if x >= ntilesPerRow {
			x = 0
			y++
		}
	}
}
