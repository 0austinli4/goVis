package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type CustomMutex struct {
	sync.Mutex
	locked bool
}

func (cm *CustomMutex) Lock() {
	cm.Mutex.Lock()
	cm.locked = true
}

func (cm *CustomMutex) Unlock() {
	cm.Mutex.Unlock()
	cm.locked = false
}
func (cm *CustomMutex) IsLocked() bool {
	return cm.locked
}

type Philosopher struct {
	name  string // name of philosopher
	count int    // count number of times they eat
	left  int    // fork number on the left
	right int    // fork number on the right
}

func (p *Philosopher) Think() {
	// fmt.Println(p.name, "is thinking.")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	// fmt.Println(p.name, "is done thinking.")
}

func (p *Philosopher) Eat() {
	p.count++
	// fmt.Printf("%s is eating round: %d\n", p.name, p.count)
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	// fmt.Println(p.name, "is done eating.")
}

// golang example solution
// always pick up left fork and then right fork
// can result in deadlock when everyone picks up their left fork
func (p *Philosopher) algoLeft(table []sync.Mutex) {
	for {
		p.Think()
		table[p.left].Lock()
		fmt.Printf("%s picks up fork %d.\n", p.name, p.left)
		runtime.Gosched() // hack to yield to the next goroutine
		table[p.right].Lock()
		fmt.Printf("%s picks up %d.\n", p.name, p.right)
		p.Eat()
		table[p.right].Unlock()
		fmt.Printf("%s puts down fork %d.\n", p.name, p.right)
		table[p.left].Unlock()
		fmt.Printf("%s puts down fork %d.\n", p.name, p.left)
	}
}

// always pick up left fork and then right fork
// can result in deadlock when everyone picks up their left fork
func (p *Philosopher) djikstra(table []CustomMutex) {
	for i := 0; i < 5; i += 1 {
		funcTime := time.Now()
		p.Think()
		table[p.left].Lock()
		fmt.Printf("%s picks up fork %d.\n", p.name, p.left)
		runtime.Gosched() // hack to yield to the next goroutine

		if !table[p.right].IsLocked() {
			table[p.right].Lock()
			fmt.Printf("%s picks up %d.\n", p.name, p.right)
			p.Eat()
			table[p.right].Unlock()
			fmt.Printf("%s puts down fork %d.\n", p.name, p.right)
			table[p.left].Unlock()
			fmt.Printf("%s puts down fork %d.\n", p.name, p.left)
			funcEndTime := time.Now()
			events["Person"+fmt.Sprintln(p.name)+fmt.Sprintln(p.count)+fmt.Sprintln(getGID())] = [2]time.Time{funcTime, funcEndTime}
		} else {
			table[p.left].Unlock()
		}
	}
}

func runPhilosopher(win *pixelgl.Window) {
	// var wg sync.WaitGroup

	// create a wait group so main won't end
	startTime := time.Now()
	philosophers := []*Philosopher{
		{"Michelle", 0, 0, 1},
		{"Bill", 0, 1, 2},
		{"Sonia", 0, 2, 3},
		{"Brooke", 0, 3, 4},
		{"Eric", 0, 4, 0},
	}

	table := make([]CustomMutex, len(philosophers))
	// wg.Add(25)
	for i := 0; i < 5; i += 1 {
		for _, philosopher := range philosophers {
			go func(p *Philosopher) {
				p.djikstra(table)
				// defer wg.Done()
			}(philosopher)
		}
	}

	// wg.Wait()
	time.Sleep(10 * time.Second)
	receiveTime := time.Now()
	events["Main"] = [2]time.Time{startTime, receiveTime}
	fmt.Println("lenght of events", len(events))

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
	runPhilosopher(win)

}

func main() {
	pixelgl.Run(run)
}
