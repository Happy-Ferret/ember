// Tribute to OpenBSD: https://www.openbsd.org/lyrics.html#41

// Open Source-ami converts the original Diablo 1 game assets into the file
// formats used by Ember.
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
// assets into the file formats used by Ember.
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
if [ ! -d "../mods/ember/images/tileset" ]; then
	mkdir -p ../mods/ember/images/tileset
	# Cathedral.
	echo "Generate Cathedral tilesets."
	montage _dump_/_dpieces_/l1/l1_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_cathedral_theme_1.png
	montage _dump_/_dpieces_/l1/l1_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_cathedral_theme_2.png
	montage _dump_/_dpieces_/l1/l1_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_cathedral_theme_3.png
	montage _dump_/_dpieces_/l1/l1_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_cathedral_theme_4.png
	montage _dump_/_dpieces_/l1/l1_5.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_cathedral_theme_5.png
	montage _dump_/_dpieces_/l1/l1palg.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_cathedral_gray.png
	# Catacombs.
	echo "Generate Catacombs tilesets."
	montage _dump_/_dpieces_/l2/l2_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_catacombs_theme_1.png
	montage _dump_/_dpieces_/l2/l2_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_catacombs_theme_2.png
	montage _dump_/_dpieces_/l2/l2_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_catacombs_theme_3.png
	montage _dump_/_dpieces_/l2/l2_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_catacombs_theme_4.png
	montage _dump_/_dpieces_/l2/l2_5.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_catacombs_theme_5.png
	montage _dump_/_dpieces_/l2/l2palg.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_catacombs_gray.png
	# Caves.
	echo "Generate Caves tilesets."
	montage _dump_/_dpieces_/l3/l3_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_theme_1.png
	montage _dump_/_dpieces_/l3/l3_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_theme_2.png
	montage _dump_/_dpieces_/l3/l3_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_theme_3.png
	montage _dump_/_dpieces_/l3/l3_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_theme_4.png
	montage _dump_/_dpieces_/l3/l3_i.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_theme_ice.png
	montage _dump_/_dpieces_/l3/l3palg.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_gray.png
	montage _dump_/_dpieces_/l3/l3pfoul.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_theme_foul_water.png
	montage _dump_/_dpieces_/l3/l3pwater.pal/dpiece_*.png -background none -tile 32x -geometry 64x160 ../mods/ember/images/tileset/tileset_caves_theme_water.png
	# Hell.
	echo "Generate Hell tilesets."
	montage _dump_/_dpieces_/l4/l4_1.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/ember/images/tileset/tileset_hell_theme_1.png
	montage _dump_/_dpieces_/l4/l4_2.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/ember/images/tileset/tileset_hell_theme_2.png
	montage _dump_/_dpieces_/l4/l4_3.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/ember/images/tileset/tileset_hell_theme_3.png
	montage _dump_/_dpieces_/l4/l4_4.pal/dpiece_*.png -background none -tile 32x -geometry 64x256 ../mods/ember/images/tileset/tileset_hell_theme_4.png
	# Tristram.
	echo "Generate Tristram tilesets."
	montage _dump_/_dpieces_/town/ltpalg.pal/dpiece_*.png -background none -tile 64x -geometry 64x256 ../mods/ember/images/tileset/tileset_tristram_gray.png
	montage _dump_/_dpieces_/town/town.pal/dpiece_*.png -background none -tile 64x -geometry 64x256 ../mods/ember/images/tileset/tileset_tristram.png
fi

# Generate tileset definitions.
gentilesetdef -dtype town > ../mods/ember/tileset/tileset_tristram.txt
gentilesetdef -dtype l1 > ../mods/ember/tileset/tileset_cathedral_theme_1.txt
gentilesetdef -dtype l2 > ../mods/ember/tileset/tileset_catacombs_theme_1.txt
gentilesetdef -dtype l3 > ../mods/ember/tileset/tileset_caves_theme_1.txt
gentilesetdef -dtype l4 > ../mods/ember/tileset/tileset_hell_theme_1.txt

