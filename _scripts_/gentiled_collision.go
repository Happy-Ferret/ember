//+build ignore

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/mewkiz/pkg/imgutil"
	"github.com/pkg/errors"
	"golang.org/x/image/colornames"
)

func usage() {
	const use = `
Generate tiled_collision.png from a given mask image.

Usage:

	gentiled_collision [OPTION]... FILE.png
`
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	maskPath := flag.Arg(0)
	mask, err := imgutil.ReadFile(maskPath)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	const ncollisions = 40
	i := 0
	for _, c := range colornames.Map {
		if i >= ncollisions {
			break
		}
		dst := gen(mask, c)
		out := fmt.Sprintf("mask_%04d.png", i)
		if err := imgutil.WriteFile(out, dst); err != nil {
			log.Fatalf("%+v", errors.WithStack(err))
		}
		i++
	}
}

func gen(mask image.Image, c color.Color) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, 64, 32))
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			c0 := mask.At(x, y)
			if _, _, _, a := c0.RGBA(); a != 0 {
				dst.Set(x, y, c)
			}
		}
	}
	return dst
}
