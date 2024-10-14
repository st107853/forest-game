package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
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

	grassSprite  rl.Texture2D
	fencedSprite rl.Texture2D
	hillSprite   rl.Texture2D
	houseSprite  rl.Texture2D
	tilledSprite rl.Texture2D
	waterSprite  rl.Texture2D

	tex rl.Texture2D

	playerSprite rl.Texture2D
	playerSrc    rl.Rectangle
	playerDest   rl.Rectangle

	playerMoving bool
	playerDir    PlayerDirection
	playerFrame  int
	frameCount   int

	musicPaused = false
	music       rl.Music

	tileDest   rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	cam rl.Camera2D
)

type PlayerDirection int

func loadMap(mapFile string) {
	fmt.Printf("Loading map: %s\n", mapFile)
	file, err := os.ReadFile(mapFile)
	if err != nil {
		fmt.Printf("Error reading map file: %s: %s\n", mapFile, err)
		os.Exit(1)
	}
	sliced := strings.Split(strings.ReplaceAll(string(file), "\n", " "), " ")
	mapW, _ = strconv.Atoi(sliced[0])
	mapH, _ = strconv.Atoi(sliced[1])
	tileMap = make([]int, mapW*mapH)
	srcMap = make([]string, mapW*mapH)
	for i := 0; i < mapW*mapH; i++ {
		m, _ := strconv.Atoi(sliced[i+2])
		tileMap[i] = m
	}
	for i := 0; i < mapW*mapH; i++ {
		srcMap[i] = sliced[i+2+mapW*mapH]
	}
}

func Init() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Forest game")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("assets/Tilesets/Grass.png")
	fencedSprite = rl.LoadTexture("assets/Tilesets/Fences.png")
	hillSprite = rl.LoadTexture("assets/Tilesets/Hills.png")
	houseSprite = rl.LoadTexture("assets/Tilesets/House.png")
	tilledSprite = rl.LoadTexture("assets/Tilesets/Tilled.png")
	waterSprite = rl.LoadTexture("assets/Tilesets/Water.png")

	tileSrc = rl.NewRectangle(0, 0, 16, 16)
	tileDest = rl.NewRectangle(0, 0, 48, 48)

	playerSprite = rl.LoadTexture("assets/Characters/Spritesheet.png")
	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 150, 150)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/music/ForestWalk.mp3")
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(screenWidth/2.0, screenHeight/2.0),
		rl.NewVector2(playerDest.X-playerDest.Width/2, playerDest.Y-playerDest.
			Height/2), 0.0, 1)

	loadMap("world.txt")
	rl.BeginMode2D(cam)
}

func update(velocity rl.Vector2) {
	running = !rl.WindowShouldClose()

	if playerMoving {
		playerDest.X += rl.Vector2Normalize(velocity).X * playerSpeed
		playerDest.Y += rl.Vector2Normalize(velocity).Y * playerSpeed
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
		playerFrame = 2
	}
	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(playerDest.X, playerDest.Y)
	playerMoving = false
}

func input() rl.Vector2 {
	velocity := rl.NewVector2(0, 0)

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = Up
		velocity = rl.Vector2Add(velocity, rl.NewVector2(0, -1))
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = Down
		velocity = rl.Vector2Add(velocity, rl.NewVector2(0, 1))
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = Left
		velocity = rl.Vector2Add(velocity, rl.NewVector2(-1, 0))
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = Right
		velocity = rl.Vector2Add(velocity, rl.NewVector2(1, 0))
	}
	if rl.IsKeyDown(rl.KeyM) {
		musicPaused = !musicPaused
	}

	if rl.IsKeyDown(rl.KeyC) {
		if cam.Zoom < 2.0 {
			cam.Zoom = cam.Zoom + 0.1
		}
	}
	if rl.IsKeyDown(rl.KeyV) {
		if cam.Zoom > 0.6 {
			cam.Zoom = cam.Zoom - 0.1
		}
	}
	return velocity
}

func quit() {
	rl.UnloadMusicStream(music)
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.CloseWindow()
}

func render() {
	rl.BeginDrawing()
	rl.BeginMode2D(cam)
	rl.ClearBackground(backgroungColor)
	drawScene()
	rl.EndDrawing()
	rl.EndMode2D()
}

func drawScene() {

	for i := 0; i < len(tileMap); i++ {
		tileDest.X = tileDest.Width * float32(i%mapW)
		tileDest.Y = tileDest.Height * float32(i/mapW)

		switch srcMap[i] {
		case "g":
			tex = grassSprite
		case "l":
			tex = hillSprite
		case "f":
			tex = fencedSprite
		case "h":
			tex = houseSprite
		case "w":
			tex = waterSprite
		case "t":
			tex = tilledSprite
		default:
			tex = grassSprite
		}
		tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))
		tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Height)))

		rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
	}
	// Drawing the player same as before
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest,
		rl.NewVector2(playerDest.Width/2, playerDest.Height/2), 0, rl.White)
}

func main() {
	Init()
	for running {
		vel := input()
		update(vel)
		render()
	}

	quit()
}
