// Tribute to OpenBSD: https://www.openbsd.org/lyrics.html#41

// Open Source-ami converts the original Diablo 1 game assets into the file
// formats used by Spark.
package main

import (
	"flag"
	"fmt"
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
	w := os.Stdout
	if len(output) > 0 {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0755)
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
	dir, err := goutil.SrcDir("github.com/mewrnd/blizzconv")
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := fmt.Fprintf(w, script[1:], dir, dir, dir, dir); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

const script = `
#!/bin/bash

# Locate extracted diabdat.mpq
if [ ! -f "mpqdump/levels/towndata/town.cel" ]; then
	echo "Unable to locate \"mpqdump\" directory containing the contents of diabdat.mpq"
	echo ""
	echo "   Please extract diabdat.mpq using Ladislav Zezula's MPQ Editor [1]."
	echo ""
	echo "   [1]: http://www.zezula.net/en/mpq/download.html"
	exit 1
fi

# Add symlinks to cel.ini, cl2.ini, dun.ini and mpq.ini.
ln -s %s/images/imgconf/cel.ini
ln -s %s/images/imgconf/cl2.ini
ln -s %s/configs/dunconf/dun.ini
ln -s %s/mpq/mpq.ini

# Convert CEL, CL2, MIN, TIL and DUN files to PNG images.
echo "*.cel"
time img_dump -a -imgini=cel.ini
#echo "*.cl2"
#time img_dump -a -imgini=cl2.ini
echo "*.min"
time min_dump town.min l1.min l2.min l3.min l4.min
#echo "*.til"
#time til_dump town.til l1.til l2.til l3.til l4.til
#echo "*.dun"
#time dun_dump -a

# Generate the tileset for Tristram.
mkdir -p ../mods/spark/images/tilesets
montage _dump_/_pillars_/town/pillar_*.png -background none -tile x16 -geometry 64x256 ../mods/spark/images/tilesets/tileset_town.png
montage _dump_/_pillars_/l1/pillar_*.png -background none -tile x16 -geometry 64x160 ../mods/spark/images/tilesets/tileset_l1.png
montage _dump_/_pillars_/l2/pillar_*.png -background none -tile x16 -geometry 64x160 ../mods/spark/images/tilesets/tileset_l2.png
montage _dump_/_pillars_/l3/pillar_*.png -background none -tile x16 -geometry 64x160 ../mods/spark/images/tilesets/tileset_l3.png
montage _dump_/_pillars_/l4/pillar_*.png -background none -tile x16 -geometry 64x256 ../mods/spark/images/tilesets/tileset_l4.png

# Generate the music for Tristram.
mkdir -p ../mods/spark/music
ffmpeg -i mpqdump/music/dtowne.wav ../mods/spark/music/town.ogg
`
