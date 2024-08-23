package main

import (
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	penguinSprite rl.Texture2D

	spriteRect rl.Rectangle
	playerRec  rl.Rectangle

	spriteFaceDir = faceDown
	spriteSeq     = 0
	keyDownCnt    = 0

	music rl.Music
	cam   rl.Camera2D
)

const (
	faceUp float32 = iota
	faceRight
	faceDown
	faceLeft
)

const (
	spriteHeight float32 = 32
	spriteWidth  float32 = 24

	playerSpeed   float32 = 3
	maxKeyDownCnt int     = 10
)

func load() {
	baseDir, err := os.Executable()
	if err != nil {
		panic("what happened? " + err.Error())
	}
	baseDir = filepath.Dir(baseDir)

	penguinSprite = rl.LoadTexture(filepath.Join(baseDir, "sprites/penguin.png"))

	music = rl.LoadMusicStream(filepath.Join(baseDir, "musics/some-random-music.mp3"))
	rl.SetMusicVolume(music, 1.0)
	rl.PlayMusicStream(music)

	spriteRect = rl.NewRectangle(0, spriteFaceDir*spriteHeight, 24, 32)
	playerRec = rl.NewRectangle(300, 300, 48, 64)

	cam = rl.NewCamera2D(rl.NewVector2(300, 300), rl.NewVector2(playerRec.X-(playerRec.Width/2), playerRec.Y-(playerRec.Height/2)), 0, 1.0)
}

func unload() {
	rl.UnloadTexture(penguinSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
}

func getInput() {
	keyDownCnt++

	switch {
	case rl.IsKeyDown(rl.KeyW):
		spriteRect.Y = faceUp * spriteHeight
		playerRec.Y -= playerSpeed
	case rl.IsKeyDown(rl.KeyA):
		spriteRect.Y = faceLeft * spriteHeight
		playerRec.X -= playerSpeed
	case rl.IsKeyDown(rl.KeyS):
		spriteRect.Y = faceDown * spriteHeight
		playerRec.Y += playerSpeed
	case rl.IsKeyDown(rl.KeyD):
		spriteRect.Y = faceRight * spriteHeight
		playerRec.X += playerSpeed
	default:
		keyDownCnt--
	}

	if keyDownCnt == maxKeyDownCnt {
		spriteSeq = (spriteSeq + 1) & 3
		spriteRect.X = float32(spriteSeq) * spriteWidth
		keyDownCnt = 0
	}

	cam.Target = rl.NewVector2(playerRec.X-(playerRec.Width/2), playerRec.Y-(playerRec.Height/2))
}

const spawnText = "This is where you spawn"
const fontSize = 20

func drawScene() {
	rl.DrawText(spawnText, 300-rl.MeasureText(spawnText, fontSize)/2, 300, fontSize, rl.Gray)

	rl.DrawTexturePro(penguinSprite, spriteRect, playerRec, rl.NewVector2(playerRec.Width, playerRec.Height), 0, rl.White)
}

func main() {
	rl.InitWindow(600, 600, "Hello, World!")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)
	rl.InitAudioDevice()

	load()

	for !rl.WindowShouldClose() {
		getInput()

		rl.BeginDrawing()
		rl.BeginMode2D(cam)
		rl.ClearBackground(rl.RayWhite)

		rl.UpdateMusicStream(music)

		drawScene()

		rl.EndMode2D()
		rl.EndDrawing()
	}

	unload()

	rl.CloseWindow()
}
