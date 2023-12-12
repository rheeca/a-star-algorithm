package main

import (
	"embed"
	"fmt"

	"a-star/src/game"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var EmbeddedAssets embed.FS

func main() {
	gameObj := game.NewGame(EmbeddedAssets)

	// pathfinding
	// origin := "A"
	// dest := "G"

	// pathFound, err := AStar(origin, dest)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(pathFound)

	err := ebiten.RunGame(gameObj)
	if err != nil {
		fmt.Println("failed to run game:", err)
	}
}
