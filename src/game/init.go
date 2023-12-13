package game

import (
	"a-star/src/astar"
	"a-star/src/utils"
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/lafriks/go-tiled"
)

type Game struct {
	GameMap        *tiled.Map
	GridMap        *astar.GridMap
	Tilesets       map[string]*ebiten.Image
	Player         *Player
	Chickens       []*Chicken
	CurrentFrame   int
	EmbeddedAssets embed.FS
}

func NewGame(embeddedAssets embed.FS) *Game {
	gameMap, err := loadMapFromEmbedded(embeddedAssets, "assets/map.tmx")
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	windowWidth := gameMap.Width * gameMap.TileWidth
	windowHeight := gameMap.Height * gameMap.TileHeight
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("A Star Algorithm Implementation")

	spawnPoint := gameMap.ObjectGroups[utils.PlayerSpawnPoint].Objects[0]
	player := NewPlayer(embeddedAssets, int(spawnPoint.X), int(spawnPoint.Y))
	chickens := NewChickens(embeddedAssets, gameMap)

	return &Game{
		GameMap:        gameMap,
		GridMap:        astar.NewGridMap(gameMap),
		Tilesets:       getTilesets(embeddedAssets),
		Player:         player,
		Chickens:       chickens,
		EmbeddedAssets: embeddedAssets,
	}
}

func (g *Game) Update() error {
	g.CurrentFrame += 1
	g.Player.UpdateFrame(g.CurrentFrame)
	getPlayerInput(g)

	// update chickens
	for i, c := range g.Chickens {
		g.Chickens[i].UpdateFrame(g.CurrentFrame)

		// if chicken has a path, walk to path
		if c.Path != nil {
			nextCell := c.Path.GetCurrentCell()
			if nextCell == nil {
				// already at destination
				if c.State == utils.ChickenWalkState {
					g.Chickens[i].State = utils.ChickenIdleState
				}
				continue
			}

			if nextCell.X*32 == c.XLoc && nextCell.Y*32 == c.YLoc {
				g.Chickens[i].Path.Next()
				continue
			}

			if nextCell.X*32 < c.XLoc {
				// walk left
				g.Chickens[i].Direction = utils.ChickenLeft
				g.Chickens[i].State = utils.ChickenWalkState
				g.Chickens[i].Dx -= utils.ChickenMovementSpeed
				g.Chickens[i].UpdateLocation()
			} else if nextCell.X*32 > c.XLoc {
				// walk right
				g.Chickens[i].Direction = utils.ChickenRight
				g.Chickens[i].State = utils.ChickenWalkState
				g.Chickens[i].Dx += utils.ChickenMovementSpeed
				g.Chickens[i].UpdateLocation()
			} else if nextCell.Y*32 < c.YLoc {
				// walk down
				g.Chickens[i].State = utils.ChickenWalkState
				g.Chickens[i].Dy -= utils.ChickenMovementSpeed
				g.Chickens[i].UpdateLocation()
			} else if nextCell.Y*32 > c.YLoc {
				// walk up
				g.Chickens[i].State = utils.ChickenWalkState
				g.Chickens[i].Dy += utils.ChickenMovementSpeed
				g.Chickens[i].UpdateLocation()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawOptions := ebiten.DrawImageOptions{}
	drawMap(g.GameMap, g.Tilesets, screen, drawOptions)
	drawPlayer(g.Player, screen, drawOptions)
	drawChickens(g.Chickens, screen, drawOptions)
	drawButtons(g, screen, drawOptions)
}

func (g *Game) Layout(oWidth, oHeight int) (sWidth, sHeight int) {
	return oWidth, oHeight
}

func getPlayerInput(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.Player.Sprite.X > 0 {
		g.Player.Direction = utils.Left
		g.Player.State = utils.WalkState
		g.Player.Dx -= utils.PlayerMovementSpeed
		if !playerHasCollisions(g, g.Player) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) &&
		g.Player.Sprite.X < (g.GameMap.Width*g.GameMap.TileWidth)-g.Player.Sprite.Width {
		g.Player.Direction = utils.Right
		g.Player.State = utils.WalkState
		g.Player.Dx += utils.PlayerMovementSpeed
		if !playerHasCollisions(g, g.Player) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dx = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyW) && g.Player.Sprite.Y > 0 {
		g.Player.Direction = utils.Back
		g.Player.State = utils.WalkState
		g.Player.Dy -= utils.PlayerMovementSpeed
		if !playerHasCollisions(g, g.Player) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dy = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) &&
		g.Player.Sprite.Y < (g.GameMap.Height*g.GameMap.TileHeight)-g.Player.Sprite.Height {
		g.Player.Direction = utils.Front
		g.Player.State = utils.WalkState
		g.Player.Dy += utils.PlayerMovementSpeed
		if !playerHasCollisions(g, g.Player) {
			g.Player.UpdateLocation()
		} else {
			g.Player.Dy = 0
		}
	} else if g.Player.StateTTL == 0 {
		g.Player.State = utils.IdleState
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		if isClicked(mouseX, mouseY, CollisionBody{
			X:      6,
			Y:      4,
			Width:  180,
			Height: 54,
		}) {
			for i := range g.Chickens {
				g.Chickens[i].SetPath(g, g.Player.XLoc, g.Player.YLoc)
			}
		} else if isClicked(mouseX, mouseY, CollisionBody{
			X:      6,
			Y:      4 + 58,
			Width:  180,
			Height: 54,
		}) {
			g.Player.Restart(g)
		}
	}
}

func isClicked(x, y int, body CollisionBody) bool {
	// check if mouse clicked on a body
	aBounds := collision.BoundingBox{
		X:      float64(x),
		Y:      float64(y),
		Width:  1,
		Height: 1,
	}
	bBounds := collision.BoundingBox{
		X:      float64(body.X),
		Y:      float64(body.Y),
		Width:  float64(body.Width),
		Height: float64(body.Height),
	}
	if collision.AABBCollision(aBounds, bBounds) {
		return true
	}
	return false
}

func loadImage(EmbeddedAssets embed.FS, filepath string) *ebiten.Image {
	embeddedFile, err := EmbeddedAssets.Open(filepath)
	if err != nil {
		return nil
	}
	image, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		return nil
	}
	return image
}

func loadMapFromEmbedded(embeddedAssets embed.FS, name string) (*tiled.Map, error) {
	embeddedMap, err := tiled.LoadFile(name,
		tiled.WithFileSystem(embeddedAssets))
	if err != nil {
		return nil, err
	}
	return embeddedMap, nil
}

func getTilesets(embeddedAssets embed.FS) map[string]*ebiten.Image {
	tilesets := map[string]*ebiten.Image{}
	for _, tsPath := range utils.TilesetSheets {
		embeddedFile, err := embeddedAssets.Open(path.Join("assets", tsPath))
		if err != nil {
			log.Fatal("failed to load embedded image:", embeddedFile, err)
		}
		tsImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
		if err != nil {
			fmt.Println("error loading tileset image")
		}
		tilesets[strings.Split(tsPath, ".")[0]] = tsImage
	}
	return tilesets
}
