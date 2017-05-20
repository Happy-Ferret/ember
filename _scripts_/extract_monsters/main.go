// The extract_monsters tool extracts monsters assets from the Diablo 1 game.
//
// Note, this tool requires an original copy of diablo.exe. None of the Diablo 1
// game assets are provided by this project.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mewkiz/pkg/pathutil"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"github.com/sanctuary/exp/d1"
)

// dbg represents a logger with the "extract_monsters:" prefix, which logs debug
// messages to standard error.
var dbg = log.New(os.Stderr, term.MagentaBold("extract_monsters:")+" ", 0)

func usage() {
	const use = `
Extract monsters assets from the Diablo 1 game.

Usage:

	extract_monsters [OPTION]... diablo.exe

Flags:
`
	fmt.Fprint(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Parse command line arguments.
	var (
		// quiet specifies whether to suppress non-error messages.
		quiet bool
	)
	flag.Usage = usage
	flag.BoolVar(&quiet, "q", false, "suppress non-error messages")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	exePath := flag.Arg(0)
	// Mute debug messages if `-q` is set.
	if quiet {
		dbg.SetOutput(ioutil.Discard)
	}

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
	fmt.Println("#!/bin/bash")
	for _, monster := range exe.Monsters {
		if err := extractMonster(monster); err != nil {
			return errors.WithStack(err)
		}
	}
	//pretty.Println(exe.Monsters)
	return nil
}

// extractMonster extracts the assets of the given monster.
func extractMonster(monster d1.MonsterData) error {
	dbg.Printf("extracting assets of %q.", monster.Name)
	// Extract monster graphics.
	if err := extractMonsterGraphics(monster); err != nil {
		return errors.WithStack(err)
	}
	// Extract monster sounds.
	if err := extractMonsterSounds(monster); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// extractMonsterGraphics extracts the graphics of the given monster.
func extractMonsterGraphics(monster d1.MonsterData) error {
	switch monster.Name {
	case "Wyrm", "Cave Slug", "Devil Wyrm", "Devourer":
		// TODO: Try to locate the Wyrm monster graphics in another MPQ archive
		// than diabdat.mpq.

		// Skip monster; Wyrm monster graphics missing from diabdat.mpq.
		return nil
	}

	actions := []d1.MonsterAction{
		d1.MonsterActionStand,
		d1.MonsterActionWalk,
		d1.MonsterActionAttack,
		d1.MonsterActionHit,
		d1.MonsterActionDie,
	}
	if monster.HasSpecialGraphic {
		actions = append(actions, d1.MonsterActionSpecial)
	}
	// # Spitting Terror
	//
	//    montage \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_2/*.png \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_3/*.png \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_4/*.png \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_5/*.png \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_6/*.png \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_7/*.png \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_0/*.png \
	//       _dump_/monsters/acid/acid{a,d,h,n,s,w}/*_1/*.png \
	//       -geometry +0+0 -tile x8 \
	//       -background none \
	//       ../mods/spark/images/enemies/spitting_terror.png
	script := &bytes.Buffer{}
	fmt.Fprintf(script, "echo 'Extracting graphics for %s'\n", monster.Name)
	script.WriteString("montage \\\n")
	for i := 0; i < 8; i++ {
		direction := (2 + i) % 8
		for _, action := range actions {
			format := strings.ToLower(monster.CL2Path)
			format = strings.Replace(format, `\`, "/", -1)
			relCL2Path := fmt.Sprintf(format, action.Rune())
			relCL2Dir := pathutil.TrimExt(relCL2Path)
			trnDir := ""
			if monster.HasTrn {
				relTrnPath := strings.ToLower(monster.TrnPath)
				relTrnPath = strings.Replace(relTrnPath, `\`, "/", -1)
				dbg.Printf("using colour transition: %q.", relTrnPath)
				trnDir = fmt.Sprintf("%s/", path.Base(relTrnPath))
			}
			switch relCL2Dir {
			case "monsters/darkmage/dmagew":
				// Skip action; darkmage has no walk animation.
				continue
			case "monsters/bigfall/fallgs":
				// TODO: Try to locate the bigfall special action graphics in
				// another MPQ archive than diabdat.mpq.

				// Skip action; bigfall special action graphics missing from
				// diabdat.mpq.
				continue
			case "monsters/golem/golemn", "monsters/golem/golemh":
				// Skip actions; golem has no stand, hit or foo animation.
				continue
			}
			if relCL2Dir == "monsters/golem/golemd" || relCL2Dir == "monsters/golem/golems" {
				// Golem has only one direction for die and special actions.
				relPngPath := fmt.Sprintf("%s/%s*.png", relCL2Dir, trnDir)
				pngPath := filepath.Join("_dump_", relPngPath)
				fmt.Fprintf(script, "\t%s \\\n", pngPath)
				break
			}
			name := path.Base(relCL2Dir)
			relPngPath := fmt.Sprintf("%s/%s%s_%d/*.png", relCL2Dir, trnDir, name, direction)
			pngPath := filepath.Join("_dump_", relPngPath)
			fmt.Fprintf(script, "\t%s \\\n", pngPath)
		}
	}
	fmt.Fprintf(script, "\t-gravity south -geometry %dx+0+0 \\\n", monster.FrameWidth)
	script.WriteString("\t-tile x8 \\\n")
	script.WriteString("\t-background none \\\n")
	dstName := monsterName(monster)
	dstPath := fmt.Sprintf("../mods/spark/images/monster/%s.png", dstName)
	fmt.Fprintf(script, "\t%s", dstPath)

	fmt.Println(script)
	return nil
}

// monsterName returns the unique file name of the given monster.
func monsterName(monster d1.MonsterData) string {
	name := snakeCase(monster.Name)
	cl2Path := strings.ToLower(monster.CL2Path)
	cl2Path = strings.Replace(cl2Path, `\`, "/", -1)
	// Resolve monster name collisions.
	switch {
	case strings.HasPrefix(cl2Path, "monsters/skelaxe/"):
		return name + "_axe"
	case strings.HasPrefix(cl2Path, "monsters/skelbow/"):
		return name + "_bow"
	case strings.HasPrefix(cl2Path, "monsters/falspear/"):
		return name + "_spear"
	case strings.HasPrefix(cl2Path, "monsters/falsword/"):
		return name + "_sword"
	case strings.HasPrefix(cl2Path, "monsters/goatmace/"):
		return name + "_mace"
	case strings.HasPrefix(cl2Path, "monsters/goatbow/"):
		return name + "_bow"
	}
	return name
}

// snakeCase returns a snake case version of the given monster name.
func snakeCase(name string) string {
	// TODO: Let monster categories (four different kinds of zombies) use the
	// same graphic.
	s := strings.ToLower(name)
	return strings.Replace(s, " ", "_", -1)
}

// extractMonsterSounds extracts the sounds of the given monster.
func extractMonsterSounds(monster d1.MonsterData) error {
	actions := []d1.MonsterAction{
		//d1.MonsterActionStand,
		//d1.MonsterActionWalk,
		d1.MonsterActionAttack,
		d1.MonsterActionHit,
		d1.MonsterActionDie,
	}
	if monster.HasSpecialSound {
		actions = append(actions, d1.MonsterActionSpecial)
	}
	script := &bytes.Buffer{}
	fmt.Fprintf(script, "echo 'Extracting sounds for %s'\n", monster.Name)
	// # Spitting Terror
	//
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acida1.wav ../mods/spark/sounds/monster/spitting_terror_attack_1.flac
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acida2.wav ../mods/spark/sounds/monster/spitting_terror_attack_2.flac
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidh1.wav ../mods/spark/sounds/monster/spitting_terror_hit_1.flac
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidh2.wav ../mods/spark/sounds/monster/spitting_terror_hit_2.flac
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidd1.wav ../mods/spark/sounds/monster/spitting_terror_die_1.flac
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidd2.wav ../mods/spark/sounds/monster/spitting_terror_die_2.flac
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acids1.wav ../mods/spark/sounds/monster/spitting_terror_special_1.flac
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acids2.wav ../mods/spark/sounds/monster/spitting_terror_special_2.flac
	for _, action := range actions {
		for i := 1; i <= 2; i++ {
			format := strings.ToLower(monster.WavPath)
			format = strings.Replace(format, `\`, "/", -1)
			format = strings.Replace(format, "%i", "%d", -1)
			relWavPath := fmt.Sprintf(format, action.Rune(), i)
			wavPath := filepath.Join("diabdat", relWavPath)
			fmt.Fprintf(script, "ffmpeg -loglevel error -y -i %s ../mods/spark/sounds/monster/%s_%s_%d.flac\n", wavPath, monsterName(monster), action.String(), i)
		}
	}
	fmt.Println(script)
	return nil
}
