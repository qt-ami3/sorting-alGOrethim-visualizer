# a sorting algorithm visualiser for the ebitengine.

## To run:

Instructions are for debian/ubuntu based distrobutions, instructions for other distrobutions can be found on ebitengine website.

DO NOT unmute unless you aim to optimize or otherwise run g.playSound() on its own thread(s). This works easily with go g.playSound() but uses a considerable amount of memory and, statistically, is practically guarenteed to crash whatever system your using before it completes the sort, it also loops rhythmically after enough Goroutines are started.

```
sudo apt install gcc

sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

go mod init github.com/name/gamename

go mod tidy

go run .
```
