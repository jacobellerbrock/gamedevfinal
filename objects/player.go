package objects

import (
	// rl "github.com/gen2brain/raylib-go/raylib"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Position struct {
	X int
	Y int
}

const (
	UP    = 0
	DOWN  = 1
	LEFT  = 2
	RIGHT = 3
)

type Player struct {
	Position  Position
	Speed     int
	Direction int // const
	Color     rl.Color
	ColorID   Color
	Health    float32
	MaxHealth float32
	IsDead    bool
}

func NewPlayer(x int, y int, color rl.Color, maxHealth float32) *Player {
	player := Player{}
	player.Position = Position{X: x, Y: y}
	player.Speed = 1
	player.Direction = DOWN
	player.Color = color
	if color == rl.Red {
		player.ColorID = 1
	} else if color == rl.Blue {
		player.ColorID = 2
	} else {
		player.ColorID = 0
	}
	player.Health = maxHealth
	player.MaxHealth = maxHealth
	return &player
}

func (p *Player) Update(dungeon Dungeon, deadSprite rl.Texture2D) {
	if p.Health <= 0 {
		if !p.IsDead {
			go p.Dead(dungeon)
		}
	} else {
		p.PaintFloor(dungeon)
	}
	p.Draw(dungeon, deadSprite)
}

func (p *Player) Draw(d Dungeon, deadSprite rl.Texture2D) {
	if !p.IsDead {
		rl.DrawRectangle(int32(p.Position.X*d.BlockSize), int32(p.Position.Y*d.BlockSize), int32(d.BlockSize), int32(d.BlockSize), p.Color)
	} else {
		rl.DrawTexture(deadSprite, int32(p.Position.X*d.BlockSize), int32(p.Position.Y*d.BlockSize), p.Color)
	}
}

func (p *Player) PaintFloor(d Dungeon) {
	d.Colors[p.Position.X][p.Position.Y] = p.ColorID
}

func (p *Player) UseAbility1(d Dungeon, bombs *[]Bomb) {
	*bombs = append(*bombs, NewBomb(d.GetBlockPosition(p.Position.X, p.Position.Y), 3, 3, p.ColorID))
}

func (p *Player) Move(direction int, dungeon Dungeon, otherPlayer *Player) {
	if p.IsDead {
		return
	}
	// move direction, check for collision with walls and other player
	switch direction {
	case UP:
		if dungeon.Blocks[p.Position.X][p.Position.Y-1] != WALL && (p.Position.X != otherPlayer.Position.X || p.Position.Y-1 != otherPlayer.Position.Y) {
			p.Position.Y--
		}
	case DOWN:
		if dungeon.Blocks[p.Position.X][p.Position.Y+1] != WALL && (p.Position.X != otherPlayer.Position.X || p.Position.Y+1 != otherPlayer.Position.Y) {
			p.Position.Y++
		}
	case LEFT:
		if dungeon.Blocks[p.Position.X-1][p.Position.Y] != WALL && (p.Position.X-1 != otherPlayer.Position.X || p.Position.Y != otherPlayer.Position.Y) {
			p.Position.X--
		}
	case RIGHT:
		if dungeon.Blocks[p.Position.X+1][p.Position.Y] != WALL && (p.Position.X+1 != otherPlayer.Position.X || p.Position.Y != otherPlayer.Position.Y) {
			p.Position.X++
		}
	}
}

func (p *Player) Dead(dungeon Dungeon) {
	p.IsDead = true

	deathTimer := NewTimer(2, true)
	for !deathTimer.IsDone() {
		time.Sleep(1 * time.Millisecond)
	}
	// Get spawn coordinates
	spawnX, spawnY := dungeon.GetSpawnPoint(p.ColorID)
	println(spawnX, spawnY)

	p.IsDead = false
	p.Health = p.MaxHealth
	// Move player back to spawn
	p.Position.X = spawnX
	p.Position.Y = spawnY
}
