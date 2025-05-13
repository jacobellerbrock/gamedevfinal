package main

import (
	. "final/project/objects"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"strconv"
)

type GameState int

const (
	PLAYING GameState = iota
	OVER
)

const (
	TIMER = 60
)

func main() {
	rl.InitWindow(1920, 1080, "Final Project")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	Keymaps := map[string]int32{
		"player1Left":     rl.KeyA,
		"player1Right":    rl.KeyD,
		"player1Up":       rl.KeyW,
		"player1Down":     rl.KeyS,
		"player1Ability1": rl.KeyQ,
		"player2Left":     rl.KeyLeft,
		"player2Right":    rl.KeyRight,
		"player2Up":       rl.KeyUp,
		"player2Down":     rl.KeyDown,
	}

	dungeon := NewDungeon(60, 33)
	dungeon.Generate(uint64(420), uint64(69))

	camera := rl.NewCamera2D(rl.Vector2Zero(), rl.Vector2Zero(), 0, 2)

	x, y := dungeon.GetSpawnPoint(REDSPAWN)
	player1 := NewPlayer(x, y, rl.Red)
	x, y = dungeon.GetSpawnPoint(BLUESPAWN)
	player2 := NewPlayer(x, y, rl.Blue)

	var timer float32 = TIMER
	gameState := PLAYING

	for !rl.WindowShouldClose() {

		if timer > 1 {
			timer -= rl.GetFrameTime()
		} else {
			// TODO: timer end
			gameState = OVER
		}

		ManageInput(Keymaps, dungeon, player1, player2, gameState)

		rl.BeginDrawing()
		rl.BeginMode2D(camera)
		rl.ClearBackground(rl.RayWhite)

		dungeon.DrawDungeon()

		player1.PaintFloor(dungeon)
		player2.PaintFloor(dungeon)
		player1.Draw(dungeon)
		player2.Draw(dungeon)

		rl.EndMode2D()

		if gameState == OVER {
			rl.DrawText("Game Over", int32(rl.GetScreenWidth()/2)-100, int32(rl.GetScreenHeight()/2)-50, 40, rl.Red)
			rl.DrawText("Press R to restart", int32(rl.GetScreenWidth()/2)-100, int32(rl.GetScreenHeight()/2)+50, 20, rl.Red)
			winner := dungeon.GetWinner()
			var winnerText string
			if winner == RED {
				winnerText = "Red wins!"
			} else if winner == BLUE {
				winnerText = "Blue wins!"
			} else {
				winnerText = "It's a draw!"
			}
			rl.DrawText(winnerText, int32(rl.GetScreenWidth()/2)-100, int32(rl.GetScreenHeight()/2)+100, 50, rl.Black)
		} else {
			DrawText(strconv.Itoa(int(math.Ceil(float64(timer)))), 50, rl.Black)
		}

		rl.EndDrawing()
	}
}

func DrawText(text string, size int32, color rl.Color) {
	textWidth := rl.MeasureText(text, size)
	rl.DrawText(text, int32(rl.GetScreenWidth()/2)-(textWidth/2), 10, size, color)
}

func ManageInput(keys map[string]int32, dungeon Dungeon, player1 *Player, player2 *Player, gameState GameState) {
	if gameState == PLAYING {
		if rl.IsKeyPressed(keys["player1Up"]) {
			player1.Move(UP, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player1Down"]) {
			player1.Move(DOWN, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player1Left"]) {
			player1.Move(LEFT, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player1Right"]) {
			player1.Move(RIGHT, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player2Up"]) {
			player2.Move(UP, dungeon, player1)
		} else if rl.IsKeyPressed(keys["player2Down"]) {
			player2.Move(DOWN, dungeon, player1)
		} else if rl.IsKeyPressed(keys["player2Left"]) {
			player2.Move(LEFT, dungeon, player1)
		} else if rl.IsKeyPressed(keys["player2Right"]) {
			player2.Move(RIGHT, dungeon, player1)
		}
	}
}
