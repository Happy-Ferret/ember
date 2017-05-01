// The gentilesetdef tool generates tileset definitions based on dungeon type.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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
		// Tileset title.
		title string
	)
	switch dtype {
	case "town":
		tileHeight = 256
		title = "tristram"
	case "l1":
		tileHeight = 160
		title = "cathedral"
	case "l2":
		tileHeight = 160
		title = "catacombs"
	case "l3":
		tileHeight = 160
		title = "caves"
	case "l4":
		tileHeight = 256
		title = "hell"
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

	fmt.Printf("img=images/tilesets/tileset_%s.png\n\n", title)
	x, y := 0, 0
	n := int(math.Ceil(float64(ndpieces) / 16))
	const firstID = 41
	for i := 0; i < ndpieces; i++ {
		id := firstID + i
		fmt.Printf("tile=%d,%d,%d,64,%d,32,%d\n", id, x*64, y*tileHeight, tileHeight, tileHeight-16)
		x++
		if x >= n {
			x = 0
			y++
		}
	}
}
