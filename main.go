package main

import (
	"github.com/gen2brain/raylib-go/raylib"
	"os"
	"strconv"
	"strings"
)

type direction int

const (
	screenWidth  = 800
	screenHeight = 450
)

const (
	playerDown direction = iota
	playerUp
	playerLeft
	playerRight
)

var (
	running      = true
	bkgColor     = rl.NewColor(147, 211, 196, 255)
	grassSprite  rl.Texture2D
	playerSprite rl.Texture2D

	playerSrc    rl.Rectangle
	playerDst    rl.Rectangle
	playerSpeed  float32 = 3
	playerMoving bool
	playerDir    int
	playerFrame  int

	frameCount int

	tileDst    rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

func drawScene() {
	//rl.DrawTexture(grassSprite, 100, 50, rl.White)

	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDst.X = tileDst.Width * float32(i%mapW)
			tileDst.Y = tileDst.Height * float32(i/mapH)
			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(grassSprite.Width/int32(tileSrc.Width)))
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(grassSprite.Width/int32(tileSrc.Width)))

			rl.DrawTexturePro(grassSprite, tileSrc, tileDst, rl.NewVector2(tileDst.Width, tileDst.Height), 0, rl.White)
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
	running = !rl.WindowShouldClose()

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

func loadMap(mapFile string) {
	file, err := os.ReadFile(mapFile)
	if err != nil {
		panic(err)
	}
	remNewLines := strings.Replace(string(file), "\r\n", " ", -1)
	sliced := strings.Split(remNewLines, " ")
	mapW = -1
	mapH = -1
	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.Atoi(sliced[i])
		if mapW == -1 {
			mapW = s
		} else if mapH == -1 {
			mapH = s
		} else if i < mapW*mapH+2 {
			tileMap = append(tileMap, s)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}
}

func init() {
	rl.InitWindow(screenWidth, screenHeight, "Hashnimals")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	grassSprite = rl.LoadTexture("assets/Tilesets/ground tiles/old tiles/Grass.png")

	tileDst = rl.NewRectangle(0, 0, 16, 16)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	playerSprite = rl.LoadTexture("assets/Characters/Basic Charakter Spritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDst = rl.NewRectangle(200, 200, 100, 100)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/Sound/AverysFarmLoopable.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2), rl.NewVector2(playerDst.X-(playerDst.Width/2), playerDst.Y-(playerDst.Height/2)), 0.0, 1.0)
	loadMap("level1.map")
}
func quit() {
	rl.UnloadTexture(grassSprite)
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
