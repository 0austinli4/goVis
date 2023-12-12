// Program with race condition fixed by mutex
package main

import (
	"fmt"
	"sync" // to import sync later on
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var GFG = 0

var globalLock sync.Mutex

// This is the function we’ll run in every
// goroutine. Note that a WaitGroup must
// be passed to functions by pointer.
func workerLock(wg *sync.WaitGroup, m *sync.Mutex) {
	// Lock() the mutex to ensure
	// exclusive access to the state,
	// increment the value,
	// Unlock() the mutex
	m.Lock()
	funcTime := time.Now()
	GFG = GFG + 1
	time.Sleep(1 * time.Second)
	funcEndTime := time.Now()
	events["Work"+fmt.Sprintln(GFG)+fmt.Sprintln(getGID())] = [2]time.Time{funcTime, funcEndTime}
	m.Unlock()
	// On return, notify the
	// WaitGroup that we’re done.
	wg.Done()
}

// Program with race condition

// This is the function we’ll run in every
// goroutine. Note that a WaitGroup must
// be passed to functions by pointer.
func worker(wg *sync.WaitGroup) {
	funcTime := time.Now()
	GFG = GFG + 1
	time.Sleep(1 * time.Second)
	funcEndTime := time.Now()

	globalLock.Lock()
	defer globalLock.Unlock()
	events["Work"+fmt.Sprintln(GFG)+fmt.Sprintln(getGID())] = [2]time.Time{funcTime, funcEndTime}
	// On return, notify the
	// WaitGroup that we’re done.
	wg.Done()
}

func runWorkersRace(win *pixelgl.Window) {
	GFG = 0
	startTime := time.Now()
	// wait group to make sure goroutines finish
	var w sync.WaitGroup
	// Launch several goroutines and increment
	for i := 0; i < 15; i++ {
		// add one per go routine
		w.Add(1)
		go worker(&w)
	}
	// Block until waitgroup is at 0
	w.Wait()
	fmt.Println("Value of x", GFG)
	receiveTime := time.Now()

	events["Main"] = [2]time.Time{startTime, receiveTime}
	animateChannel(win)
}

// https://www.geeksforgeeks.org/mutex-in-golang-with-examples/#
func runWorkersNoRace(win *pixelgl.Window) {
	startTime := time.Now()
	// This WaitGroup is used to wait for
	// all the goroutines launched here to finish.
	var w sync.WaitGroup

	// This mutex will synchronize access to state.
	var m sync.Mutex

	// Launch several goroutines and increment
	// the WaitGroup counter for each
	for i := 0; i < 15; i++ {
		w.Add(1)
		go workerLock(&w, &m)
	}
	// Block until the WaitGroup counter
	// goes back to 0; all the workers
	// notified they’re done.
	w.Wait()
	fmt.Println("Value of x", GFG)
	receiveTime := time.Now()
	events["Main"] = [2]time.Time{startTime, receiveTime}
	animateChannel(win)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear((colornames.White))
	runWorkersNoRace(win)
	runWorkersRace(win)
}

func main() {
	pixelgl.Run(run)
}
