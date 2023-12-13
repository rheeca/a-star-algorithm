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

	err := ebiten.RunGame(gameObj)
	if err != nil {
		fmt.Println("failed to run game:", err)
	}
}
