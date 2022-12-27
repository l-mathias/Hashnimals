package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"os"
)

type direction int

const (
	screenWidth  = 800
	screenHeight = 600
)

const (
	playerDown direction = iota
	playerUp
	playerLeft
	playerRight
)

var (
	bkgColor     = rl.NewColor(147, 211, 196, 255)
	textureIndex map[string]rl.Texture2D

	playerSprite rl.Texture2D

	playerSrc    rl.Rectangle
	playerDst    rl.Rectangle
	playerSpeed  float32 = 1.4
	playerMoving bool
	playerDir    int
	playerFrame  int

	frameCount int

	tileDst rl.Rectangle
	tileSrc rl.Rectangle

	tiledMap *tiled.Map

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

func drawScene() {
	for _, layer := range tiledMap.Layers {
		if layer.Visible {
			for i, tile := range layer.Tiles {
				if !tile.IsNil() {
					tileSrc.X = float32(tile.GetTileRect().Min.X)
					tileSrc.Y = float32(tile.GetTileRect().Min.Y)
					tileSrc.Width = float32(tile.GetTileRect().Size().X)
					tileSrc.Height = float32(tile.GetTileRect().Size().Y)

					tileDst.X = tileDst.Width * float32(i%tiledMap.Width)
					tileDst.Y = tileDst.Width * float32(i/tiledMap.Width)

					rl.DrawTexturePro(textureIndex[tile.Tileset.Name], tileSrc, tileDst, rl.NewVector2(tileDst.Width, tileDst.Height), 0, rl.White)
				} else {

				}
			}
		}
	}

	rl.DrawTexturePro(playerSprite, playerSrc, playerDst, rl.NewVector2(playerDst.Width, playerDst.Height), 0, rl.White)
}

func input() {
	if rl.IsKeyDown(rl.KeyUp) {
		playerDst.Y -= playerSpeed
		playerMoving = true
		playerDir = int(playerUp)
	}
	if rl.IsKeyDown(rl.KeyDown) {
		playerDst.Y += playerSpeed
		playerMoving = true
		playerDir = int(playerDown)
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		playerDst.X -= playerSpeed
		playerMoving = true
		playerDir = int(playerLeft)
	}
	if rl.IsKeyDown(rl.KeyRight) {
		playerDst.X += playerSpeed
		playerMoving = true
		playerDir = int(playerRight)
	}
	if rl.IsKeyDown(rl.KeyQ) {
		musicPaused = !musicPaused
	}
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)

	rl.BeginMode2D(cam)
	drawScene()
	rl.EndMode2D()

	rl.EndDrawing()
}

func update() {
	//running = !rl.WindowShouldClose()

	playerSrc.X = 0
	playerSrc.X = playerSrc.Width * float32(playerFrame)
	if playerMoving {
		switch playerDir {
		case int(playerUp):
			playerDst.Y -= playerSpeed
		case int(playerDown):
			playerDst.Y += playerSpeed
		case int(playerLeft):
			playerDst.X -= playerSpeed
		case int(playerRight):
			playerDst.X += playerSpeed
		}
		if frameCount%8 == 1 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		playerFrame++
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
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

	cam.Target = rl.NewVector2(playerDst.X-(playerDst.Width/2), playerDst.Y-(playerDst.Height/2))

	playerMoving = false
}

func loadTMX(mapFile string) {
	var err error
	textureIndex = make(map[string]rl.Texture2D)

	tiledMap, err = tiled.LoadFile(mapFile)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	for _, layer := range tiledMap.Layers {
		//fmt.Println(layer.Name)
		for _, tile := range layer.Tiles {
			if !tile.IsNil() {
				_, ok := textureIndex[tile.Tileset.Name]
				if !ok {
					textureIndex[tile.Tileset.Name] = rl.LoadTexture(tile.Tileset.Image.Source)
				}
			}
		}
	}
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "Hashnimals")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	tileDst = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSprite = rl.LoadTexture("assets/Characters/Basic Charakter Spritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDst = rl.NewRectangle(200, 200, 60, 60)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/Sound/AverysFarmLoopable.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2), rl.NewVector2(playerDst.X-(playerDst.Width/2), playerDst.Y-(playerDst.Height/2)), 0.0, 1.0)
	cam.Zoom = 3

	loadTMX("island.tmx")
}
func quit() {
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
func main() {
	for !rl.WindowShouldClose() {
		input()
		update()
		render()
	}
	quit()

}
