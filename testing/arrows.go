// Program with race condition fixed by mutex
package main

import (
	"fmt"
	"sync" // to import sync later on
	"time"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var wg sync.WaitGroup

type Bank struct {
	owner *User
	mu sync.Mutex
	balance int
}

type User struct {
	name string
	balance int
}

func (user *User) deposit(bank *Bank, amount int) {
	defer wg.Done()
	bank.owner = user
	bank.mu.Lock()

	user.balance -= amount
	bank.balance += amount

	// build-in delay
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	bank.owner = nil

	fmt.Println(user.name, "deposited $20. balance:", user.balance)
	bank.mu.Unlock()
}

func (user *User) withdraw(bank *Bank, amount int) {
	defer wg.Done()
	bank.mu.Lock()
	bank.owner = user

	user.balance += amount
	bank.balance -= amount

	// built-in delay
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

	bank.owner = nil
	bank.mu.Unlock()
}


func runSim(win *pixelgl.Window) {
	startTime := time.Now()

	users := []*User{
		{"Brandon", 100},
		{"Austin", 100},
	}

	bank := Bank{
		balance: 0,
	}

	//table := make([]CustomMutex, len(philosophers))
	// wg.Add(25)
	for i := 0; i < 5; i++ {
		for _, user := range users {
			wg.Add(1)
			go user.deposit(&bank, 20)
			time.Sleep(1 * time.Second)
		}
	}

	// wg.Wait()
	receiveTime := time.Now()

	events["Main"] = [2]time.Time{startTime, receiveTime}
	fmt.Println("length of events", len(events))

	animateChannel(win)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Bank",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear((colornames.White))
	runSim(win)

}

func main() {
	pixelgl.Run(run)
}
