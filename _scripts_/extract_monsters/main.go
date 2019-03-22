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

// Global command line flags.
var (
	// def specifies whether to extract monster definitions.
	def bool
	// graphics specifies whether to extract monster graphics.
	graphics bool
	// sounds specifies whether to extract monster sounds.
	sounds bool
)

func main() {
	// Parse command line arguments.
	var (
		// quiet specifies whether to suppress non-error messages.
		quiet bool
	)
	flag.Usage = usage
	flag.BoolVar(&def, "def", false, "extract monster definitions")
	flag.BoolVar(&graphics, "graphics", false, "extract monster graphics")
	flag.BoolVar(&quiet, "q", false, "suppress non-error messages")
	flag.BoolVar(&sounds, "sounds", false, "extract monster sounds")
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
	return nil
}

// extractMonster extracts the assets of the given monster.
func extractMonster(monster d1.MonsterData) error {
	dbg.Printf("extracting assets of %q.", monster.Name)
	// Extract monster graphics.
	if graphics {
		if err := extractMonsterGraphics(monster); err != nil {
			return errors.WithStack(err)
		}
	}
	// Extract monster sounds.
	if sounds {
		if err := extractMonsterSounds(monster); err != nil {
			return errors.WithStack(err)
		}
	}
	if def {
		if err := extractMonsterDef(monster); err != nil {
			return errors.WithStack(err)
		}
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
	//       ../mods/tristram/images/enemies/spitting_terror.png
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
	dstPath := fmt.Sprintf("../mods/tristram/images/monster/%s.png", dstName)
	fmt.Fprintf(script, "\t%s", dstPath)

	fmt.Println(script)
	return nil
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
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acida1.wav ../mods/tristram/sounds/monster/spitting_terror_attack_1.ogg
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acida2.wav ../mods/tristram/sounds/monster/spitting_terror_attack_2.ogg
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidh1.wav ../mods/tristram/sounds/monster/spitting_terror_hit_1.ogg
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidh2.wav ../mods/tristram/sounds/monster/spitting_terror_hit_2.ogg
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidd1.wav ../mods/tristram/sounds/monster/spitting_terror_die_1.ogg
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acidd2.wav ../mods/tristram/sounds/monster/spitting_terror_die_2.ogg
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acids1.wav ../mods/tristram/sounds/monster/spitting_terror_special_1.ogg
	//    ffmpeg -loglevel error -y -i diabdat/monsters/acid/acids2.wav ../mods/tristram/sounds/monster/spitting_terror_special_2.ogg
	for _, action := range actions {
		for i := 1; i <= 2; i++ {
			format := strings.ToLower(monster.WavPath)
			format = strings.Replace(format, `\`, "/", -1)
			format = strings.Replace(format, "%i", "%d", -1)
			relWavPath := fmt.Sprintf(format, action.Rune(), i)
			wavPath := filepath.Join("diabdat", relWavPath)
			fmt.Fprintf(script, "ffmpeg -loglevel error -y -i %s ../mods/tristram/sounds/monster/%s_%s_%d.ogg\n", wavPath, monsterName(monster), action.String(), i)
		}
	}
	fmt.Println(script)
	return nil
}

// extractMonsterDef extracts the definition of the given monster.
func extractMonsterDef(monster d1.MonsterData) error {
	// Create enemies/base/spitting_terror.txt
	//    sfx_attack=swing,sounds/monster/spitting_terror_attack_1.ogg
	//    sfx_attack=swing,sounds/monster/spitting_terror_attack_2.ogg
	//    sfx_attack=shoot,sounds/monster/spitting_terror_special_1.ogg
	//    sfx_attack=shoot,sounds/monster/spitting_terror_special_2.ogg
	//    sfx_attack=cast,sounds/monster/spitting_terror_special_1.ogg
	//    sfx_attack=cast,sounds/monster/spitting_terror_special_2.ogg
	//    sfx_block=sounds/powers/block.ogg
	//    sfx_critdie=sounds/monster/spitting_terror_die_1.ogg
	//    sfx_critdie=sounds/monster/spitting_terror_die_2.ogg
	//    sfx_die=sounds/monster/spitting_terror_die_1.ogg
	//    sfx_die=sounds/monster/spitting_terror_die_2.ogg
	//    sfx_hit=sounds/monster/spitting_terror_hit_1.ogg
	//    sfx_hit=sounds/monster/spitting_terror_hit_2.ogg

	// attack
	buf := &bytes.Buffer{}
	name := monsterName(monster)
	fmt.Fprintf(buf, "sfx_attack=swing,sounds/monster/%s_attack_1.ogg\n", name)
	// special
	if monster.HasSpecialSound {
		fmt.Fprintf(buf, "sfx_attack=shoot,sounds/monster/%s_special_1.ogg\n", name)
		fmt.Fprintf(buf, "sfx_attack=cast,sounds/monster/%s_special_1.ogg\n", name)
	}
	// block

	// TODO: figure out if monsters can block.
	//    sfx_block=sounds/powers/block.ogg
	fmt.Fprintf(buf, "sfx_block=soundfx/powers/block.ogg\n")
	// hit
	fmt.Fprintf(buf, "sfx_hit=sounds/monster/%s_hit_1.ogg\n", name)
	// die
	fmt.Fprintf(buf, "sfx_die=sounds/monster/%s_die_1.ogg\n", name)

	//    animations=animations/monster/spitting_terror.txt
	buf.WriteString("\n")
	fmt.Fprintf(buf, "animations=animations/monster/%s.txt\n", name)

	// TODO: figure out what melee_range and thread_range do and if thread_range
	// may cause performance problems when set too high.
	buf.WriteString("\n")
	buf.WriteString("melee_range=1.2\n")
	buf.WriteString("threat_range=600.0\n")

	// Store output.
	basePath := fmt.Sprintf("../mods/tristram/enemies/base/%s.txt", name)
	if err := ioutil.WriteFile(basePath, buf.Bytes(), 0644); err != nil {
		return errors.WithStack(err)
	}

	// Create enemies/spitting_terror.txt
	buf = &bytes.Buffer{}
	fmt.Fprintf(buf, "INCLUDE enemies/base/%s.txt\n", name)
	buf.WriteString("\n")
	fmt.Fprintf(buf, "name=%s\n", monster.Name)
	fmt.Fprintf(buf, "level=%d\n", monster.Level)
	fmt.Fprintf(buf, "categories=%s,dungeon\n", name)
	fmt.Fprintf(buf, "rarity=common\n")
	fmt.Fprintf(buf, "xp=%d\n", monster.Exp)
	buf.WriteString("\n")

	hp := monster.MinHP + (monster.MaxHP-monster.MinHP)/2
	buf.WriteString("# combat stats\n")
	fmt.Fprintf(buf, "stat=hp,%d\n", hp)
	// TODO: set speed from monster.Rate.
	fmt.Fprintf(buf, "speed=2\n")
	fmt.Fprintf(buf, "turn_delay=400ms\n")
	fmt.Fprintf(buf, "chance_pursue=10\n")
	buf.WriteString("\n")
	fmt.Fprintf(buf, "power=melee,1,2\n")
	fmt.Fprintf(buf, "power=ranged,32,2\n")
	buf.WriteString("\n")
	fmt.Fprintf(buf, "stat=accuracy,69\n")
	fmt.Fprintf(buf, "stat=avoidance,19\n")
	buf.WriteString("\n")
	fmt.Fprintf(buf, "stat=dmg_melee_min,%d\n", monster.MinDamage)
	fmt.Fprintf(buf, "stat=dmg_melee_max,%d\n", monster.MaxDamage)
	if monster.HasSpecialGraphic && monster.MinDamageSpecial != 0 {
		fmt.Fprintf(buf, "stat=dmg_ranged_min,%d\n", monster.MinDamageSpecial)
		fmt.Fprintf(buf, "stat=dmg_ranged_max,%d\n", monster.MaxDamageSpecial)
	}
	fmt.Fprintf(buf, "cooldown=1s\n")
	buf.WriteString("\n")
	buf.WriteString("# loot\n")
	fmt.Fprintf(buf, "loot=loot/leveled_low.txt\n")

	// Store output.
	defPath := fmt.Sprintf("../mods/tristram/enemies/%s.txt", name)
	if err := ioutil.WriteFile(defPath, buf.Bytes(), 0644); err != nil {
		return errors.WithStack(err)
	}

	// Create animations/monster/spitting_terror.txt
	buf = &bytes.Buffer{}
	//    image=images/monster/spitting_terror.png
	//    render_size=128,96
	//    render_offset=64,84
	//
	//    [swing]
	//    position=0
	//    frames=12
	//    duration=600ms
	//    type=play_once
	//
	//    [die]
	//    position=12
	//    frames=24
	//    duration=1200ms
	//    type=play_once
	//
	//    [hit]
	//    position=36
	//    frames=8
	//    duration=400ms
	//    type=play_once
	//
	//    [stance]
	//    position=44
	//    frames=13
	//    duration=650ms
	//    type=back_forth
	//
	//    [shoot]
	//    position=57
	//    frames=12
	//    duration=600ms
	//    type=play_once
	//
	//    [run]
	//    position=69
	//    frames=8
	//    duration=400ms
	//    type=looped

	fmt.Fprintf(buf, "image=images/monster/%s.png\n", name)
	// TODO: figure out how to get the value of height.
	frameHeight := 96
	fmt.Fprintf(buf, "render_size=%d,%d\n", monster.FrameWidth, frameHeight)
	// TODO: set offset to
	fmt.Fprintf(buf, "render_offset=%d,%d\n", monster.FrameWidth/2, frameHeight-16)
	buf.WriteString("\n")

	// stand
	position := int32(0)
	nframes := monster.NFrames[d1.MonsterActionStand]
	// Diablo 1 runs at 20 FPS; thus 50ms per frame.
	duration := 50 * nframes
	fmt.Fprintf(buf, "[stance]\n")
	fmt.Fprintf(buf, "position=%d\n", position)
	fmt.Fprintf(buf, "frames=%d\n", nframes)
	fmt.Fprintf(buf, "duration=%dms\n", duration)
	fmt.Fprintf(buf, "type=back_forth\n")
	position += nframes

	// walk
	nframes = monster.NFrames[d1.MonsterActionWalk]
	duration = 50 * nframes
	fmt.Fprintf(buf, "[run]\n")
	fmt.Fprintf(buf, "position=%d\n", position)
	fmt.Fprintf(buf, "frames=%d\n", nframes)
	fmt.Fprintf(buf, "duration=%dms\n", nframes)
	fmt.Fprintf(buf, "type=looped\n")
	position += nframes

	// attack
	nframes = monster.NFrames[d1.MonsterActionAttack]
	duration = 50 * nframes
	fmt.Fprintf(buf, "[swing]\n")
	fmt.Fprintf(buf, "position=%d\n", position)
	fmt.Fprintf(buf, "frames=%d\n", nframes)
	fmt.Fprintf(buf, "duration=%dms\n", nframes)
	fmt.Fprintf(buf, "type=play_once\n")
	position += nframes

	// hit
	nframes = monster.NFrames[d1.MonsterActionHit]
	duration = 50 * nframes
	fmt.Fprintf(buf, "[hit]\n")
	fmt.Fprintf(buf, "position=%d\n", position)
	fmt.Fprintf(buf, "frames=%d\n", nframes)
	fmt.Fprintf(buf, "duration=%dms\n", nframes)
	fmt.Fprintf(buf, "type=play_once\n")
	position += nframes

	// die
	nframes = monster.NFrames[d1.MonsterActionDie]
	duration = 50 * nframes
	fmt.Fprintf(buf, "[die]\n")
	fmt.Fprintf(buf, "position=%d\n", position)
	fmt.Fprintf(buf, "frames=%d\n", nframes)
	fmt.Fprintf(buf, "duration=%dms\n", nframes)
	fmt.Fprintf(buf, "type=play_once\n")
	position += nframes

	// special
	nframes = monster.NFrames[d1.MonsterActionSpecial]
	if monster.HasSpecialGraphic {
		duration = 50 * nframes
		fmt.Fprintf(buf, "[shoot]\n")
		fmt.Fprintf(buf, "position=%d\n", position)
		fmt.Fprintf(buf, "frames=%d\n", nframes)
		fmt.Fprintf(buf, "duration=%dms\n", nframes)
		fmt.Fprintf(buf, "type=play_once\n")
	}

	animPath := fmt.Sprintf("../mods/tristram/animations/monster/%s.txt", name)
	if err := ioutil.WriteFile(animPath, buf.Bytes(), 0644); err != nil {
		return errors.WithStack(err)
	}

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
