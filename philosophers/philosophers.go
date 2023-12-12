package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
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

// djisktras solution???
// always pick up left fork and then right fork
// can result in deadlock when everyone picks up their left fork
func (p *Philosopher) djikstra(table []CustomMutex) {
	for {
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
		} else {
			table[p.left].Unlock()
		}
	}
}

func main() {

	// create a wait group so main won't end

	philosophers := []*Philosopher{
		{"Michelle", 0, 0, 1},
		{"Bill", 0, 1, 2},
		{"Sonia", 0, 2, 3},
		{"Brooke", 0, 3, 4},
		{"Eric", 0, 4, 0},
	}

	table := make([]CustomMutex, len(philosophers))

	for {
		for _, philosopher := range philosophers {
			go func(p *Philosopher) {
				// p.algoLeft(table)
				p.djikstra(table)
			}(philosopher)
		}

	}

	// time.Sleep(20 * time.Second)
	// for _, philosopher := range philosophers {
	// 	name := fmt.Sprint(philosopher.name)
	// 	count := fmt.Sprint(philosopher.count)
	// 	fmt.Println(name + "EATEN " + count)
	// }
}
