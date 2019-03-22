// The gentmx tool generates TMX maps from a sequence of dungeon pieces (i.e.
// miniture tiles).
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mewkiz/pkg/osutil"
	"github.com/pkg/errors"
)

func usage() {
	const use = `
Generate TMX maps from a sequence of dungeon pieces (i.e. miniture tiles).

Usage:

	gentmx [OPTION]... FILE.bin

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
		// output specifies the output path.
		output string
	)
	flag.StringVar(&dtype, "dtype", "l1", "dungeon type (town, l1, l2, l3 or l4)")
	flag.StringVar(&mpqDir, "mpqdir", "diabdat", `path to extracted "diabdat.mpq"`)
	flag.StringVar(&output, "o", "", "output path")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	binPath := flag.Arg(0)
	if !osutil.Exists(mpqDir) {
		log.Fatalf("unable to locate %q directory", mpqDir)
	}

	// Create output file if specified by `-o`.
	w := os.Stdout
	if len(output) > 0 {
		f, err := os.Create(output)
		if err != nil {
			log.Fatalf("unable to create %q; %v", output, err)
		}
		defer f.Close()
		w = f
	}

	// Generate TMX map.
	if err := gentmx(w, binPath, dtype, mpqDir); err != nil {
		log.Fatalf("%+v", err)
	}
}

// gentmx generates a TMX map for the specified dungeon type, based on the
// dungeon pieces contained within the given file.
func gentmx(w io.Writer, binPath, dtype, mpqDir string) error {
	// Determine dungeon type specific properties.
	var (
		// Map width in number of cels.
		mapWidth int
		// Map height in number of cels
		mapHeight int
		// Map title.
		title string
		// Name of tileset.
		tileset string
		// Number of tiles per row in tileset.
		ntilesPerRow int
		// Tile height in pixels of each tile within the tileset.
		tileHeight int
	)
	switch dtype {
	case "town":
		mapWidth = 96
		mapHeight = 96
		title = "tristram"
		tileset = "tileset_tristram"
		ntilesPerRow = 64
		tileHeight = 256
	case "l1":
		mapWidth = 112
		mapHeight = 112
		title = "cathedral"
		tileset = "tileset_cathedral_theme_1"
		ntilesPerRow = 32
		tileHeight = 160
	case "l2":
		mapWidth = 112
		mapHeight = 112
		title = "catacombs"
		tileset = "tileset_catacombs_theme_1"
		ntilesPerRow = 32
		tileHeight = 160
	case "l3":
		mapWidth = 112
		mapHeight = 112
		title = "caves"
		tileset = "tileset_caves_theme_1"
		ntilesPerRow = 32
		tileHeight = 160
	case "l4":
		mapWidth = 112
		mapHeight = 112
		title = "hell"
		tileset = "tileset_hell_theme_1"
		ntilesPerRow = 32
		tileHeight = 256
	default:
		panic(fmt.Errorf("support for dungeon type %q not yet implemented", dtype))
	}

	// Parse file containing sequence of dungeon pieces (i.e. miniture tiles).
	bin, err := ioutil.ReadFile(binPath)
	if err != nil {
		return errors.WithStack(err)
	}
	got := len(bin)
	want := 4 * mapWidth * mapHeight
	if got != want {
		return errors.Errorf("mismatch between number of dungeon pieces and dungeon size %dx%d; expected %d, got %d", mapWidth, mapHeight, want, got)
	}

	// Parse SOL file.
	relSolPath := fmt.Sprintf("levels/%sdata/%s.sol", dtype, dtype)
	solPath := filepath.Join(mpqDir, relSolPath)
	sol, err := ioutil.ReadFile(solPath)
	if err != nil {
		return errors.WithStack(err)
	}

	// Number of dungeon pieces contained within <dtype>.MIN
	ndpieces := len(sol)
	// Tile width in pixels of each tile within the tileset.
	const tileWidth = 64
	// Tileset width in pixels.
	tilesetWidth := tileWidth * ntilesPerRow
	// Tileset height in pixels.
	tilesetHeight := tileHeight * int(math.Ceil(float64(ndpieces)/float64(ntilesPerRow)))
	background := make([][]int, mapWidth)
	for i := range background {
		background[i] = make([]int, mapHeight)
	}
	collision := make([][]int, mapWidth)
	for i := range collision {
		collision[i] = make([]int, mapHeight)
	}
	r := bytes.NewReader(bin)
	const firstID = 41
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			var dpieceID int32
			if err := binary.Read(r, binary.LittleEndian, &dpieceID); err != nil {
				return errors.WithStack(err)
			}
			collision[x][y] = solid(sol, dpieceID)
			if dpieceID != 0 {
				background[x][y] = firstID - 1 + int(dpieceID)
			}
		}
	}

	funcMap := map[string]interface{}{
		"title": strings.Title,
	}
	t, err := template.New("tmx").Funcs(funcMap).Parse(tmxData[1:])
	if err != nil {
		return errors.WithStack(err)
	}
	m := map[string]interface{}{
		"MapWidth":      mapWidth,
		"MapHeight":     mapHeight,
		"Title":         title,
		"Tileset":       tileset,
		"FirstID":       firstID,
		"TileHeight":    tileHeight,
		"TilesetWidth":  tilesetWidth,
		"TilesetHeight": tilesetHeight,
		"Background":    background,
		"Collision":     collision,
	}
	if err := t.Execute(w, m); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

const tmxData = `
<?xml version="1.0" encoding="UTF-8"?>
<map version="1.0" orientation="isometric" width="{{ .MapWidth }}" height="{{ .MapHeight }}" tilewidth="64" tileheight="32">
 <properties>
  <property name="music" value="music/{{ .Title }}.ogg"/>
  <property name="tileset" value="tilesetdefs/{{ .Tileset }}.txt"/>
  <property name="title" value="{{ title .Title }}"/>
 </properties>
 <tileset firstgid="1" name="collision" tilewidth="64" tileheight="32">
  <image source="../tiled_collision.png" width="512" height="160"/>
 </tileset>
 <tileset firstgid="{{ .FirstID }}" name="{{ .Title }}" tilewidth="64" tileheight="{{ .TileHeight }}">
  <image source="../../mods/ember/images/tilesets/{{ .Tileset }}.png" width="{{ .TilesetWidth }}" height="{{ .TilesetHeight }}"/>
 </tileset>
 <layer name="background" width="{{ .MapWidth }}" height="{{ .MapWidth }}">
  <data encoding="csv">
{{ range $i, $v := .Background }}
	{{- if ne $i 0 }}
		{{- printf ",\n" }}
	{{- end }}
	{{- range $j, $u := . }}
		{{- if ne $j 0 }}
			{{- printf "," }}
		{{- end }}
		{{- printf "%d" $u }}
	{{- end }}
{{- end }}
  </data>
 </layer>
 <layer name="collision" width="{{ .MapWidth }}" height="{{ .MapWidth }}" visible="0">
  <data encoding="csv">
{{ range $i, $v := .Collision }}
	{{- if ne $i 0 }}
		{{- printf ",\n" }}
	{{- end }}
	{{- range $j, $u := . }}
		{{- if ne $j 0 }}
			{{- printf "," }}
		{{- end }}
		{{- printf "%d" $u }}
	{{- end }}
{{- end }}
  </data>
 </layer>
</map>
`

const (
	BLOCKS_NONE            = 0
	BLOCKS_ALL             = 1 // block all
	BLOCKS_MOVEMENT        = 2 // block movement
	BLOCKS_ALL_HIDDEN      = 3 // block all (not visible on mini map)
	BLOCKS_MOVEMENT_HIDDEN = 4 // block movement (not visible on mini map)
)

func solid(sol []byte, dpieceID int32) int {
	if dpieceID == 0 {
		// TODO: set collision later.
		return BLOCKS_ALL
	}
	if isL1Door(dpieceID) {
		// TODO: Handle doors by replacing their tiles with open doors, and adding
		// interactable objects (with their own collision) to display the doors.

		// Return 0 to skip collision for now.
		return 0
	}
	col := sol[dpieceID-1]
	const (
		solBlockWalk    = 0x01 // block walk
		sol02           = 0x02 // lighting?
		solBlockMissile = 0x04 // block missile
		sol08           = 0x08 // transparency?
		sol10           = 0x10 // sw wall
		sol20           = 0x20 // se wall
		sol40           = 0x40
		sol80           = 0x80 // fit shrine
	)

	switch {
	// prioritize block movement over block all.
	case col&solBlockWalk != 0:
		if col&solBlockMissile != 0 {
			return BLOCKS_ALL
		}
		return BLOCKS_ALL
	case col&solBlockMissile != 0:
		return BLOCKS_MOVEMENT
	//case col&sol02 != 0:
	//	return 2
	//case col&sol08 != 0:
	//	return 4
	//case col&sol10 != 0:
	//	return 5
	//case col&sol20 != 0:
	//	return 6
	//case col&sol40 != 0:
	//	return 7
	//case col&sol80 != 0:
	//	return 8
	default:
		return 0
	}
}

func isL1Door(dpieceID int32) bool {
	if dpieceID == 0 {
		return false
	}
	switch dpieceID {
	case 44, 46, 51, 56, 214, 393, 395, 408:
		return true
	}
	return false
}
