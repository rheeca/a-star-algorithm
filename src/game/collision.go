package game

import (
	"a-star/src/utils"

	"github.com/co0p/tankism/lib/collision"
)

func playerHasCollisions(g *Game, player *Player) bool {
	if hasMapCollisions(g, player.Dx, player.Dy, player.Collision) {
		return true
	}
	return false
}

func hasMapCollisions(g *Game, dx, dy int, collisionBody CollisionBody) bool {
	for tileY := 0; tileY < g.GameMap.Height; tileY += 1 {
		for tileX := 0; tileX < g.GameMap.Width; tileX += 1 {
			for _, layer := range utils.CollisionLayers {
				tile := g.GameMap.Layers[layer].Tiles[tileY*g.GameMap.Width+tileX]
				if tile.IsNil() {
					continue
				}
				tileXpos := g.GameMap.TileWidth * tileX
				tileYpos := g.GameMap.TileHeight * tileY

				tileCollision := CollisionBody{
					X:      tileXpos,
					Y:      tileYpos,
					Width:  g.GameMap.TileWidth,
					Height: g.GameMap.TileHeight,
				}
				if hasCollision(dx, dy, collisionBody, tileCollision) {
					return true
				}
			}
		}
	}
	return false
}

func hasCollision(dx, dy int, bodyA, bodyB CollisionBody) bool {
	// check if movement of bodyA collides with bodyB
	aBounds := collision.BoundingBox{
		X:      float64(bodyA.X + dx),
		Y:      float64(bodyA.Y + dy),
		Width:  float64(bodyA.Width),
		Height: float64(bodyA.Height),
	}
	bBounds := collision.BoundingBox{
		X:      float64(bodyB.X),
		Y:      float64(bodyB.Y),
		Width:  float64(bodyB.Width),
		Height: float64(bodyB.Height),
	}
	if collision.AABBCollision(aBounds, bBounds) {
		return true
	}
	return false
}
