package main

import (
	. "final/project/objects"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
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
		"player2Ability1": rl.KeyRightControl,
		"cheatPlayer2Die": rl.KeyF1,
	}

	deadSprite := rl.LoadTexture("resources/sprites/deadplayer.png")

	dungeon := NewDungeon(60, 33)
	dungeon.Generate(uint64(420), uint64(69))

	camera := rl.NewCamera2D(rl.Vector2Zero(), rl.Vector2Zero(), 0, 2)

	x, y := dungeon.GetSpawnPoint(RED)
	player1 := NewPlayer(x, y, rl.Red, 100)
	x, y = dungeon.GetSpawnPoint(BLUE)
	player2 := NewPlayer(x, y, rl.Blue, 100)

	dungeon.AddPlayer(player1)
	dungeon.AddPlayer(player2)

	bombs := make([]Bomb, 0)

	var timer float32 = TIMER
	gameState := PLAYING

	for !rl.WindowShouldClose() {

		if timer > 1 {
			timer -= rl.GetFrameTime()
		} else {
			// TODO: timer end
			gameState = OVER
		}

		ManageInput(Keymaps, dungeon, player1, player2, gameState, &bombs)

		rl.BeginDrawing()
		rl.BeginMode2D(camera)
		rl.ClearBackground(rl.RayWhite)

		dungeon.DrawDungeon()

		for _, player := range dungeon.Players {
			player.Update(dungeon, deadSprite)
		}

		for _, bomb := range bombs {
			bomb.Update(&bombs, dungeon)
			bomb.Draw()
		}

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

func ManageInput(keys map[string]int32, dungeon Dungeon, player1 *Player, player2 *Player, gameState GameState, bombs *[]Bomb) {
	if gameState == PLAYING {
		if rl.IsKeyPressed(keys["player1Up"]) {
			player1.Move(UP, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player1Down"]) {
			player1.Move(DOWN, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player1Left"]) {
			player1.Move(LEFT, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player1Right"]) {
			player1.Move(RIGHT, dungeon, player2)
		} else if rl.IsKeyPressed(keys["player1Ability1"]) {
			player1.UseAbility1(dungeon, bombs)
		} else if rl.IsKeyPressed(keys["player2Up"]) {
			player2.Move(UP, dungeon, player1)
		} else if rl.IsKeyPressed(keys["player2Down"]) {
			player2.Move(DOWN, dungeon, player1)
		} else if rl.IsKeyPressed(keys["player2Left"]) {
			player2.Move(LEFT, dungeon, player1)
		} else if rl.IsKeyPressed(keys["player2Right"]) {
			player2.Move(RIGHT, dungeon, player1)
		} else if rl.IsKeyPressed(keys["player2Ability1"]) {
			player2.UseAbility1(dungeon, bombs)
		}

		if rl.IsKeyPressed(keys["cheatPlayer2Die"]) {
			player2.Health = 0
		}
	}
}
