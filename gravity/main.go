package main

import (
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	faceUp float32 = iota
	faceRight
	faceDown
	faceLeft
)

const (
	windowHeight, windowWidth = 600, 600

	spriteHeight float32 = 32
	spriteWidth  float32 = 24

	defaultPlayerSpeed float32 = 3
	maxKeyDownCnt      int     = 10

	floorHeight float32 = windowHeight * 0.6

	gravity float32 = 0.2
)

type playerObject struct {
	spriteRect    rl.Rectangle
	spriteFaceDir float32
	spriteSeq     int

	positionRect rl.Rectangle

	keyDownCnt int

	cam rl.Camera2D

	deltaX    float32
	deltaY    float32
	isOnFloor bool
}

var (
	penguinSprite rl.Texture2D

	player = &playerObject{
		spriteRect:    rl.Rectangle{},
		spriteFaceDir: faceDown,
		spriteSeq:     0,
		positionRect:  rl.Rectangle{},
		keyDownCnt:    0,
		cam:           rl.Camera2D{},
		isOnFloor:     false,
		deltaX:        0,
		deltaY:        0,
	}
)

func load() {
	baseDir, err := os.Executable()
	if err != nil {
		panic("what happened? " + err.Error())
	}
	baseDir = filepath.Dir(filepath.Dir(baseDir))

	penguinSprite = rl.LoadTexture(filepath.Join(baseDir, "sprites/penguin.png"))

	player.spriteRect = rl.NewRectangle(0, player.spriteFaceDir*spriteHeight, 24, 32)
	player.positionRect = rl.NewRectangle(300, 300, 24, 32)

	player.cam = rl.NewCamera2D(rl.NewVector2(300, 300), rl.NewVector2(
		player.positionRect.X-(player.positionRect.Width/2),
		player.positionRect.Y-(player.positionRect.Height/2)),
		0, 1.0)
}

func unload() {
	rl.UnloadTexture(penguinSprite)
}

func getInputAndAct() {
	switch {
	case player.isOnFloor && rl.IsKeyDown(rl.KeyW):
		// player.spriteRect.Y = faceUp * spriteHeight
		player.deltaY = -defaultPlayerSpeed * 2
		player.isOnFloor = false
	// case rl.IsKeyDown(rl.KeyS):
	// 	player.spriteRect.Y = faceDown * spriteHeight
	// 	player.positionRect.Y += playerSpeed
	case rl.IsKeyDown(rl.KeyA):
		player.spriteRect.Y = faceLeft * spriteHeight
		player.deltaX = -defaultPlayerSpeed
		player.keyDownCnt++
	case rl.IsKeyDown(rl.KeyD):
		player.spriteRect.Y = faceRight * spriteHeight
		player.deltaX = defaultPlayerSpeed
		player.keyDownCnt++
	default:
		player.deltaX = 0
	}

	if player.deltaX != 0 {
		player.positionRect.X += player.deltaX
	}

	if !player.isOnFloor {
		player.positionRect.Y = min(floorHeight, player.positionRect.Y+player.deltaY)
		if player.positionRect.Y == floorHeight {
			player.deltaY = 0
			player.isOnFloor = true
		}

		player.deltaY += gravity
	}

	if player.keyDownCnt == maxKeyDownCnt {
		player.spriteSeq = (player.spriteSeq + 1) & 3
		player.spriteRect.X = float32(player.spriteSeq) * spriteWidth
		player.keyDownCnt = 0
	}

	player.cam.Target = rl.NewVector2(
		player.positionRect.X-(player.positionRect.Width/2),
		player.positionRect.Y-(player.positionRect.Height/2),
	)
}

const spawnText = "This is where you spawn"
const fontSize = 20

func drawScene() {
	rl.DrawText(spawnText, 300-rl.MeasureText(spawnText, fontSize)/2, 300, fontSize, rl.Gray)

	rl.DrawTexturePro(penguinSprite, player.spriteRect, player.positionRect,
		rl.NewVector2(
			player.positionRect.Width,
			player.positionRect.Height,
		),
		0, rl.White)
}

func main() {
	rl.InitWindow(windowWidth, windowHeight, "Hello, World!")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	load()

	for !rl.WindowShouldClose() {
		getInputAndAct()

		rl.BeginDrawing()
		rl.BeginMode2D(player.cam)
		rl.ClearBackground(rl.RayWhite)

		drawScene()

		rl.EndMode2D()
		rl.EndDrawing()
	}

	unload()

	rl.CloseWindow()
}
