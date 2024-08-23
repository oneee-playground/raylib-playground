package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	penguinSprite rl.Texture2D

	spriteRect rl.Rectangle
	playerRec  rl.Rectangle

	spriteFaceDir = faceDown
	spriteSeq     = 0
	keyDownCnt    = 0
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
	penguinSprite = rl.LoadTexture("sprites/penguin.png")
	spriteRect = rl.NewRectangle(0, spriteFaceDir*spriteHeight, 24, 32)
	playerRec = rl.NewRectangle(300, 300, 48, 64)
}

func unload() {
	rl.UnloadTexture(penguinSprite)
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

}

func drawScene() {
	rl.DrawTexturePro(penguinSprite, spriteRect, playerRec, rl.NewVector2(playerRec.Width, playerRec.Height), 0, rl.White)
}

func main() {
	rl.InitWindow(600, 600, "Hello, World!")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	load()

	for !rl.WindowShouldClose() {
		getInput()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		drawScene()

		rl.EndDrawing()
	}

	unload()

	rl.CloseWindow()
}
