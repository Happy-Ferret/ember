// Tribute to OpenBSD: https://www.openbsd.org/lyrics.html#41

// Open Source-ami converts the original Diablo 1 game assets into the file
// formats used by Spark.
package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/mewkiz/pkg/goutil"
	"github.com/pkg/errors"
)

func main() {
	// Parse command line arguments.
	var (
		// output specifies the output path.
		output string
	)
	flag.StringVar(&output, "o", "", "output path")
	flag.Parse()

	// Create output file if specified by `-o`.
	w := os.Stdout
	if len(output) > 0 {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatalf("unable to create %q; %v", output, err)
		}
		defer f.Close()
		w = f
	}

	// Generate game assets script.
	if err := opensourceami(w); err != nil {
		log.Fatalf("%+v", err)
	}
}

// opensourceami generates a script for converting the original Diablo 1 game
// assets into the file formats used by Spark.
func opensourceami(w io.Writer) error {
	blizzconvDir, err := goutil.SrcDir("github.com/mewrnd/blizzconv")
	if err != nil {
		return errors.WithStack(err)
	}
	t, err := template.New("opensourceami").Parse(script[1:])
	if err != nil {
		return errors.WithStack(err)
	}
	m := map[string]interface{}{
		"BlizzconvDir": blizzconvDir,
	}
	if err := t.Execute(w, m); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

const script = `
#!/bin/bash

# Locate extracted diabdat.mpq
if [ ! -f "diabdat/levels/towndata/town.cel" ]; then
	echo "Unable to locate \"diabdat\" directory containing the contents of diabdat.mpq"
	echo ""
	echo "   Please extract diabdat.mpq to \"_assets_/diabdat/\" using"
	echo "   Ladislav Zezula's MPQ Editor [1]."
	echo ""
	echo "   [1]: http://www.zezula.net/en/mpq/download.html"
	exit 1
fi

# Convert CEL, CL2 and MIN files to PNG images.
echo "Converting CEL, CL2 and MIN files to PNG images."
if [ ! -d "_dump_" ]; then
	mkdir -p _dump_
	time cel_dump -a
	time min_dump -a
fi

# Draw arches onto tileset dungeon pieces.
echo "Draw arches onto tileset dungeon pieces."
fixarches

# Generate tilesets.
echo "Generate tilesets."
if [ ! -d "../mods/spark/images/tilesets" ]; then
	mkdir -p ../mods/spark/images/tilesets
	# Cathedral.
	echo "Generate Cathedral tilesets."
	montage _dump_/_dpieces_/l1/l1_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_cathedral_theme_1.png
	montage _dump_/_dpieces_/l1/l1_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_cathedral_theme_2.png
	montage _dump_/_dpieces_/l1/l1_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_cathedral_theme_3.png
	montage _dump_/_dpieces_/l1/l1_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_cathedral_theme_4.png
	montage _dump_/_dpieces_/l1/l1_5.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_cathedral_theme_5.png
	montage _dump_/_dpieces_/l1/l1palg.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_cathedral_gray.png
	# Catacombs.
	echo "Generate Catacombs tilesets."
	montage _dump_/_dpieces_/l2/l2_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_catacombs_theme_1.png
	montage _dump_/_dpieces_/l2/l2_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_catacombs_theme_2.png
	montage _dump_/_dpieces_/l2/l2_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_catacombs_theme_3.png
	montage _dump_/_dpieces_/l2/l2_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_catacombs_theme_4.png
	montage _dump_/_dpieces_/l2/l2_5.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_catacombs_theme_5.png
	montage _dump_/_dpieces_/l2/l2palg.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_catacombs_gray.png
	# Caves.
	echo "Generate Caves tilesets."
	montage _dump_/_dpieces_/l3/l3_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_theme_1.png
	montage _dump_/_dpieces_/l3/l3_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_theme_2.png
	montage _dump_/_dpieces_/l3/l3_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_theme_3.png
	montage _dump_/_dpieces_/l3/l3_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_theme_4.png
	montage _dump_/_dpieces_/l3/l3_i.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_theme_ice.png
	montage _dump_/_dpieces_/l3/l3palg.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_gray.png
	montage _dump_/_dpieces_/l3/l3pfoul.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_theme_foul_water.png
	montage _dump_/_dpieces_/l3/l3pwater.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/spark/images/tilesets/tileset_caves_theme_water.png
	# Hell.
	echo "Generate Hell tilesets."
	montage _dump_/_dpieces_/l4/l4_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/spark/images/tilesets/tileset_hell_theme_1.png
	montage _dump_/_dpieces_/l4/l4_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/spark/images/tilesets/tileset_hell_theme_2.png
	montage _dump_/_dpieces_/l4/l4_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/spark/images/tilesets/tileset_hell_theme_3.png
	montage _dump_/_dpieces_/l4/l4_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/spark/images/tilesets/tileset_hell_theme_4.png
	# Tristram.
	echo "Generate Tristram tilesets."
	montage _dump_/_dpieces_/town/ltpalg.pal/dpiece_*.png -background none -tile 64x -geometry 64x256 ../mods/spark/images/tilesets/tileset_tristram_gray.png
	montage _dump_/_dpieces_/town/town.pal/dpiece_*.png -background none -tile 64x -geometry 64x256 ../mods/spark/images/tilesets/tileset_tristram.png
fi

# Generate tileset definitions.
gentilesetdef -dtype town > ../mods/spark/tilesetdefs/tileset_tristram.txt
gentilesetdef -dtype l1 > ../mods/spark/tilesetdefs/tileset_cathedral_theme_1.txt
gentilesetdef -dtype l2 > ../mods/spark/tilesetdefs/tileset_catacombs_theme_1.txt
gentilesetdef -dtype l3 > ../mods/spark/tilesetdefs/tileset_caves_theme_1.txt
gentilesetdef -dtype l4 > ../mods/spark/tilesetdefs/tileset_hell_theme_1.txt

# Convert music from wav to ogg.
echo "Converting music from wav to ogg."
if [ ! -d "../mods/spark/music" ]; then
	mkdir -p ../mods/spark/music
	ffmpeg -loglevel quiet -y -i diabdat/music/dintro.wav ../mods/spark/music/intro.ogg
	ffmpeg -loglevel quiet -y -i diabdat/music/dlvla.wav ../mods/spark/music/cathedral.ogg
	ffmpeg -loglevel quiet -y -i diabdat/music/dlvlb.wav ../mods/spark/music/catacombs.ogg
	ffmpeg -loglevel quiet -y -i diabdat/music/dlvlc.wav ../mods/spark/music/caves.ogg
	ffmpeg -loglevel quiet -y -i diabdat/music/dlvld.wav ../mods/spark/music/hell.ogg
	ffmpeg -loglevel quiet -y -i diabdat/music/dtowne.wav ../mods/spark/music/tristram.ogg
fi
`
