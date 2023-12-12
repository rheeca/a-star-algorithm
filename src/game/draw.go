package game

import (
	"a-star/src/utils"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

func drawMap(gMap *tiled.Map, tilesets map[string]*ebiten.Image, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, layer := range gMap.Layers {
		for tileY := 0; tileY < gMap.Height; tileY += 1 {
			for tileX := 0; tileX < gMap.Width; tileX += 1 {
				// find img of tile to draw
				tileToDraw := layer.Tiles[tileY*gMap.Width+tileX]
				if tileToDraw.IsNil() {
					continue
				}

				tileToDrawX := int(tileToDraw.ID) % tileToDraw.Tileset.Columns
				tileToDrawY := int(tileToDraw.ID) / tileToDraw.Tileset.Columns

				ebitenTileToDraw := tilesets[tileToDraw.Tileset.Name].SubImage(image.Rect(tileToDrawX*gMap.TileWidth,
					tileToDrawY*gMap.TileHeight,
					tileToDrawX*gMap.TileWidth+gMap.TileWidth,
					tileToDrawY*gMap.TileHeight+gMap.TileHeight)).(*ebiten.Image)

				// draw tile
				drawOptions.GeoM.Reset()
				TileXpos := float64(gMap.TileWidth * tileX)
				TileYpos := float64(gMap.TileHeight * tileY)
				drawOptions.GeoM.Translate(TileXpos, TileYpos)
				screen.DrawImage(ebitenTileToDraw, &drawOptions)
			}
		}
	}
}

func drawPlayer(player *Player, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(float64(player.XLoc), float64(player.YLoc))
	screen.DrawImage(player.SpriteSheet.SubImage(image.Rect(player.Frame*utils.PlayerSpriteWidth,
		(player.State*utils.NumOfDirections+player.Direction)*utils.PlayerSpriteHeight,
		player.Frame*utils.PlayerSpriteWidth+utils.PlayerSpriteWidth,
		(player.State*utils.NumOfDirections+player.Direction)*utils.PlayerSpriteHeight+utils.PlayerSpriteHeight)).(*ebiten.Image), &drawOptions)
}

func drawChickens(chickens []*Chicken, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, c := range chickens {
		drawOptions.GeoM.Reset()

		drawOptions.GeoM.Translate(float64(c.XLoc), float64(c.YLoc))
		screen.DrawImage(c.SpriteSheet.SubImage(image.Rect(c.Frame*c.Sprite.Width,
			(c.State*utils.ChickenNumOfDirections+c.Direction)*c.Sprite.Height,
			c.Frame*c.Sprite.Width+c.Sprite.Width,
			(c.State*utils.ChickenNumOfDirections+c.Direction)*c.Sprite.Height+c.Sprite.Height)).(*ebiten.Image), &drawOptions)
	}
}

func drawButtons(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	playImg := loadImage(g.EmbeddedAssets, "assets/play.png")
	restartImg := loadImage(g.EmbeddedAssets, "assets/restart.png")

	x := 0.0
	y := 0.0

	// draw play button
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(x, y)
	screen.DrawImage(playImg, &drawOptions)

	// draw restart button
	drawOptions.GeoM.Reset()
	drawOptions.GeoM.Translate(x, y+58)
	screen.DrawImage(restartImg, &drawOptions)
}
