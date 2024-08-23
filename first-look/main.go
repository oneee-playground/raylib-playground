package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	rl.InitWindow(300, 300, "Hello, World!")
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Hello, World!", 0, 0, 20, rl.LightGray)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
