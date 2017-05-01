# Spark

The aim of this project is to get Diablo 1 running on the FLARE game engine.

## Current progress

* [x] Convert map to TMX format.
    - [x] Tristram
    - [x] Cathedral
    - [ ] Catacombs
    - [ ] Caves
    - [ ] Hell

![Tristram](http://oi65.tinypic.com/juk2ed.jpg)

![Cathedral](http://oi68.tinypic.com/iof5es.jpg)

## Installation

Note, this game requires an original copy of `diabdat.mpq`. None of the Diablo 1 game assets are provided by this project.

### Install FLARE engine

```bash
# Install dependencies of FLARE.
pacman -S sdl2 sdl2_image sdl2_mixer sdl2_ttf cmake ffmpeg

# Clone FLARE engine and game assets.
git clone https://github.com/clintbellanger/flare-engine
git clone https://github.com/clintbellanger/flare-game

# Clone the Spark game.
git clone https://github.com/sanctuary/spark

# Build FLARE engine.
cd flare-engine
cmake .
make

# Add symlinks to default, fantasycore and alpha_demo mods.
cd ../spark/mods
ln -s ../../flare-engine/mods/default
ln -s ../../flare-game/mods/fantasycore
ln -s ../../flare-game/mods/alpha_demo
cd ..
ln -s ../flare-engine/flare
```

### Generate game assets

```bash
# Get assets conversion tools.
go get github.com/mewrnd/blizzconv/...
go get github.com/sanctuary/spark/_scripts_/opensourceami

# Create "spark/_assets_" directory.
mkdir _assets_
cd _assets_

# Extract diabdat.mpq to the "_assets_/mpqdump" directory.
#
# #############################################################
# ### NOTE: This step requires manual intervention for now. ###
# #############################################################
#
# You may use Ladislav Zezula's MPQ Editor to extract the contents of diabdat.mpq.

# Extract game assets. Takes roughly 5 minutes.
opensourceami -o extract_assets.sh
./extract_assets.sh
cd ..
```

### Run the game

```bash
# Standing in the spark directory, run `./flare --mods=spark`
# to start the game.
./flare --mods=spark
```

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
