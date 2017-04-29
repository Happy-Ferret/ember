// generate "tilesetdefs/tileset_town.txt"

package main

import "fmt"

func main() {
	fmt.Println("img=images/tilesets/tileset_town.png")
	fmt.Println()
	x, y := 0, 0
	for id := 1; id <= 1258; id++ {
		fmt.Printf("tile=%d,%d,%d,64,256,32,240\n", id+100, x*64, y*256)
		x++
		if x >= 79 {
			x = 0
			y++
		}
	}
}