# Generate monster graphics.
echo "Generate monster graphics."
if [ ! -d "../mods/ember/images/monster" ]; then
	mkdir -p ../mods/ember/images/monster
	# Spitting Terror
	echo "Generating Spitting Terror graphics."
	montage _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/spitting_terror.png
	# Winged Fiend
	echo "Generating Winged Fiend graphics."
	montage _dump_/monsters/bat/bat{a,d,h,n,w}/*_2/*.png _dump_/monsters/bat/bat{a,d,h,n,w}/*_3/*.png _dump_/monsters/bat/bat{a,d,h,n,w}/*_4/*.png _dump_/monsters/bat/bat{a,d,h,n,w}/*_5/*.png _dump_/monsters/bat/bat{a,d,h,n,w}/*_6/*.png _dump_/monsters/bat/bat{a,d,h,n,w}/*_7/*.png _dump_/monsters/bat/bat{a,d,h,n,w}/*_0/*.png _dump_/monsters/bat/bat{a,d,h,n,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/winged_fiend.png
	# Devil Kin Brute
	echo "Generating Devil Kin Brute graphics."
	montage _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_2/*.png _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_3/*.png _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_4/*.png _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_5/*.png _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_6/*.png _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_7/*.png _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_0/*.png _dump_/monsters/bigfall/fallg{a,d,h,n,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/devil_kin_brute.png
	# Black Knight
	echo "Generating Black Knight graphics."
	montage _dump_/monsters/black/black{a,d,h,n,w}/*_2/*.png _dump_/monsters/black/black{a,d,h,n,w}/*_3/*.png _dump_/monsters/black/black{a,d,h,n,w}/*_4/*.png _dump_/monsters/black/black{a,d,h,n,w}/*_5/*.png _dump_/monsters/black/black{a,d,h,n,w}/*_6/*.png _dump_/monsters/black/black{a,d,h,n,w}/*_7/*.png _dump_/monsters/black/black{a,d,h,n,w}/*_0/*.png _dump_/monsters/black/black{a,d,h,n,w}/*_1/*.png -gravity south -geometry 160x160+0+0 -tile x8 -background none ../mods/ember/images/monster/black_knight.png
	# Dark Mage
	echo "Generating Dark Mage graphics."
	montage _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_2/*.png _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_3/*.png _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_4/*.png _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_5/*.png _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_6/*.png _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_7/*.png _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_0/*.png _dump_/monsters/darkmage/dmage{a,d,h,n,s}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/dark_mage.png
	# Bone Demon
	echo "Generating Bone Demon graphics."
	montage _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/demskel/demskl{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/bone_demon.png
	# Diablo
	echo "Generating Diablo graphics."
	montage _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/diablo/diablo{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/diablo.png
	# Fallen One Spear Wielder
	echo "Generating Fallen One Spear Wielder graphics."
	montage _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/falspear/phall{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/fallen_one_spear_wielder.png
	# Fallen One Sword Wielder
	echo "Generating Fallen One Sword Wielder graphics."
	montage _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/falsword/fall{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/fallen_one_sword_wielder.png
	# Overlord
	echo "Generating Overlord graphics."
	montage _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/fat/fat{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/overlord.png
	# Butcher
	echo "Generating Butcher graphics."
	montage _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_2/*.png _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_3/*.png _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_4/*.png _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_5/*.png _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_6/*.png _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_7/*.png _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_0/*.png _dump_/monsters/fatc/fatc{a,d,h,n,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/butcher.png
	# Fireman
	echo "Generating Fireman graphics."
	montage _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/fireman/firem{a,d,h,n,s,w}/*_1/*.png -gravity south -geometry 128x171+0+0 -tile x8 -background none ../mods/ember/images/monster/fireman.png
	# Gargoyle
	echo "Generating Gargoyle graphics."
	montage _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/gargoyle/gargo{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/gargoyle.png
	# Goat Archer
	echo "Generating Goat Archer graphics."
	montage _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_2/*.png _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_3/*.png _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_4/*.png _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_5/*.png _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_6/*.png _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_7/*.png _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_0/*.png _dump_/monsters/goatbow/goatb{a,d,h,n,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/goat_archer.png
	# Goat Lord
	echo "Generating Goat Lord graphics."
	montage _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_2/*.png _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_3/*.png _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_4/*.png _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_5/*.png _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_6/*.png _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_7/*.png _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_0/*.png _dump_/monsters/goatlord/goatl{a,d,h,n,w}/*_1/*.png -gravity south -geometry 160x160+0+0 -tile x8 -background none ../mods/ember/images/monster/goat_lord.png
	# Goat Mace Wielder
	echo "Generating Goat Mace Wielder graphics."
	montage _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/goatmace/goat{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/goat_mace_wielder.png
	# Golem
	echo "Generating Golem graphics."
	montage _dump_/monsters/golem/golema/*_2/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_2/*.png _dump_/monsters/golem/golema/*_3/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_3/*.png _dump_/monsters/golem/golema/*_4/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_4/*.png _dump_/monsters/golem/golema/*_5/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_5/*.png _dump_/monsters/golem/golema/*_6/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_6/*.png _dump_/monsters/golem/golema/*_7/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_7/*.png _dump_/monsters/golem/golema/*_0/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_0/*.png _dump_/monsters/golem/golema/*_1/*.png _dump_/monsters/golem/golemd/*.png _dump_/monsters/golem/golems/*.png _dump_/monsters/golem/golemw/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/golem.png
	# Mage
	echo "Generating Mage graphics."
	montage _dump_/monsters/mage/mage{a,d,h,n,s}/*_2/*.png _dump_/monsters/mage/mage{a,d,h,n,s}/*_3/*.png _dump_/monsters/mage/mage{a,d,h,n,s}/*_4/*.png _dump_/monsters/mage/mage{a,d,h,n,s}/*_5/*.png _dump_/monsters/mage/mage{a,d,h,n,s}/*_6/*.png _dump_/monsters/mage/mage{a,d,h,n,s}/*_7/*.png _dump_/monsters/mage/mage{a,d,h,n,s}/*_0/*.png _dump_/monsters/mage/mage{a,d,h,n,s}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/mage.png
	# Magma Demon
	echo "Generating Magma Demon graphics."
	montage _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/magma/magma{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/magma_demon.png
	# Balrog
	echo "Generating Balrog graphics."
	montage _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/mega/mega{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/balrog.png
	# Horned Demon
	echo "Generating Horned Demon graphics."
	montage _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/rhino/rhino{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/horned_demon.png
	# Scavenger
	echo "Generating Scavenger graphics."
	montage _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/scav/scav{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/scavenger.png
	# Skeleton Axe Wielder
	echo "Generating Skeleton Axe Wielder graphics."
	montage _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/skelaxe/sklax{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/skeleton_axe_wielder.png
	# Skeleton Archer
	echo "Generating Skeleton Archer graphics."
	montage _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/skelbow/sklbw{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/skeleton_archer.png
	# Skeleton Sword Wielder
	echo "Generating Skeleton Sword Wielder graphics."
	montage _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/skelsd/sklsr{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/skeleton_sword_wielder.png
	# Skeleton King
	echo "Generating Skeleton King graphics."
	montage _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/sking/sking{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/skeleton_king.png
	# Viper
	echo "Generating Viper graphics."
	montage _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/snake/snake{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/viper.png
	# Hidden
	echo "Generating Hidden graphics."
	montage _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/sneak/sneak{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/hidden.png
	# Succubus
	echo "Generating Succubus graphics."
	montage _dump_/monsters/succ/scbs{a,d,h,n,w}/*_2/*.png _dump_/monsters/succ/scbs{a,d,h,n,w}/*_3/*.png _dump_/monsters/succ/scbs{a,d,h,n,w}/*_4/*.png _dump_/monsters/succ/scbs{a,d,h,n,w}/*_5/*.png _dump_/monsters/succ/scbs{a,d,h,n,w}/*_6/*.png _dump_/monsters/succ/scbs{a,d,h,n,w}/*_7/*.png _dump_/monsters/succ/scbs{a,d,h,n,w}/*_0/*.png _dump_/monsters/succ/scbs{a,d,h,n,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/succubus.png
	# Litch Demon
	echo "Generating Litch Demon graphics."
	montage _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/thin/thin{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/litch_demon.png
	# Invisible Lord
	echo "Generating Invisible Lord graphics."
	montage _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_2/*.png _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_3/*.png _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_4/*.png _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_5/*.png _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_6/*.png _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_7/*.png _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_0/*.png _dump_/monsters/tsneak/tsneak{a,d,h,n,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/invisible_lord.png
	# Unraveler
	echo "Generating Unraveler graphics."
	montage _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/unrav/unrav{a,d,h,n,s,w}/*_1/*.png -gravity south -geometry 96x128+0+0 -tile x8 -background none ../mods/ember/images/monster/unraveler.png
	# Zombie
	echo "Generating Zombie graphics."
	montage _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_2/*.png _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_3/*.png _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_4/*.png _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_5/*.png _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_6/*.png _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_7/*.png _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_0/*.png _dump_/monsters/zombie/zombie{a,d,h,n,s,w}/*_1/*.png -geometry +0+0 -tile x8 -background none ../mods/ember/images/monster/zombie.png
fi

# Copy cursor graphics.
if [ ! -d "../mods/ember/images/cursor" ]; then
	mkdir -p ../mods/ember/images/cursor
	cp _dump_/data/inv/objcurs/objcurs_0001.png ../mods/ember/images/cursor/cursor_hand.png
fi

# Convert music from wav to ogg.
echo "Converting music from wav to ogg."
if [ ! -d "../mods/ember/music" ]; then
	mkdir -p ../mods/ember/music
	ffmpeg -loglevel error -y -i diabdat/music/dintro.wav ../mods/ember/music/intro.ogg
	ffmpeg -loglevel error -y -i diabdat/music/dlvla.wav ../mods/ember/music/cathedral.ogg
	ffmpeg -loglevel error -y -i diabdat/music/dlvlb.wav ../mods/ember/music/catacombs.ogg
	ffmpeg -loglevel error -y -i diabdat/music/dlvlc.wav ../mods/ember/music/caves.ogg
	ffmpeg -loglevel error -y -i diabdat/music/dlvld.wav ../mods/ember/music/hell.ogg
	ffmpeg -loglevel error -y -i diabdat/music/dtowne.wav ../mods/ember/music/tristram.ogg
fi
`
