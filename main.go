package main

import (
	"time"
	"strings"
	"fmt"
	"log"
	"os"
	"bufio"
	"unsafe"
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
	swapped  bool
	forward  bool
	oddPhase bool
	gap      int
}



var (	
	muted = bool(true)
	sortSelected = bool(false)
	gameSpeed = int(100000)
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
	
	fmt.Println(GetSelfRAM())

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
		fmt.Println ("b.) bubble")
		fmt.Println ("i.) insertion")
		fmt.Println ("s.) selection")
		fmt.Println ("c.) cocktail (bidirectional bubble)")
		fmt.Println ("g.) gnome")
		fmt.Println ("o.) odd-even")
		fmt.Println ("C.) comb (uppercase C)")
		fmt.Println ("B.) bogo (uppercase B - WARNING: MAY NEVER FINISH!)")

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

	// DOUBLE SORT - 'd'
	if sortSelection == 'd' {
		if g.data[g.j] > g.data[g.j+1] { //one comparison per tick
    	g.data[g.j], g.data[g.j+1] = g.data[g.j+1], g.data[g.j]	
			g.lastSwap = g.j + 1 // highlight the value that moved

			if !muted {
				g.playSound()
			}
		}
		g.j++ // Advance inner index
		
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

	// BUBBLE SORT - 'b'
	if sortSelection == 'b' {
		// If done
		if g.i >= len(g.data)-1 {
			g.sorted = true
			g.fillIndex = 0
			return nil
		}
		
		// One comparison per tick
		if g.j < len(g.data)-g.i-1 {
			if g.data[g.j] > g.data[g.j+1] {
				g.data[g.j], g.data[g.j+1] = g.data[g.j+1], g.data[g.j]
				g.lastSwap = g.j + 1
				if !muted {
					g.playSound()
				}
			}
			g.j++
			return nil
		}
		
		// End of pass
		g.j = 0
		g.i++
		return nil
	}

	// INSERTION SORT - 'i'
	if sortSelection == 'i' {
		// Initialize insertion sort
		if g.i == 0 {
			g.i = 1
			g.j = g.i
		}

		// If done
		if g.i >= len(g.data) {
			g.sorted = true
			g.fillIndex = 0
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

	// SELECTION SORT - 's'
	if sortSelection == 's' {
		// If done
		if g.i >= len(g.data)-1 {
			g.sorted = true
			g.fillIndex = 0
			return nil
		}
		
		// Initialize min index for this pass
		if g.j == g.i {
			// Store min index in a custom field or reuse lastSwap temporarily
			// We'll track it via j's behavior
		}
		
		// One comparison per tick - find minimum
		if g.j < len(g.data) {
			if g.data[g.j] < g.data[g.i] {
				// Found new minimum, swap immediately for visualization
				g.data[g.i], g.data[g.j] = g.data[g.j], g.data[g.i]
				g.lastSwap = g.i
				if !muted {
					g.playSound()
				}
			}
			g.j++
			return nil
		}
		
		// End of pass
		g.i++
		g.j = g.i
		return nil
	}

	// COCKTAIL SORT (bidirectional bubble) - 'c'
	if sortSelection == 'c' {
		// Initialize
		if g.i == 0 && g.j == 0 {
			g.forward = true
			g.swapped = false
		}
		
		// If done
		if g.i >= len(g.data)/2 {
			g.sorted = true
			g.fillIndex = 0
			return nil
		}
		
		// Forward pass
		if g.forward {
			if g.j < len(g.data)-g.i-1 {
				if g.data[g.j] > g.data[g.j+1] {
					g.data[g.j], g.data[g.j+1] = g.data[g.j+1], g.data[g.j]
					g.lastSwap = g.j + 1
					g.swapped = true
					if !muted {
						g.playSound()
					}
				}
				g.j++
				return nil
			}
			// Switch to backward pass
			g.forward = false
			g.j = len(g.data) - g.i - 2
			return nil
		}
		
		// Backward pass
		if g.j >= g.i {
			if g.data[g.j] > g.data[g.j+1] {
				g.data[g.j], g.data[g.j+1] = g.data[g.j+1], g.data[g.j]
				g.lastSwap = g.j
				g.swapped = true
				if !muted {
					g.playSound()
				}
			}
			g.j--
			return nil
		}
		
		// End of complete pass
		if !g.swapped {
			g.sorted = true
			g.fillIndex = 0
			return nil
		}
		g.swapped = false
		g.forward = true
		g.i++
		g.j = g.i
		return nil
	}

	// GNOME SORT - 'g'
	if sortSelection == 'g' {
		// Initialize
		if g.i == 0 {
			g.i = 1
		}
		
		// If done
		if g.i >= len(g.data) {
			g.sorted = true
			g.fillIndex = 0
			return nil
		}
		
		// One comparison per tick
		if g.i > 0 && g.data[g.i] < g.data[g.i-1] {
			g.data[g.i], g.data[g.i-1] = g.data[g.i-1], g.data[g.i]
			g.lastSwap = g.i - 1
			g.i-- // move backward
			if !muted {
				g.playSound()
			}
			return nil
		}
		
		// Move forward
		g.i++
		return nil
	}

	// ODD-EVEN SORT - 'o'
	if sortSelection == 'o' {
		// Initialize
		if g.i == 0 && g.j == 0 {
			g.oddPhase = false
			g.swapped = false
		}
		
		// Calculate starting position based on phase
		startPos := 0
		if g.oddPhase {
			startPos = 1
		}
		
		// One comparison per tick
		if g.j < len(g.data)-1 {
			// Only compare if we're at the right position for this phase
			if g.j >= startPos && (g.j-startPos)%2 == 0 {
				if g.data[g.j] > g.data[g.j+1] {
					g.data[g.j], g.data[g.j+1] = g.data[g.j+1], g.data[g.j]
					g.lastSwap = g.j + 1
					g.swapped = true
					if !muted {
						g.playSound()
					}
				}
			}
			g.j++
			return nil
		}
		
		// End of phase - need two clean passes to be done
		if !g.swapped {
			g.i++ // Count clean passes
			if g.i >= 2 {
				g.sorted = true
				g.fillIndex = 0
				return nil
			}
		} else {
			g.i = 0 // Reset clean pass counter
		}
		
		g.oddPhase = !g.oddPhase
		g.j = 0
		g.swapped = false
		return nil
	}

	// COMB SORT - 'C' (uppercase)
	if sortSelection == 'C' {
		// Initialize gap
		if g.gap == 0 {
			g.gap = len(g.data)
			g.swapped = false
		}
		
		// Update gap
		if g.j == 0 && g.gap > 1 {
			g.gap = (g.gap * 10) / 13
			if g.gap < 1 {
				g.gap = 1
			}
		}
		
		// One comparison per tick
		if g.j < len(g.data)-g.gap {
			if g.data[g.j] > g.data[g.j+g.gap] {
				g.data[g.j], g.data[g.j+g.gap] = g.data[g.j+g.gap], g.data[g.j]
				g.lastSwap = g.j + g.gap
				g.swapped = true
				if !muted {
					g.playSound()
				}
			}
			g.j++
			return nil
		}
		
		// End of pass
		if g.gap == 1 && !g.swapped {
			g.sorted = true
			g.fillIndex = 0
			return nil
		}
		
		g.j = 0
		g.swapped = false
		return nil
	}

	// BOGO SORT - 'B' (uppercase)
	if sortSelection == 'B' {
		// Only do full sorted check once per complete shuffle cycle
		// Use g.i to track if we need to check
		if g.i == 0 {
			// Check if sorted using pointers
			isSorted := true
			dataPtr := &g.data[0]
			dataLen := len(g.data)
			
			for i := 0; i < dataLen-1; i++ {
				if *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(dataPtr)) + uintptr(i)*unsafe.Sizeof(int(0)))) > 
				   *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(dataPtr)) + uintptr(i+1)*unsafe.Sizeof(int(0)))) {
					isSorted = false
					break
				}
			}
			
			if isSorted {
				g.sorted = true
				g.fillIndex = 0
				return nil
			}
			
			// Start new shuffle cycle
			g.i = len(g.data) / 2 // Do half-array worth of swaps before checking again
		}
		
		// One shuffle per tick - swap two random elements using pointers
		idx1 := randInt(0, len(g.data)-1)
		idx2 := randInt(0, len(g.data)-1)
		
		if idx1 != idx2 {
			// Pointer-based swap
			dataPtr := &g.data[0]
			ptr1 := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(dataPtr)) + uintptr(idx1)*unsafe.Sizeof(int(0))))
			ptr2 := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(dataPtr)) + uintptr(idx2)*unsafe.Sizeof(int(0))))
			
			*ptr1, *ptr2 = *ptr2, *ptr1
			
			g.lastSwap = idx1
			if idx2 < idx1 {
				g.lastSwap = idx2
			}
			if !muted {
				go g.playSound()
			}
		}
		
		g.i-- // Decrement shuffle counter
		return nil
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
	ebiten.SetWindowTitle("Sorting algorithm visualizer")
	ebiten.SetFullscreen(false)

	game := NewGame() // <-- use the constructor
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
