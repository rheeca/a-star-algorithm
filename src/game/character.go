package game

import (
	"a-star/src/astar"
	"a-star/src/utils"
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

type Player struct {
	SpriteSheet *ebiten.Image
	XLoc        int
	YLoc        int
	Dx          int
	Dy          int
	State       int
	StateTTL    int
	Direction   int
	Frame       int
	Sprite      CollisionBody
	Collision   CollisionBody
}

type Chicken struct {
	SpriteSheet *ebiten.Image
	XLoc        int
	YLoc        int
	Dx          int
	Dy          int
	State       int
	StateTTL    int
	Direction   int
	Frame       int
	Sprite      CollisionBody
	Collision   CollisionBody
	Path        *astar.Path
}

func NewPlayer(embeddedAssets embed.FS, x, y int) *Player {
	return &Player{
		SpriteSheet: loadImage(embeddedAssets, "assets/player.png"),
		XLoc:        x - 39,
		YLoc:        y - 35,
		Sprite: CollisionBody{
			X:      x,
			Y:      y,
			Width:  20,
			Height: 30,
		},
		Collision: CollisionBody{
			X:      x,
			Y:      y + 15,
			Width:  18,
			Height: 16,
		},
	}
}

func NewChickens(embeddedAssets embed.FS, gameMap *tiled.Map) []*Chicken {
	chickens := []*Chicken{}
	spritesheet := loadImage(embeddedAssets, "assets/chicken.png")
	for _, spawnPoint := range gameMap.ObjectGroups[utils.ChickenSpawnPoints].Objects {
		xLoc := int(spawnPoint.X)
		yLoc := int(spawnPoint.Y)
		chicken := &Chicken{
			SpriteSheet: spritesheet,
			XLoc:        xLoc,
			YLoc:        yLoc,
			Sprite: CollisionBody{
				X:      xLoc,
				Y:      yLoc,
				Width:  32,
				Height: 32,
			},
			Collision: CollisionBody{
				X:      xLoc + 8,
				Y:      yLoc + 16,
				Width:  16,
				Height: 14,
			},
		}
		chickens = append(chickens, chicken)
	}
	return chickens
}

func (c *Chicken) UpdateFrame(currentFrame int) {
	if currentFrame%utils.ChickenFrameDelay == 0 {
		if c.StateTTL > 1 {
			c.StateTTL -= 1
		} else if c.StateTTL == 1 {
			c.StateTTL -= 1
			c.State = utils.IdleState
		}

		c.Frame += 1
		if c.Frame >= utils.ChickenFrameCount {
			c.Frame = 0
		}
	}
}

func (c *Chicken) SetPath(g *Game, x, y int) {
	cx, cy := c.GetCenterPoint()
	px, py := g.Player.GetCenterPoint()
	originCell := astar.GetCell(cx, cy)
	destCell := astar.GetCell(px, py)

	c.Path = astar.AStar(g.GridMap, originCell, destCell)
}

func (c *Chicken) UpdateLocation() {
	c.XLoc += c.Dx
	c.YLoc += c.Dy
	c.Sprite.X += c.Dx
	c.Sprite.Y += c.Dy
	c.Collision.X += c.Dx
	c.Collision.Y += c.Dy
	c.Dx = 0
	c.Dy = 0
}

func (p *Player) UpdateFrame(currentFrame int) {
	if currentFrame%utils.PlayerFrameDelay == 0 {
		if p.StateTTL > 1 {
			p.StateTTL -= 1
		} else if p.StateTTL == 1 {
			p.StateTTL -= 1
			p.State = utils.IdleState
		}

		p.Frame += 1
		if p.Frame >= utils.PlayerFrameCount {
			p.Frame = 0
		}
	}
}

func (p *Player) UpdateLocation() {
	p.XLoc += p.Dx
	p.YLoc += p.Dy
	p.Sprite.X += p.Dx
	p.Sprite.Y += p.Dy
	p.Collision.X += p.Dx
	p.Collision.Y += p.Dy
	p.Dx = 0
	p.Dy = 0
}

func (p *Player) GetCenterPoint() (x, y int) {
	return (p.Sprite.X + p.Sprite.Width/2), (p.Sprite.Y + p.Sprite.Height - 8)
}

func (c *Chicken) GetCenterPoint() (x, y int) {
	return (c.Sprite.X + c.Sprite.Width/2), (c.Sprite.Y + c.Sprite.Height/2)
}

func (p *Player) Restart(g *Game) {
	spawnPoint := g.GameMap.ObjectGroups[utils.PlayerSpawnPoint].Objects[0]
	x := int(spawnPoint.X)
	y := int(spawnPoint.Y)

	p.XLoc = x - 39
	p.YLoc = y - 35
	p.Sprite.X = x
	p.Sprite.Y = y
	p.Collision.X = x
	p.Collision.Y = y + 15

	p.Dx = 0
	p.Dy = 0
	p.State = 0
	p.StateTTL = 0
	p.Direction = 0
	p.Frame = 0
}

type CollisionBody struct {
	X      int
	Y      int
	Width  int
	Height int
}
