package main
import (
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"os"
	"log"
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
