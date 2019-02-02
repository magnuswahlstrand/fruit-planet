package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/fruit-planet/world"
)

func main() {
	// screenWidth, screenHeight := 12*16,6*16
	screenWidth, screenHeight := 24*16, 18*16
	w := world.New("world-2", screenWidth, screenHeight)

	if err := ebiten.Run(w.Update, screenWidth, screenHeight, 2, "Fruit World"); err != nil {
		log.Fatal(err)
	}
}
