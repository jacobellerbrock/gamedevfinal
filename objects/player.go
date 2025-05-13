package objects

import (
	// rl "github.com/gen2brain/raylib-go/raylib"
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
}

func NewPlayer(x int, y int, color rl.Color) *Player {
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
	return &player
}

func (p *Player) Draw(d Dungeon) {
	rl.DrawRectangle(int32(p.Position.X*d.BlockSize), int32(p.Position.Y*d.BlockSize), int32(d.BlockSize), int32(d.BlockSize), p.Color)
}

func (p *Player) PaintFloor(d Dungeon) {
	d.Colors[p.Position.X][p.Position.Y] = p.ColorID
}

func (p *Player) UseAbility1(d Dungeon) {

}

func (p *Player) Move(direction int, dungeon Dungeon, otherPlayer *Player) {
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
