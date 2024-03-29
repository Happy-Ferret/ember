#!/bin/bash
gentilesetdef -mpqdir=../_assets_/diabdat -dtype town > ../mods/ember/tilesetdefs/tileset_tristram.txt
gentilesetdef -mpqdir=../_assets_/diabdat -dtype l1 > ../mods/ember/tilesetdefs/tileset_cathedral.txt
gentilesetdef -mpqdir=../_assets_/diabdat -dtype l2 > ../mods/ember/tilesetdefs/tileset_catacombs.txt
gentilesetdef -mpqdir=../_assets_/diabdat -dtype l3 > ../mods/ember/tilesetdefs/tileset_caves.txt
gentilesetdef -mpqdir=../_assets_/diabdat -dtype l4 > ../mods/ember/tilesetdefs/tileset_hell.txt

gentmx -mpqdir=../_assets_/diabdat ../_assets_/testdata/l1/l1_pillars_00000000.bin > ../tiled/cathedral/cathedral_00000000.tmx
