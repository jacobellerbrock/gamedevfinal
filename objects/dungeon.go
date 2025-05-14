package objects

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Block int
type Color int

const (
	NONE      = 0
	FLOOR     = 1
	WALL      = 2
	REDSPAWN  = 3
	BLUESPAWN = 4
)

const (
	WHITE = 0
	RED   = 1
	BLUE  = 2
)

type Dungeon struct {
	Blocks         [][]Block
	Colors         [][]Color
	BlockSize      int
	BlockSprites   []rl.Texture2D
	BlockColors    []rl.Color
	FloorPositions []rl.Vector2
	Width          int
	Height         int
	Players        []*Player
}

type Room struct {
	PivotX int
	PivotY int
	Width  int
	Height int
}

func NewDungeon(width, height int) Dungeon {
	dungeon := Dungeon{}

	dungeon.Width = width
	dungeon.Height = height

	dungeon.Blocks = make([][]Block, width)
	dungeon.Colors = make([][]Color, width)
	for x := 0; x < width; x++ {
		dungeon.Blocks[x] = make([]Block, height)
		dungeon.Colors[x] = make([]Color, height)
	}
	noneTexture := rl.LoadTexture("resources/sprites/none.png")
	dotTexture := rl.LoadTexture("resources/sprites/dot.png")
	wallTexture := rl.LoadTexture("resources/sprites/rounded.png")
	spawnTexture := rl.LoadTexture("resources/sprites/spawn.png")
	dungeon.BlockSprites = append(dungeon.BlockSprites, noneTexture)
	dungeon.BlockSprites = append(dungeon.BlockSprites, dotTexture)
	dungeon.BlockSprites = append(dungeon.BlockSprites, wallTexture)
	dungeon.BlockSprites = append(dungeon.BlockSprites, spawnTexture)

	dungeon.BlockColors = append(dungeon.BlockColors, rl.NewColor(0, 0, 0, 10))
	dungeon.BlockColors = append(dungeon.BlockColors, rl.NewColor(0, 0, 0, 64))
	dungeon.BlockColors = append(dungeon.BlockColors, rl.NewColor(0, 0, 0, 128))

	dungeon.BlockSize = int(dotTexture.Width)

	dungeon.FloorPositions = make([]rl.Vector2, 0)
	dungeon.Players = make([]*Player, 0)
	return dungeon
}

func (d *Dungeon) PlaceBlock(block Block, x, y int) {
	d.Blocks[x][y] = block

	if block == FLOOR {
		d.FloorPositions = append(d.FloorPositions, rl.NewVector2(float32(x), float32(y)))
	}
}

func (d *Dungeon) GetBlock(x, y int) Block {
	return d.Blocks[x][y]
}

func (d *Dungeon) GetBlockPosition(x, y int) rl.Vector2 {
	return rl.Vector2{X: float32(x * d.BlockSize), Y: float32(y * d.BlockSize)}
}

func (d *Dungeon) DrawDungeon() {
	for x := 0; x < len(d.Blocks); x++ {
		for y := 0; y < len(d.Blocks[x]); y++ {
			if d.Blocks[x][y] == NONE {
				continue
			} else if d.Blocks[x][y] == REDSPAWN {
				rl.DrawTexture(
					d.BlockSprites[3],
					int32(x*d.BlockSize),
					int32(y*d.BlockSize),
					rl.NewColor(255, 0, 0, 255),
				)
			} else if d.Blocks[x][y] == BLUESPAWN {
				rl.DrawTexture(
					d.BlockSprites[3],
					int32(x*d.BlockSize),
					int32(y*d.BlockSize),
					rl.NewColor(0, 0, 255, 255),
				)
			} else if d.Blocks[x][y] == FLOOR && d.Colors[x][y] == RED {
				rl.DrawRectangle(int32(x*d.BlockSize), int32(y*d.BlockSize), int32(d.BlockSize), int32(d.BlockSize), rl.NewColor(255, 0, 0, 64))
			} else if d.Blocks[x][y] == FLOOR && d.Colors[x][y] == BLUE {
				rl.DrawRectangle(int32(x*d.BlockSize), int32(y*d.BlockSize), int32(d.BlockSize), int32(d.BlockSize), rl.NewColor(0, 0, 255, 64))
			} else {
				rl.DrawTexture(
					d.BlockSprites[d.Blocks[x][y]],
					int32(x*d.BlockSize),
					int32(y*d.BlockSize),
					d.BlockColors[d.Blocks[x][y]],
				)
			}
		}
	}
}

func RandInt(drand *rand.Rand, min, max int) int {
	return drand.IntN(max-min) + min
}

func (d *Dungeon) StampRoom(room Room) {
	for x := 0; x < room.Width; x++ {
		for y := 0; y < room.Height; y++ {
			if x == 0 || y == 0 || x == room.Width-1 || y == room.Height-1 {
				d.PlaceBlock(WALL, x+room.PivotX, y+room.PivotY)
			} else {
				d.PlaceBlock(FLOOR, x+room.PivotX, y+room.PivotY)
			}
		}
	}
}

func (r *Room) MoveTo(moveTo rl.Vector2) {
	if r.PivotX < int(moveTo.X) {
		r.PivotX += 1
	}
	if r.PivotX > int(moveTo.X) {
		r.PivotX -= 1
	}
	if r.PivotY < int(moveTo.Y) {
		r.PivotY += 1
	}
	if r.PivotY > int(moveTo.Y) {
		r.PivotY -= 1
	}
}

