package objects

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/slices"
)

type Bomb struct {
	Position    rl.Vector2
	Radius      int
	Timer       *Timer
	PlayerColor Color
	Color       rl.Color
}

func NewBomb(position rl.Vector2, radius int, timer float32, playerColor Color) Bomb {
	var bombColor rl.Color
	if playerColor == RED {
		bombColor = rl.NewColor(120, 20, 0, 255)
	} else if playerColor == BLUE {
		bombColor = rl.NewColor(0, 64, 127, 255)
	}
	return Bomb{
		Position:    position,
		Radius:      radius,
		Timer:       NewTimer(timer, false),
		PlayerColor: playerColor,
		Color:       bombColor,
	}
}

func (b *Bomb) Update(bombs *[]Bomb, dungeon Dungeon) {
	b.Timer.Update()
	if b.Timer.IsDone() {
		b.Explode(bombs, dungeon)
	}
}

func (b *Bomb) Draw() {
	rl.DrawCircle(int32(b.Position.X), int32(b.Position.Y), float32(b.Radius), b.Color)
}

func (b *Bomb) Explode(bombs *[]Bomb, dungeon Dungeon) {

	// Convert bomb position to grid coordinates
	gridX := int(b.Position.X) / dungeon.BlockSize
	gridY := int(b.Position.Y) / dungeon.BlockSize

	// Check all tiles in a square around the bomb up to radius
	for x := gridX - b.Radius; x <= gridX+b.Radius; x++ {
		for y := gridY - b.Radius; y <= gridY+b.Radius; y++ {
			// Skip if out of bounds
			if x < 0 || y < 0 || x >= dungeon.Width || y >= dungeon.Height {
				continue
			}

			// Skip if not a floor tile
			if dungeon.Blocks[x][y] != FLOOR {
				continue
			}

			// Calculate distance from bomb
			dx := x - gridX
			dy := y - gridY
			distSq := dx*dx + dy*dy

			// Skip if outside radius
			if distSq > b.Radius*b.Radius {
				continue
			}

			// Ray trace to check line of sight
			rayX := float32(gridX)
			rayY := float32(gridY)
			stepX := float32(dx) / float32(b.Radius)
			stepY := float32(dy) / float32(b.Radius)

			blocked := false
			for step := 0; step < b.Radius; step++ {
				checkX := int(rayX + stepX*float32(step))
				checkY := int(rayY + stepY*float32(step))

				if dungeon.Blocks[checkX][checkY] == WALL {
					blocked = true
					break
				}
			}

			if !blocked {
				dungeon.Colors[x][y] = b.PlayerColor

				// Check if any player is in this tile and damage them if they're the enemy
				for _, player := range dungeon.Players {
					if player.Position.X == x && player.Position.Y == y && player.ColorID != b.PlayerColor {
						player.Health -= 100
					}
				}
			}
		}
	}

	for i, bomb := range *bombs {
		if bomb.Position.X == b.Position.X && bomb.Position.Y == b.Position.Y {
			*bombs = slices.Delete(*bombs, i, i+1)
			return
		}
	}
}
