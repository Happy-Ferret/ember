package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mewkiz/pkg/imgutil"
	"github.com/mewkiz/pkg/term"
	"github.com/pkg/errors"
	"github.com/sanctuary/formats/image/cel/config"
)

// dbg represents a logger with the "fixarches:" prefix, which logs debug
// messages to standard error.
var dbg = log.New(os.Stderr, term.BlueBold("fixarches:")+" ", 0)

func main() {
	// Parse command line flags.
	var (
		// mpqDir specifies the path to an extracted "diabdat.mpq".
		mpqDir string
	)
	flag.StringVar(&mpqDir, "mpqdir", "diabdat", `path to extracted "diabdat.mpq"`)
	flag.Parse()

	dtypes := []string{"l1", "l2", "l3", "l4", "town"}
	for _, dtype := range dtypes {
		if err := fixArches(dtype, mpqDir); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

// fixArches draws arches on the relevant tiles of the given dungeon type.
func fixArches(dtype, mpqDir string) error {
	// Parse SOL file.
	relSolPath := fmt.Sprintf("levels/%sdata/%s.sol", dtype, dtype)
	solPath := filepath.Join(mpqDir, relSolPath)
	sol, err := ioutil.ReadFile(solPath)
	if err != nil {
		return errors.WithStack(err)
	}

	// Number of dungeon pieces in the tileset.
	ndpieces := len(sol)
	for dpieceID := 1; dpieceID <= ndpieces; dpieceID++ {
		archID := getArch(dpieceID, dtype)
		if archID == ArchNone {
			continue
		}
		relCelPath := fmt.Sprintf("levels/%sdata/%ss.cel", dtype, dtype)
		celName := filepath.Base(relCelPath)
		conf, err := config.Get(celName)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, relPalPath := range conf.Pals {
			palName := filepath.Base(relPalPath)
			dpiecePath := filepath.Join("_dump_", fmt.Sprintf("_dpieces_/%s/%s/dpiece_%04d.png", dtype, palName, dpieceID))
			dpieceImg, err := imgutil.ReadFile(dpiecePath)
			if err != nil {
				return errors.WithStack(err)
			}
			archPath := filepath.Join("_dump_", fmt.Sprintf("levels/%sdata/%ss/%s/%ss_%04d.png", dtype, dtype, palName, dtype, archID))
			archImg, err := imgutil.ReadFile(archPath)
			if err != nil {
				return errors.WithStack(err)
			}
			bounds := dpieceImg.Bounds()
			rect := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
			dst := image.NewRGBA(rect)
			draw.Draw(dst, rect, dpieceImg, image.ZP, draw.Src)
			draw.Draw(dst, rect, archImg, image.ZP, draw.Over)
			dbg.Printf("Drawing arch ID %d onto dungeon piece ID %d with palette %q.", archID, dpieceID, relPalPath)
			if err := imgutil.WriteFile(dpiecePath, dst); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

// getArch returns the arch ID of the given dungeon piece.
func getArch(dpieceID int, dtype string) int {
	switch dtype {
	case "l1":
		return l1ArchID(dpieceID)
	case "l2":
		// TODO: Implement support for layout 2.
		return 0
	case "l3":
		// nothing to do; layout 3 has no arches.
		return ArchNone
	case "l4":
		// nothing to do; layout 4 has no arches.
		return ArchNone
	case "town":
		// TODO: Implement support for town.
		return 0
	default:
		panic(fmt.Errorf("support for dungeon type %q not yet implemented", dtype))
	}
}

const ArchNone = 0

// Arch IDs for layout 1.
const (
	L1ArchSe        = 2
	L1ArchSeBroken  = 3
	L1ArchSeDoor    = 8
	L1ArchSw        = 1
	L1ArchSw2       = 5
	L1ArchSwBroken  = 6
	L1ArchSwBroken2 = 4
	L1ArchSwDoor    = 7
)

// l1ArchID returns the arch ID of the given dungeon piece.
//
// ref: 46E9E2
func l1ArchID(dpieceID int) int {
	// Dungeon piece IDs for layout 1.
	const (
		// Floor shadows for arches.
		L1DPieceNone                       = 0
		L1DPieceFloorShadowArchSe_1        = 11
		L1DPieceFloorShadowArchSe_2        = 249
		L1DPieceFloorShadowArchSe_3        = 325
		L1DPieceFloorShadowArchSe_4        = 331
		L1DPieceFloorShadowArchSe_5        = 344
		L1DPieceFloorShadowArchSe_6        = 421
		L1DPieceFloorShadowArchSw2_1       = 259
		L1DPieceFloorShadowArchSwBroken2_1 = 255
		L1DPieceFloorShadowArchSw_1        = 12
		L1DPieceFloorShadowArchSw_2        = 71
		L1DPieceFloorShadowArchSw_3        = 211
		L1DPieceFloorShadowArchSw_4        = 321
		L1DPieceFloorShadowArchSw_5        = 341
		L1DPieceFloorShadowArchSw_6        = 418
	)
	switch dpieceID {
	case L1DPieceFloorShadowArchSw_1, L1DPieceFloorShadowArchSw_2, L1DPieceFloorShadowArchSw_3, L1DPieceFloorShadowArchSw_4, L1DPieceFloorShadowArchSw_5, L1DPieceFloorShadowArchSw_6:
		return L1ArchSw
	case L1DPieceFloorShadowArchSe_1, L1DPieceFloorShadowArchSe_2, L1DPieceFloorShadowArchSe_3, L1DPieceFloorShadowArchSe_4, L1DPieceFloorShadowArchSe_5, L1DPieceFloorShadowArchSe_6:
		return L1ArchSe
	case L1DPieceFloorShadowArchSwBroken2_1:
		return L1ArchSwBroken2
	case L1DPieceFloorShadowArchSw2_1:
		return L1ArchSw2
	default:
		return ArchNone
	}
}
