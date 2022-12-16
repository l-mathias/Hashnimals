package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/melonfunction/ebiten-camera"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"log"
	"time"
)

const (
	screenWidth  = 1000
	screenHeight = 480
	sampleRate   = 32000
)

var (
	running     = true
	bkgColor    = color.RGBA{R: 147, G: 211, B: 196, A: 255}
	grassSprite *ebiten.Image
	//go:embed "assets/Tilesets/ground tiles/old tiles/Grass.png"
	GrassPng []byte

	//playerSprite *ebiten.Image
	//go:embed "assets/Characters/Basic Charakter Spritesheet.png"
	PlayerPng []byte

	//go:embed "assets/Sound/AverysFarmLoopable.mp3"
	AverysFarmMP3 []byte
)

type musicType int

const (
	typeOgg musicType = iota
	typeMP3
)

type direction int

const (
	down direction = iota
	up
	left
	right
)

type Pos struct {
	X, Y float64
}

type Player struct {
	sprite *ebiten.Image
	src    image.Rectangle
	dst    image.Rectangle
	Pos
	speed      float64
	moving     bool
	dir        direction
	frameCount int
	frame      int
	size       int
}

type Game struct {
	keys []ebiten.Key
	p    *Player

	//audio
	musicPlayer   *Player
	musicPlayerCh chan *Player
	errCh         chan error

	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
	musicType    musicType

	//camera
	cam        *camera.Camera
	zoomFactor float64

	//tiles
	tileDest   image.Rectangle
	tileSrc    image.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int
}

func (g *Game) Update() error {
	g.cam.SetPosition(g.p.X+float64(48)/2, g.p.Y+float64(48)/2)

	if g.p.moving {
		switch g.p.dir {
		case up:
			g.p.Y -= g.p.speed
		case down:
			g.p.Y += g.p.speed
		case left:
			g.p.X -= g.p.speed
		case right:
			g.p.X += g.p.speed
		}

		if g.p.frameCount%8 == 1 {
			g.p.frame++
		}
	} else if g.p.frameCount%45 == 1 {
		g.p.frame++
	}

	g.p.frameCount++
	if g.p.frame > 3 {
		g.p.frame = 0
	}
	if !g.p.moving && g.p.frame > 1 {
		g.p.frame = 0
	}

	g.p.src.Min.X = g.p.size * g.p.frame
	g.p.src.Min.Y = g.p.size * int(g.p.dir)

	g.p.src.Max.X = g.p.size*g.p.frame + g.p.size
	g.p.src.Max.Y = g.p.size*int(g.p.dir) + g.p.size

	// Zoom
	_, scrollAmount := ebiten.Wheel()
	if scrollAmount > 0 && g.zoomFactor <= 5.0 {
		g.zoomFactor += 0.1
		g.cam.SetZoom(g.zoomFactor)
	} else if scrollAmount < 0 && g.zoomFactor >= 0.5 {
		g.zoomFactor -= 0.1
		g.cam.SetZoom(g.zoomFactor)
	}

	g.p.moving = false

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cam.Surface.Clear()
	g.cam.Surface.Fill(bkgColor)
	tileOps := &ebiten.DrawImageOptions{}
	g.cam.Surface.DrawImage(grassSprite, g.cam.GetTranslation(tileOps, 0, 0))

	playerOps := &ebiten.DrawImageOptions{}
	//playerOps = g.cam.GetRotation(playerOps, 0, -float64(48)/2, -float64(48)/2)
	//playerOps = g.cam.GetScale(playerOps, 1, 1)
	//playerOps = g.cam.GetSkew(playerOps, 0, 0)
	playerOps = g.cam.GetTranslation(playerOps, g.p.X, g.p.Y)
	g.cam.Surface.DrawImage(g.p.sprite.SubImage(g.p.src).(*ebiten.Image), playerOps)
	g.cam.Blit(screen)
	g.cam.SetZoom(g.zoomFactor)
	g.input()
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (g *Game) input() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.p.Y -= g.p.speed
		g.p.moving = true
		g.p.dir = up
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.p.Y += g.p.speed
		g.p.moving = true
		g.p.dir = down
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.p.X -= g.p.speed
		g.p.moving = true
		g.p.dir = left
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.p.X += g.p.speed
		g.p.moving = true
		g.p.dir = right
	}
	if repeatingKeyPressed(ebiten.KeyM) {
		if g.audioPlayer.IsPlaying() {
			g.audioPlayer.Pause()
		} else {
			g.audioPlayer.Play()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if g.zoomFactor >= 0.5 {
			g.zoomFactor -= 0.1
			g.cam.SetZoom(g.zoomFactor)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if g.zoomFactor <= 5.0 {
			g.zoomFactor += 0.1
			g.cam.SetZoom(g.zoomFactor)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		//g.cam.Rotate(0.1)
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		//g.camera.Reset()
	}
}

func NewGame() *Game {
	audioContext := audio.NewContext(sampleRate)
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4

	var s audioStream

	s, err := mp3.DecodeWithoutResampling(bytes.NewReader(AverysFarmMP3))
	if err != nil {
		panic(err)
	}

	p, err := audioContext.NewPlayer(s)

	p.Play()
	return &Game{
		keys:          nil,
		p:             NewPlayer(),
		musicPlayerCh: make(chan *Player),
		errCh:         make(chan error),
		audioContext:  audioContext,
		audioPlayer:   p,
		current:       0,
		total:         time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		seBytes:       nil,
		seCh:          make(chan []byte),
		volume128:     128,
		musicType:     typeMP3,
		cam:           camera.NewCamera(screenWidth, screenHeight, 0, 0, 0, 1.5),
		zoomFactor:    1.5,
	}
}

func NewPlayer() *Player {
	img, _, err := image.Decode(bytes.NewReader(PlayerPng))
	if err != nil {
		log.Fatal(err)
	}
	playerSprite := ebiten.NewImageFromImage(img)

	return &Player{
		sprite: playerSprite,
		src: image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: 48, Y: 48},
		},
		dst:   image.Rectangle{},
		Pos:   Pos{},
		speed: 3,
		size:  48,
	}

}

func init() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hashnimals")

	img, _, err := image.Decode(bytes.NewReader(GrassPng))
	if err != nil {
		log.Fatal(err)
	}

	grassSprite = ebiten.NewImageFromImage(img)

}

func main() {
	g := NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
