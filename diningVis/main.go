package main

import (
	_ "image/png"
	"sync"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// optional function - you can use either runOnClick() or
//	runConstantAnimation()
func runOnClick() {
	var wg sync.WaitGroup
	wg.Add(1)
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Aliceblue)

	initializeForks(win)
	initializePhilosophers(win)

	for !win.JustPressed(pixelgl.MouseButtonLeft) {
		win.Update()
	}

	win.Update()

	for _, philosopher := range philosophers {
		go func(p *Philosopher) {
			p.Dine(win)
		}(philosopher)
	}

	wg.Wait()

}

func runConstantAnimation() {
	var wg sync.WaitGroup
	wg.Add(1)
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Aliceblue)

	initializeForks(win)
	initializePhilosophers(win)

	for !win.JustPressed(pixelgl.MouseButtonLeft) {
		win.Update()
	}

	win.Update()

	for !win.Closed() {
		for _, philosopher := range philosophers {
			go func(p *Philosopher) {
				p.Dine(win)
			}(philosopher)
		}
		win.Update()
	}

	wg.Wait()

}

func main() {
	pixelgl.Run(runConstantAnimation)
}
