package main

import (
	"flag"
	"log"
	"os"

	"github.com/MagicalTux/gones/cartridge"
	"github.com/MagicalTux/gones/cpu2a03"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	cpu *cpu2a03.Cpu2A03
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 240
}

func main() {
	flag.Parse()
	arg := flag.Args()
	if len(arg) != 1 {
		log.Printf("Usage: %s file.nes", os.Args[0])
		os.Exit(1)
	}

	data, err := cartridge.Load(arg[0])
	if err != nil {
		log.Printf("Failed to load %s: %s", arg[0], err)
		os.Exit(1)
	}

	cpu := cpu2a03.New()

	err = data.Mapper.Setup(cpu.Memory)
	if err != nil {
		log.Printf("Failed to map %s: %s", arg[0], err)
		os.Exit(1)
	}

	log.Printf("CPU ready with memory: %s", cpu.Memory)

	cpu.Reset()
	go cpu.Start(cpu2a03.NTSC)

	ebiten.SetWindowSize(256*2, 240*2)
	ebiten.SetWindowTitle("goNES")

	game := &Game{
		cpu: cpu,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
