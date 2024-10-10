package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/rand"
)

const (
	screenWidth  = 800
	screenHeight = 480
	fontSize     = 36
	playerSpeed  = 3
)

const (
	Down PlayerDirection = iota
	Up
	Left
	Right
)

var (
	running         = true
	backgroungColor = rl.DarkGreen

	grassSprite rl.Texture2D

	playerSprite rl.Texture2D
	playerSrc    rl.Rectangle
	playerDest   rl.Rectangle

	playerMoving bool
	playerDir    PlayerDirection
	playerFrame  int
	frameCount   int

	musicPaused = false
	music       rl.Music

	tileDest = rl.NewRectangle(0, 0, 16, 16)
	tileSrc  = rl.NewRectangle(0, 0, 16, 16)

	tileMap []int
	mapW    = 10
	mapH    = 10

	cam rl.Camera2D
)

type PlayerDirection int

func loadMap() {
	mapW, mapH = 10, 10
	tileMap = make([]int, mapW*mapH)
	for i := 0; i < len(tileMap); i++ {
		tileMap[i] = rand.Intn(80)
	}
}

func init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Forest game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")

	playerSprite = rl.LoadTexture("assets/Characters/Spritesheet.png")
	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 150, 150)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/music/ForestWalk.mp3")
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(screenWidth/2.0, screenHeight/2.0),
		rl.NewVector2(playerDest.X-playerDest.Width/2,
			playerDest.Y-playerDest.Height/2), 0.0, 1.0)

	loadMap()
}

func update() {
	running = !rl.WindowShouldClose()

	if playerMoving {
		if playerDir == Up {
			playerDest.Y -= playerSpeed
		}
		if playerDir == Down {
			playerDest.Y += playerSpeed
		}
		if playerDir == Left {
			playerDest.X -= playerSpeed
		}
		if playerDir == Right {
			playerDest.X += playerSpeed
		}
		if frameCount%6 == 1 {
			playerFrame++
		}
	}

	// Update the player frame even when not moving
	if frameCount%30 == 1 {
		playerFrame++
	}

	// Switch between frame 0 and 1 when player is not moving
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(playerDest.X-playerDest.Width/2,
		playerDest.Y-playerDest.Height/2)

	playerMoving = false
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = Up
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = Down
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = Left
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = Right
	}
	if rl.IsKeyDown(rl.KeyM) {
		musicPaused = !musicPaused
	}

	if rl.IsKeyDown(rl.KeyC) {
		cam.Zoom = cam.Zoom + 0.1
	}
	if rl.IsKeyDown(rl.KeyV) {
		cam.Zoom = cam.Zoom - 0.1
	}
}

func quit() {
	rl.UnloadMusicStream(music)
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.CloseWindow()
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(backgroungColor)
	drawScene()
	rl.EndDrawing()
}

func drawScene() {

	for i := 0; i < len(tileMap); i++ {
		tileDest.X = tileDest.Width * float32(i%mapW)
		tileDest.Y = tileDest.Height * float32(i/mapW)
		tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(grassSprite.
			Width/int32(tileSrc.Width)))
		tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(grassSprite.
			Height/int32(tileSrc.Height)))
		rl.DrawTexturePro(grassSprite, tileSrc, tileDest,
			rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
	}
	// Drawing the player same as before
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest,
		rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func main() {

	for running {
		input()
		update()
		render()
	}

	quit()
}
