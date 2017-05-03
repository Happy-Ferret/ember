package main

// TODO: Blit arches onto their corresponding dungeon piece before generating
// the tileset.

// Dungeon piece IDs for layout 1.
const (
	// Floor shadows for arches.
	DPieceNone                       = 0
	DPieceFloorShadowArchSe_1        = 11
	DPieceFloorShadowArchSe_2        = 249
	DPieceFloorShadowArchSe_3        = 325
	DPieceFloorShadowArchSe_4        = 331
	DPieceFloorShadowArchSe_5        = 344
	DPieceFloorShadowArchSe_6        = 421
	DPieceFloorShadowArchSw2_1       = 259
	DPieceFloorShadowArchSwBroken2_1 = 255
	DPieceFloorShadowArchSw_1        = 12
	DPieceFloorShadowArchSw_2        = 71
	DPieceFloorShadowArchSw_3        = 211
	DPieceFloorShadowArchSw_4        = 321
	DPieceFloorShadowArchSw_5        = 341
	DPieceFloorShadowArchSw_6        = 418
)

// Arch IDs for layout 1.
const (
	ArchNone      = 0
	ArchSe        = 2
	ArchSeBroken  = 3
	ArchSeDoor    = 8
	ArchSw        = 1
	ArchSw2       = 5
	ArchSwBroken  = 6
	ArchSwBroken2 = 4
	ArchSwDoor    = 7
)

// archID returns the arch ID of the given dungeon piece ID.
//
// ref: 46E9E2
func archID(dpieceID, firstArchID int) int {
	switch dpieceID {
	case DPieceFloorShadowArchSw_1, DPieceFloorShadowArchSw_2, DPieceFloorShadowArchSw_3, DPieceFloorShadowArchSw_4, DPieceFloorShadowArchSw_5, DPieceFloorShadowArchSw_6:
		return firstArchID + ArchSw
	case DPieceFloorShadowArchSe_1, DPieceFloorShadowArchSe_2, DPieceFloorShadowArchSe_3, DPieceFloorShadowArchSe_4, DPieceFloorShadowArchSe_5, DPieceFloorShadowArchSe_6:
		return firstArchID + ArchSe
	case DPieceFloorShadowArchSwBroken2_1:
		return firstArchID + ArchSwBroken2
	case DPieceFloorShadowArchSw2_1:
		return firstArchID + ArchSw2
	default:
		return 0
	}
}
