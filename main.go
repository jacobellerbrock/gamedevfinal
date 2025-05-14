package main

import (
	. "final/project/objects"
	"math"
	"math/rand"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	PLAYING GameState = iota
	OVER
	MENU
	PAUSED
)

const (
	TIMER = 60
)

func main() {
	rl.InitWindow(1920, 1080, "Final Project")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	music := rl.LoadMusicStream("resources/audio/music/Boost.mp3")
	defer rl.UnloadMusicStream(music)

	rl.PlayMusicStream(music)
	rl.SetMusicVolume(music, 0.3) // Set volume to 30%

	walkSound := rl.LoadSound("resources/audio/sfx/walk.wav")
	rl.SetSoundVolume(walkSound, 0.3)
	defer rl.UnloadSound(walkSound)

	placeSound := rl.LoadSound("resources/audio/sfx/place.wav")
	rl.SetSoundVolume(placeSound, 0.3)
	defer rl.UnloadSound(placeSound)

	explodeSound := rl.LoadSound("resources/audio/sfx/explode.wav")
	rl.SetSoundVolume(explodeSound, 0.3)
	defer rl.UnloadSound(explodeSound)

	deathSound := rl.LoadSound("resources/audio/sfx/death.wav")
	rl.SetSoundVolume(deathSound, 0.3)
	defer rl.UnloadSound(deathSound)

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
		"pause":           rl.KeyP,
		"quit":            rl.KeyEscape,
		"menu":            rl.KeyM,
		"cheatPlayer2Die": rl.KeyF1,
	}

	deadSprite := rl.LoadTexture("resources/sprites/deadplayer.png")

	// load two random seeds based on time
	rand.Seed(time.Now().UnixNano())

	dungeon := NewDungeon(60, 33)

	camera := rl.NewCamera2D(rl.Vector2Zero(), rl.Vector2Zero(), 0, 2)

	var player1, player2 *Player

	bombs := make([]Bomb, 0)
	gameState := MENU

	var timer float32 = TIMER

	for !rl.WindowShouldClose() {
		if gameState == PLAYING {
			if timer > 1 {
				timer -= rl.GetFrameTime()
			} else {
				// TODO: timer end
				gameState = OVER
			}
		}

		rl.UpdateMusicStream(music) // Update music buffer

		ManageInput(Keymaps, dungeon, player1, player2, &gameState, &bombs, walkSound, placeSound, explodeSound, deathSound)

		rl.BeginDrawing()
		rl.BeginMode2D(camera)
		rl.ClearBackground(rl.RayWhite)

		if gameState == PLAYING {
			dungeon.DrawDungeon()

			for _, player := range dungeon.Players {
				player.Update(dungeon, deadSprite, (*int)(&gameState), deathSound)
			}

			for _, bomb := range bombs {
				bomb.Update(&bombs, dungeon, (*int)(&gameState), explodeSound)
			}
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

			if rl.IsKeyPressed(rl.KeyR) {
				player1, player2 = RestartGame(&gameState, &dungeon, player1, player2, &bombs)
				gameState = PLAYING
			}

		} else if gameState == MENU {
			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.RayWhite)
			rl.DrawText("Press Enter to start", int32(rl.GetScreenWidth()/2-100), 100, 50, rl.Black)
			rl.DrawText("Player 1 Controls:", int32(rl.GetScreenWidth()/2-100), 150, 50, rl.Black)
			rl.DrawText("WASD to move", int32(rl.GetScreenWidth()/2-100), 200, 50, rl.Black)
			rl.DrawText("Q to use ability", int32(rl.GetScreenWidth()/2-100), 250, 50, rl.Black)
			rl.DrawText("Player 2 Controls:", int32(rl.GetScreenWidth()/2-100), 300, 50, rl.Black)
			rl.DrawText("Arrow keys to move", int32(rl.GetScreenWidth()/2-100), 350, 50, rl.Black)
			rl.DrawText("Right Control to use ability", int32(rl.GetScreenWidth()/2-100), 400, 50, rl.Black)
			rl.DrawText("P to Pause", int32(rl.GetScreenWidth()/2-100), 450, 50, rl.Black)
			rl.DrawText("Esc to Quit", int32(rl.GetScreenWidth()/2-100), 500, 50, rl.Black)
			if rl.IsKeyPressed(rl.KeyEnter) {
				player1, player2 = RestartGame(&gameState, &dungeon, player1, player2, &bombs)
				gameState = PLAYING
			}
		} else if gameState == PAUSED {
			rl.DrawText("Game Paused", int32(rl.GetScreenWidth()/2-100), 100, 50, rl.Black)
			rl.DrawText("P to Resume", int32(rl.GetScreenWidth()/2-100), 150, 50, rl.Black)
			rl.DrawText("M to Return to Menu", int32(rl.GetScreenWidth()/2-100), 200, 50, rl.Black)
			rl.DrawText("Esc to Quit", int32(rl.GetScreenWidth()/2-100), 250, 50, rl.Black)
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

func ManageInput(keys map[string]int32, dungeon Dungeon, player1 *Player, player2 *Player, gameState *GameState, bombs *[]Bomb, walkSound rl.Sound, placeSound rl.Sound, explodeSound rl.Sound, deathSound rl.Sound) {
	if *gameState == PLAYING {
		moved := false
		if rl.IsKeyPressed(keys["player1Up"]) {
			player1.Move(UP, dungeon, player2)
			moved = true
		} else if rl.IsKeyPressed(keys["player1Down"]) {
			player1.Move(DOWN, dungeon, player2)
			moved = true
		} else if rl.IsKeyPressed(keys["player1Left"]) {
			player1.Move(LEFT, dungeon, player2)
			moved = true
		} else if rl.IsKeyPressed(keys["player1Right"]) {
			player1.Move(RIGHT, dungeon, player2)
			moved = true
		} else if rl.IsKeyPressed(keys["player1Ability1"]) {
			player1.UseAbility1(dungeon, bombs, (*int)(gameState))
			rl.PlaySound(placeSound)
		} else if rl.IsKeyPressed(keys["player2Up"]) {
			player2.Move(UP, dungeon, player1)
			moved = true
		} else if rl.IsKeyPressed(keys["player2Down"]) {
			player2.Move(DOWN, dungeon, player1)
			moved = true
		} else if rl.IsKeyPressed(keys["player2Left"]) {
			player2.Move(LEFT, dungeon, player1)
			moved = true
		} else if rl.IsKeyPressed(keys["player2Right"]) {
			player2.Move(RIGHT, dungeon, player1)
			moved = true
		} else if rl.IsKeyPressed(keys["player2Ability1"]) {
			player2.UseAbility1(dungeon, bombs, (*int)(gameState))
			rl.PlaySound(placeSound)
		} else if rl.IsKeyPressed(keys["pause"]) {
			*gameState = PAUSED
		}

		if moved {
			// Random pitch between 0.95 and 1.05
			randomPitch := 0.95 + (rand.Float32() * 0.1)
			rl.SetSoundPitch(walkSound, randomPitch)
			rl.PlaySound(walkSound)
		}

		if rl.IsKeyPressed(keys["cheatPlayer2Die"]) {
			player2.Health = 0
		}
	} else if *gameState == PAUSED {
		if rl.IsKeyPressed(keys["pause"]) {
			*gameState = PLAYING
		} else if rl.IsKeyPressed(keys["quit"]) {
			rl.CloseWindow()
		} else if rl.IsKeyPressed(keys["menu"]) {
			*gameState = MENU
		}
	}
}

func RestartGame(gameState *GameState, dungeon *Dungeon, player1 *Player, player2 *Player, bombs *[]Bomb) (*Player, *Player) {
	seed1 := rand.Uint64()
	seed2 := rand.Uint64()
	*gameState = PLAYING
	dungeon.Generate(seed1, seed2)
	x, y := dungeon.GetSpawnPoint(RED)
	newPlayer1 := NewPlayer(x, y, rl.Red, 100)
	x, y = dungeon.GetSpawnPoint(BLUE)
	newPlayer2 := NewPlayer(x, y, rl.Blue, 100)

	dungeon.Players = make([]*Player, 0)
	dungeon.AddPlayer(newPlayer1)
	dungeon.AddPlayer(newPlayer2)

	*bombs = make([]Bomb, 0)
	return newPlayer1, newPlayer2
}
