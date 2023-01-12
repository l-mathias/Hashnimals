package main

import (
	"fmt"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/lafriks/go-tiled"
	"image/color"
	"os"
)

type direction int

const (
	screenWidth  = 800
	screenHeight = 600
	mapTileSize  = 16
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
	playerTileX  int
	playerTileY  int

	frameCount int

	tileDst rl.Rectangle
	tileSrc rl.Rectangle

	tiledMap           *tiled.Map
	backgroundMap      [][]Tile
	detailsMap         [][]Tile
	collisionsMap      [][]Tile
	debugCollisionsMap []debugRect

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

type Pos struct {
	X, Y int
}

type Tile struct {
	id int
	Pos
	tileSetName string
	collision   bool
}

type debugRect struct {
	rect  rl.Rectangle
	color color.RGBA
}

func drawBackground() {
	for y := range backgroundMap {
		for x := range backgroundMap[y] {
			if !tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].IsNil() {
				tileSrc.X = float32(tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].GetTileRect().Min.X)
				tileSrc.Y = float32(tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].GetTileRect().Min.Y)
				tileSrc.Width = float32(tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].GetTileRect().Size().X)
				tileSrc.Height = float32(tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].GetTileRect().Size().Y)

				tileDst.X = tileDst.Width * float32(x)
				tileDst.Y = tileDst.Height * float32(y)
				rl.DrawTexturePro(textureIndex[tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].Tileset.Name], tileSrc, tileDst, rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}
}

func drawDetails() {
	for y := range detailsMap {
		for x := range detailsMap[y] {
			if !tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].IsNil() {
				tileSrc.X = float32(tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].GetTileRect().Min.X)
				tileSrc.Y = float32(tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].GetTileRect().Min.Y)
				tileSrc.Width = float32(tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].GetTileRect().Size().X)
				tileSrc.Height = float32(tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].GetTileRect().Size().Y)

				tileDst.X = tileDst.Width * float32(x)
				tileDst.Y = tileDst.Height * float32(y)
				rl.DrawTexturePro(textureIndex[tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].Tileset.Name], tileSrc, tileDst, rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}
}

func drawCollisions() {
	for _, rect := range debugCollisionsMap {
		rl.DrawRectangleRec(rect.rect, rect.color)
	}
}

func drawScene() {

	drawBackground()
	drawDetails()
	//drawCollisions()

	//playerTileX = int((playerDst.X + mapTileSize/2) / mapTileSize)
	//playerTileY = int((playerDst.Y + mapTileSize/2) / mapTileSize)
	//playerPos := "X: " + strconv.Itoa(playerTileX) + " , " + " Y: " + strconv.Itoa(playerTileY)

	rl.DrawTexturePro(playerSprite, playerSrc, playerDst, rl.NewVector2(16, 16), 0, rl.White)
	//rl.DrawText(playerPos, int32(playerDst.X), int32(playerDst.Y), 20, rl.Lime)
}

func input() {
	playerPos := rl.Rectangle{
		X:      playerDst.X,
		Y:      playerDst.Y,
		Width:  playerDst.Width / 5,
		Height: playerDst.Height / 5,
	}
	if rl.IsKeyDown(rl.KeyUp) {
		playerPos.Y -= playerSpeed
		if !checkCollision(playerPos) {
			playerDst.Y -= playerSpeed
			playerMoving = true
			playerDir = int(playerUp)
		}
	}
	if rl.IsKeyDown(rl.KeyDown) {
		playerPos.Y += playerSpeed
		if !checkCollision(playerPos) {
			playerDst.Y += playerSpeed
			playerMoving = true
			playerDir = int(playerDown)
		}
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		playerPos.X -= playerSpeed
		if !checkCollision(playerPos) {
			playerDst.X -= playerSpeed
			playerMoving = true
			playerDir = int(playerLeft)
		}
	}
	if rl.IsKeyDown(rl.KeyRight) {
		playerPos.X += playerSpeed
		if !checkCollision(playerPos) {
			playerDst.X += playerSpeed
			playerMoving = true
			playerDir = int(playerRight)
		}
	}
	if rl.IsKeyDown(rl.KeyQ) {
		musicPaused = !musicPaused
	}
	if rl.IsKeyDown(rl.KeyC) {
		debugCollisionsMap = make([]debugRect, 0)
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

func checkCollision(dst rl.Rectangle) bool {
	for y := range collisionsMap {
		for x, r := range collisionsMap[y] {
			if !r.collision {
				rect := rl.Rectangle{
					X:      float32(x) * tileSrc.Width,
					Y:      float32(y) * tileSrc.Height,
					Width:  tileSrc.Width,
					Height: tileSrc.Height,
				}

				if rl.CheckCollisionRecs(rect, dst) {
					//fmt.Println("Collision found !")
					//fmt.Println("Rect : ", rect)
					//fmt.Println("Player Pos : ", dst)
					//
					debugCollisionsMap = append(debugCollisionsMap, debugRect{
						rect:  dst,
						color: color.RGBA{R: 255, A: 255},
					})
					debugCollisionsMap = append(debugCollisionsMap, debugRect{
						rect:  rect,
						color: color.RGBA{G: 255, A: 255},
					})
					//
					//fmt.Println()
					return true
				}
			}
		}
	}

	return false
}

func update() {
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

	cam.Target = rl.NewVector2(playerDst.X, playerDst.Y)

	playerMoving = false
}

func loadTMX(mapFile string) {
	var err error
	tiledMap, err = tiled.LoadFile(mapFile)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	textureIndex = make(map[string]rl.Texture2D)
	backgroundMap = make([][]Tile, tiledMap.Height)
	detailsMap = make([][]Tile, tiledMap.Height)
	collisionsMap = make([][]Tile, tiledMap.Height)
	debugCollisionsMap = make([]debugRect, 0)

	for i := range backgroundMap {
		backgroundMap[i] = make([]Tile, tiledMap.Width)
		detailsMap[i] = make([]Tile, tiledMap.Width)
		collisionsMap[i] = make([]Tile, tiledMap.Width)
	}

	for y := 0; y < tiledMap.Height; y++ {
		for x := 0; x < tiledMap.Width; x++ {
			if !tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].IsNil() {
				backgroundMap[y][x].tileSetName = tiledMap.Layers[0].Tiles[y*tiledMap.Width+x].Tileset.Name
			}
			backgroundMap[y][x].id = y*tiledMap.Width + x

			if !tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].IsNil() {
				detailsMap[y][x].tileSetName = tiledMap.Layers[1].Tiles[y*tiledMap.Width+x].Tileset.Name
			}
			detailsMap[y][x].id = y*tiledMap.Width + x

			//add collision info
			collisionsMap[y][x].collision = tiledMap.Layers[2].Tiles[y*tiledMap.Width+x].Nil
			collisionsMap[y][x].id = y*tiledMap.Width + x
		}
	}

	for _, layer := range tiledMap.Layers {
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

	tileDst = rl.NewRectangle(0, 0, mapTileSize, mapTileSize)
	tileSrc = rl.NewRectangle(0, 0, mapTileSize, mapTileSize)

	playerSprite = rl.LoadTexture("assets/Characters/Basic Charakter Spritesheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDst = rl.NewRectangle(250, 150, 48, 48)

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("assets/Sound/AverysFarmLoopable.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth)/2, float32(screenHeight)/2), rl.NewVector2(playerDst.X, playerDst.Y), 0.0, 3.0)

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
