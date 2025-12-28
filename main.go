package main

import (
	"time"
	"strings"
	"fmt"
	"log"
	"os"
	"bufio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 48000 //sample rate may or may not crash program if changed
												 //on a per distro & soundsystem basis

type Game struct {
	data []int
	i, j int
	sorted bool

	fillIndex int
	lastSwap  int

	audioContext *audio.Context
	soundPlayer  *audio.Player
	soundFile    string
}



var (	
	muted = bool(false)
	sortSelected = bool(false)
	gameSpeed = int(15)
	visualizerBar *ebiten.Image
	visualizerPosition = float64(0)
	tickCount = int(0)
	screenHeight = int(960)
	screenY = float64(screenHeight)
	sortSelection rune
	insortStepOne = bool(true)
	insortStepTwo = bool(false)
	insortStepThree = bool (false)
)

type barStats struct {
	height int
}

func init() {
	var err error

	visualizerBar, _, err = ebitenutil.NewImageFromFile("assets/Sprite-0001.png")
	if err != nil {
		log.Fatal(err)
	}
}

func NewGame() *Game {
	return &Game{
		soundFile: "assets/boop.mp3",
		lastSwap:  -1,
	}
}


func (g *Game) Update() error { //game logic
	
	if g.audioContext == nil {
		g.audioContext = audio.NewContext(sampleRate)
		g.initSound()
	}


	if !sortSelected {
		fmt.Println ("WARNING! LOWER VOLUME!")
		
		fmt.Println ("program will unlock in 3...")
		time.Sleep(1 * time.Second)
		fmt.Println ("program will unlock in 2...")
		time.Sleep(1 * time.Second)
		fmt.Println ("program will unlock in 1...")
		time.Sleep(1 * time.Second)
		g.playSound()

		fmt.Println ("up key will mute, down key will unmute. muting may drastically increase performence with higher tps")
		fmt.Println ("please select sorting algorithm:")
		fmt.Println ("d.) double")
		fmt.Println ("i.) insertion")

		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		sortSelection = rune(strings.TrimSpace(line)[0])

		sortSelected = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		muted = true
	}
	
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		muted = false
	}

	if g.sorted {
		if tickCount%5 == 0 && g.fillIndex < len(g.data) {
			g.fillIndex++
		}
		tickCount++
		return nil
	}

	if g.data == nil { // Initialize once
		size := randInt(320,320)
		g.data = make([]int, size)
		for i := range g.data {
			g.data[i] = randInt(1, 220)
		}
		return nil
	}

	if g.sorted {
		return nil
	}

	if sortSelection == 'd' {
		if g.data[g.j] > g.data[g.j+1] { //one comparison per tick
    	g.data[g.j], g.data[g.j+1] = g.data[g.j+1], g.data[g.j]	
			g.lastSwap = g.j + 1 // highlight the value that moved

			if !muted {
				g.playSound()
			}
		}
		g.j++ // Advance inner index
	}	

	if sortSelection == 'i' {

		// Initialize insertion sort
		if g.i == 0 {
			g.i = 1
			g.j = g.i
		}

		// If done
		if g.i >= len(g.data) {
			g.sorted = true
			return nil
		}

		// One comparison per tick
		if g.j > 0 && g.data[g.j] < g.data[g.j-1] {
			g.data[g.j], g.data[g.j-1] = g.data[g.j-1], g.data[g.j]
			g.lastSwap = g.j - 1
			g.j-- // move left
			if !muted {
				g.playSound()
			}
			return nil
		}

	// Element is placed, move to next i
	g.i++
	g.j = g.i
	return nil
}

	if g.j >= len(g.data)-g.i-1 {	// End of inner pass
		g.j = 0
		g.i++
	}

	if g.i >= len(g.data) - 1 { // Fully sorted
		g.sorted = true

		g.fillIndex = 0 // start from rightmost bar
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	visualizerPosition = 0

	for i := range g.data {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.01, float64(g.data[i]) * (-0.01))
		op.GeoM.Translate(float64(100*(visualizerPosition*0.01)), (screenY/4))

		if g.sorted && i < g.fillIndex { // Fill green 
			op.ColorM.Scale(0, 1, 0, 1)
		}

		if i == g.lastSwap {
			op.ColorM.Scale(0, 1, 0, 1) // green
		}

		screen.DrawImage(visualizerBar, op)
		visualizerPosition++
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetTPS(gameSpeed)
	ebiten.SetWindowSize(1280, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")

	game := NewGame() // <-- use the constructor
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