func (r Room) RoomInBounds(dungeon *Dungeon) bool {
	if r.PivotX < 0 || r.PivotY < 0 {
		return false
	}
	if r.PivotX+r.Width >= dungeon.Width {
		return false
	}
	if r.PivotY+r.Height >= dungeon.Height {
		return false
	}
	return true
}

func (r Room) BadOverlap(dungeon *Dungeon) bool {
	for x := 1; x < r.Width-1; x++ {
		for y := 1; y < r.Height-1; y++ {
			if dungeon.GetBlock(x+r.PivotX, y+r.PivotY) != NONE {
				return true
			}
		}
	}
	return false
}

func (dungeon *Dungeon) AttemptStamp(r Room) bool {
	//if we overlap properly, then
	for x := 0; x < r.Width; x++ {
		for y := 0; y < r.Height; y++ {
			//left right overlap
			if (x == 0 || x == r.Width-1) && (y > 0 && y < r.Height-1) {
				if dungeon.GetBlock(x+r.PivotX, y+r.PivotY) == WALL {
					if dungeon.GetBlock(x+r.PivotX-1, y+r.PivotY) == FLOOR || dungeon.GetBlock(x+r.PivotX+1, y+r.PivotY) == FLOOR {
						dungeon.StampRoom(r)
						dungeon.PlaceBlock(FLOOR, x+r.PivotX, y+r.PivotY)
						return true
					}
				}
			}
			//top bottom overlap
			if (y == 0 || y == r.Height-1) && (x > 0 && x < r.Width-1) {
				if dungeon.GetBlock(x+r.PivotX, y+r.PivotY) == WALL {
					if dungeon.GetBlock(x+r.PivotX, y+r.PivotY-1) == FLOOR || dungeon.GetBlock(x+r.PivotX, y+r.PivotY+1) == FLOOR {
						dungeon.StampRoom(r)
						dungeon.PlaceBlock(FLOOR, x+r.PivotX, y+r.PivotY)
						return true
					}
				}
			}
		}
	}
	return false
}

func (d *Dungeon) Generate(seed1, seed2 uint64) {
	//clear old blocks out of grid
	d.Blocks = make([][]Block, d.Width)
	d.Colors = make([][]Color, d.Width)
	for x := 0; x < d.Width; x++ {
		d.Blocks[x] = make([]Block, d.Height)
		d.Colors[x] = make([]Color, d.Height)
	}

	pcg := rand.NewPCG(seed1, seed2)
	drand := rand.New(pcg)

	//place initial room
	initialRoom := Room{0, 0, RandInt(drand, 5, 8), RandInt(drand, 5, 8)}

	initialRoom.PivotX = RandInt(drand, 0, d.Width-initialRoom.Width)
	initialRoom.PivotY = RandInt(drand, 0, d.Height-initialRoom.Height)
	d.StampRoom(initialRoom)

	attemptsLeft := 50

	for attemptsLeft > 0 {
		attemptsLeft -= 1

		//create a random room
		tempRoom := Room{0, 0, RandInt(drand, 3, 8), RandInt(drand, 3, 8)}

		spawnDirection := RandInt(drand, 0, 4)
		switch spawnDirection {
		case 0:
			tempRoom.PivotX = -tempRoom.Width
			tempRoom.PivotY = RandInt(drand, -d.Height, d.Height)
		case 1:
			tempRoom.PivotX = d.Width
			tempRoom.PivotY = RandInt(drand, -d.Height, d.Height)
		case 2:
			tempRoom.PivotY = -tempRoom.Height
			tempRoom.PivotY = RandInt(drand, -d.Width, d.Width)
		case 3:
			tempRoom.PivotY = d.Height
			tempRoom.PivotY = RandInt(drand, -d.Width, d.Width)
		}

		moveTo := d.FloorPositions[RandInt(drand, 0, len(d.FloorPositions))]

		moveAttempts := int(math.Max(float64(d.Width), float64(d.Height)))
		for moveAttempts > 0 {

			tempRoom.MoveTo(moveTo)
			moveAttempts -= 1

			if !tempRoom.RoomInBounds(d) {
				continue //not inside the dungeon yet
			}
			if tempRoom.BadOverlap(d) {
				break //we've gone too far!
			}
			if d.AttemptStamp(tempRoom) {
				break //success!
			}
		}
	}

	//place spawn points
	for i := 0; i < 2; i++ {
		for {
			spawnX := RandInt(drand, 0, d.Width)
			spawnY := RandInt(drand, 0, d.Height)

			if d.GetBlock(spawnX, spawnY) == FLOOR {
				if i == 0 {
					d.PlaceBlock(REDSPAWN, spawnX, spawnY)
					break
				} else {
					d.PlaceBlock(BLUESPAWN, spawnX, spawnY)
					break
				}
			}
		}
	}
}

func (d *Dungeon) GetSpawnPoint(color Color) (int, int) {
	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			if int(d.Blocks[x][y]) == REDSPAWN && color == RED {
				return x, y
			} else if int(d.Blocks[x][y]) == BLUESPAWN && color == BLUE {
				return x, y
			}
		}
	}
	return -1, -1
}

func (d *Dungeon) Count(color Color) int {
	count := 0
	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			if int(d.Colors[x][y]) == int(color) {
				count++
			}
		}
	}
	return count
}

func (d *Dungeon) GetWinner() Color {
	redCount := d.Count(RED)
	blueCount := d.Count(BLUE)
	println("Red Count: ", redCount)
	println("Blue Count: ", blueCount)

	if redCount > blueCount {
		return RED
	} else if blueCount > redCount {
		return BLUE
	}
	return WHITE
}

func (d *Dungeon) AddPlayer(player *Player) {
	d.Players = append(d.Players, player)
}
