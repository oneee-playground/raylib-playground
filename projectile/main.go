package main

import (
	"math"
	"os"
	"path/filepath"
	"time"

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

	sideMoveCnt int
	latestShot  time.Time

	cam rl.Camera2D

	deltaX    float32
	deltaY    float32
	isOnFloor bool

	projectiles   [5]*projectile
	projectileIdx int
}

type projectile struct {
	positionRect rl.Rectangle

	deltaX   float32
	isActive bool
}

var (
	penguinSprite rl.Texture2D

	player = &playerObject{
		spriteRect:    rl.Rectangle{},
		spriteFaceDir: faceDown,
		spriteSeq:     0,
		positionRect:  rl.Rectangle{},
		sideMoveCnt:   0,
		cam:           rl.Camera2D{},
		isOnFloor:     false,
		deltaX:        0,
		deltaY:        0,
		projectiles:   [5]*projectile{{}, {}, {}, {}, {}},
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
	case rl.IsKeyDown(rl.KeySpace):
		now := time.Now()

		if now.Before(player.latestShot.Add(time.Second)) {
			break
		}
		player.latestShot = now

		posX := player.positionRect.X
		posY := player.positionRect.Y - (player.positionRect.Height / 2)

		deltaX := defaultPlayerSpeed * 2
		if player.spriteFaceDir == faceLeft {
			deltaX *= -1
			posX -= player.positionRect.Width
		}

		p := &projectile{
			positionRect: rl.NewRectangle(posX, posY, 8, 8),
			deltaX:       deltaX,
			isActive:     true,
		}

		player.projectiles[player.projectileIdx] = p
		player.projectileIdx = (player.projectileIdx + 1) & len(player.projectiles)
	case player.isOnFloor && rl.IsKeyDown(rl.KeyW):
		player.deltaY = -defaultPlayerSpeed * 2
		player.isOnFloor = false
	case rl.IsKeyDown(rl.KeyA):
		player.spriteFaceDir = faceLeft
		player.deltaX = -defaultPlayerSpeed
		player.sideMoveCnt++
	case rl.IsKeyDown(rl.KeyD):
		player.spriteFaceDir = faceRight
		player.deltaX = defaultPlayerSpeed
		player.sideMoveCnt++
	default:
		player.deltaX = 0
	}

	player.spriteRect.Y = player.spriteFaceDir * spriteHeight

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

	if player.sideMoveCnt == maxKeyDownCnt {
		player.spriteSeq = (player.spriteSeq + 1) & 3
		player.spriteRect.X = float32(player.spriteSeq) * spriteWidth
		player.sideMoveCnt = 0
	}

	player.cam.Target = rl.NewVector2(
		player.positionRect.X-(player.positionRect.Width/2),
		player.positionRect.Y-(player.positionRect.Height/2),
	)

	for _, p := range player.projectiles {
		if p.isActive {
			posX := p.positionRect.X + p.deltaX
			diff := math.Abs(float64(posX - player.positionRect.X))
			if diff > float64(windowWidth/2) {
				p.isActive = false
			}

			p.positionRect.X = posX
		}
	}
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

	for _, p := range player.projectiles {
		if p.isActive {
			rl.DrawRectangleRec(p.positionRect, rl.Red)
		}
	}
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
		rl.ClearBackground(rl.LightGray)

		drawScene()

		rl.EndMode2D()
		rl.EndDrawing()
	}

	unload()

	rl.CloseWindow()
}
