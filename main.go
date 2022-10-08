package main

import (
	"flag"
	"image"
	"log"
	"os"
	"runtime/pprof"

	"github.com/MagicalTux/gones/apu"
	"github.com/MagicalTux/gones/cartridge"
	"github.com/MagicalTux/gones/cpu2a03"
	"github.com/MagicalTux/gones/nesinput"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	cputrace   = flag.String("trace", "", "write 6502 instructions to file")
	ppudebug   = flag.String("ppudebug", "", "write PPU debug info to file, or - for stdout")
	zoom       = flag.Int("zoom", 4, "zoom level for display")
	startV     = flag.Int("start_v", 0, "define start position in RAM, for ex 0xc000")
)

type Game struct {
	cpu     *cpu2a03.Cpu2A03
	img     *ebiten.Image
	started bool
	gamepad ebiten.GamepadID
}

func (g *Game) Update() error {
	if !g.started {
		g.started = true
		g.cpu.Reset()
		if *startV != 0 {
			g.cpu.PC = uint16(*startV)
		}
		g.cpu.Start()

		snd := audio.NewContext(44100)
		if player, err := snd.NewPlayer(g.cpu.APU); err != nil {
			log.Printf("failed to create player: %s", err)
		} else {
			log.Printf("Audio: setting buffer length to %s", apu.BufferLength())
			player.SetBufferSize(apu.BufferLength())
			go player.Play()
		}
	}

	if g.gamepad != 0 {
		if inpututil.IsGamepadJustDisconnected(g.gamepad) {
			// return to keyboard control
			g.cpu.Input[0] = nesinput.NewKeyboard()
			g.gamepad = 0
		}
	} else {
		cn := inpututil.AppendJustConnectedGamepadIDs(nil)
		if len(cn) > 0 {
			// take first gamepad
			id := cn[0]
			pad, err := nesinput.NewGamepad(id)
			if err != nil {
				log.Printf("enable gamepad failed: %s", err)
			} else {
				g.gamepad = id
				g.cpu.Input[0] = pad
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cpu.PPU.Front(func(img *image.RGBA) {
		g.img.WritePixels(img.Pix)
	})
	screen.DrawImage(g.img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 240
}

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	cpu := cpu2a03.New(cpu2a03.NTSC)

	cpu.Input[0] = nesinput.NewKeyboard()

	if *cputrace != "" {
		cpu.Trace, err = os.Create(*cputrace)
		if err != nil {
			log.Printf("Failed to create %s: %s", *cputrace, err)
			os.Exit(1)
		}
	}
	if *ppudebug != "" {
		if *ppudebug == "-" {
			cpu.PPU.Trace = os.Stdout
		} else {
			cpu.PPU.Trace, err = os.Create(*ppudebug)
			if err != nil {
				log.Printf("Failed to create %s: %s", *ppudebug, err)
				os.Exit(1)
			}
		}
	}

	err = data.Setup(cpu)
	if err != nil {
		log.Printf("Failed to map %s: %s", arg[0], err)
		os.Exit(1)
	}

	log.Printf("CPU ready with memory: %s", cpu.Memory)
	log.Printf("PPU ready with memory: %s", cpu.PPU.Memory)

	ebiten.SetWindowSize(256*(*zoom), 240*(*zoom))
	ebiten.SetWindowTitle("goNES")

	game := &Game{
		cpu: cpu,
		img: ebiten.NewImage(256, 240),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
