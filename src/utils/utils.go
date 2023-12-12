package utils

const (
	UnitSize = 32

	PlayerSpawnPoint   = 0
	ChickenSpawnPoints = 1

	GroundLayer    = 0
	CollisionLayer = 1

	PlayerFrameCount    = 8
	PlayerFrameDelay    = 8
	PlayerSpriteWidth   = 96
	PlayerSpriteHeight  = 96
	PlayerMovementSpeed = 2

	ChickenFrameCount    = 8
	ChickenFrameDelay    = 12
	ChickenMovementSpeed = 1

	IdleState = 0
	WalkState = 1

	ChickenIdleState = 0
	ChickenWalkState = 2
)

// Directions
const (
	Front = iota
	Back
	Left
	Right
	NumOfDirections
)

const (
	ChickenRight = iota
	ChickenLeft
	ChickenNumOfDirections
)

var (
	TilesetSheets   = []string{"grass_hill.png", "fences.png"}
	CollisionLayers = []int{CollisionLayer}
)
