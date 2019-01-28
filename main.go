package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/fruit-planet/world"
)

func main() {
	w := world.New("default", 12*16, 6*16)

	if err := ebiten.Run(w.Update, 12*16, 6*16, 2, "Fruit World"); err != nil {
		log.Fatal(err)
	}
}
