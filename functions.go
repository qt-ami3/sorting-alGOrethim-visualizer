package main
import (
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
)

func randInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
	
func (g *Game) initSound() {
	if g.soundPlayer != nil || g.audioContext == nil {
		return
	}

	f, err := os.Open(g.soundFile)
	if err != nil {
		log.Println(err)
		return
	}

	stream, err := mp3.Decode(g.audioContext, f)
	if err != nil {
		log.Println(err)
		return
	}

	player, err := g.audioContext.NewPlayer(stream)
	if err != nil {
		log.Println(err)
		return
	}

	g.soundPlayer = player
}


func (g *Game) playSound() {
	if g.soundPlayer == nil {
		return
	}

	g.soundPlayer.Rewind()
	g.soundPlayer.Play()
}
func GetSelfRAM() float64 {
	
	file, err := os.Open("/proc/self/status")
	
	if err != nil {
		return -1
	}
	
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			kb, _ := strconv.Atoi(fields[1]) // value in KB
			return float64(kb) / 1024        // return MB
		}
	}
	return -1
}
