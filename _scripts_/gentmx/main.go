// The gentmx tool generates TMX maps from a sequence of dungeon pieces (i.e.
// miniture tiles).
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
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
	)
	flag.StringVar(&dtype, "dtype", "l1", "dungeon type (town, l1, l2, l3 or l4)")
	flag.StringVar(&mpqDir, "mpqdir", "diabdat", `path to extracted "diabdat.mpq"`)
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

	// Determine dungeon type specific metrics.
	var (
		// Map width in number of cels.
		mapWidth int
		// Map height in number of cels
		mapHeight int
		// Tile height in pixels of each tile within the tileset.
		tileHeight int
		// Map title.
		title string
	)
	switch dtype {
	case "town":
		mapWidth = 96
		mapHeight = 96
		tileHeight = 256
		title = "tristram"
	case "l1":
		mapWidth = 112
		mapHeight = 112
		tileHeight = 160
		title = "cathedral"
	case "l2":
		mapWidth = 112
		mapHeight = 112
		tileHeight = 160
		title = "catacombs"
	case "l3":
		mapWidth = 112
		mapHeight = 112
		tileHeight = 160
		title = "caves"
	case "l4":
		mapWidth = 112
		mapHeight = 112
		tileHeight = 256
		title = "hell"
	default:
		panic(fmt.Errorf("support for dungeon type %q not yet implemented", dtype))
	}

	// Parse file containing sequence of dungeon pieces (i.e. miniture tiles).
	bin, err := ioutil.ReadFile(binPath)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	got := len(bin)
	want := 4 * mapWidth * mapHeight
	if got != want {
		log.Fatalf("mismatch between number of dungeon pieces and dungeon size %dx%d; expected %d, got %d", mapWidth, mapHeight, want, got)
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
	// Tileset width in pixels.
	tilesetWidth := 64 * int(math.Ceil(float64(ndpieces)/16))
	// Tileset height in pixels.
	tilesetHeight := tileHeight * 16
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
			var v int32
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				log.Fatalf("%+v", errors.WithStack(err))
			}
			collision[x][y] = solid(sol, v)
			if v != 0 {
				background[x][y] = firstID + int(v-1)
			}
		}
	}
	funcMap := map[string]interface{}{
		"title": strings.Title,
	}
	t, err := template.New("tmx").Funcs(funcMap).Parse(tmxData[1:])
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	m := map[string]interface{}{
		"Title":         title,
		"DType":         dtype,
		"FirstID":       firstID,
		"MapWidth":      mapWidth,
		"MapHeight":     mapHeight,
		"TileHeight":    tileHeight,
		"TilesetWidth":  tilesetWidth,
		"TilesetHeight": tilesetHeight,
		"Background":    background,
		"Collision":     collision,
	}
	t.Execute(os.Stdout, m)
}

const tmxData = `
<?xml version="1.0" encoding="UTF-8"?>
<map version="1.0" orientation="isometric" width="{{ .MapWidth }}" height="{{ .MapHeight }}" tilewidth="64" tileheight="32">
 <properties>
  <property name="music" value="music/{{ .Title }}.ogg"/>
  <property name="tileset" value="tilesetdefs/tileset_{{ .Title }}.txt"/>
  <property name="title" value="{{ title .Title }}"/>
 </properties>
 <tileset firstgid="1" name="collision" tilewidth="64" tileheight="32">
  <image source="../tiled_collision.png" width="512" height="160"/>
 </tileset>
 <tileset firstgid="{{ .FirstID }}" name="{{ .Title }}" tilewidth="64" tileheight="{{ .TileHeight }}">
  <image source="../../mods/spark/images/tilesets/tileset_{{ .Title }}.png" width="{{ .TilesetWidth }}" height="{{ .TilesetHeight }}"/>
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
