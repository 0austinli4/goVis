// Program with race condition fixed by mutex
package main

import (
	"fmt"
	"math/rand"
	"sync" // to import sync later on
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var wg sync.WaitGroup
var globalLock sync.Mutex

type Bank struct {
	owner   *User
	mu      sync.Mutex
	balance int
}

type User struct {
	name    string
	balance int
}

func (user *User) deposit(bank *Bank, amount int) {
	funcTime := time.Now()
	defer wg.Done()

	currentUser := bank.owner

	bank.mu.Lock()
	acquiredTime := time.Now()
	bank.owner = user

	if user.balance >= amount {
		user.balance -= amount
		bank.balance += amount

		fmt.Println(user.name, "deposited $"+fmt.Sprint(amount))
	} else {
		fmt.Println(user.name, "does not have enough to deposit $"+
			fmt.Sprint(amount))
	}

	// built-in delay
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	bank.owner = nil

	bank.mu.Unlock()
	funcEndTime := time.Now()

	globalLock.Lock()
	if currentUser != nil && currentUser.name != user.name {
		key := "ID: " + fmt.Sprint(getGID()) + " " + user.name + ", waiting for " + currentUser.name
		events[key] = [2]time.Time{funcTime, acquiredTime}
		eventsOrder = append(eventsOrder, key)
	}
	key := "ID: " + fmt.Sprint(getGID()) + fmt.Sprint(user.name) + ", deposit"
	events[key] = [2]time.Time{funcTime, funcEndTime}
	eventsOrder = append(eventsOrder, key)

	globalLock.Unlock()

}

func (user *User) withdraw(bank *Bank, amount int) {
	funcTime := time.Now()
	defer wg.Done()

	currentUser := bank.owner

	bank.mu.Lock()

	acquiredTime := time.Now()
	bank.owner = user

	if bank.balance >= amount {
		user.balance += amount
		bank.balance -= amount

		fmt.Println(user.name, "withdrew $"+fmt.Sprint(amount))
	} else {
		fmt.Println("bank does not have enough to withdraw $" +
			fmt.Sprint(amount))
	}

	// built-in delay
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	bank.owner = nil
	bank.mu.Unlock()

	funcEndTime := time.Now()

	globalLock.Lock()

	if currentUser != nil && currentUser.name != user.name {
		key := "ID: " + fmt.Sprint(getGID()) + " " + user.name + ", waiting for " + currentUser.name
		events[key] = [2]time.Time{funcTime, acquiredTime}
		eventsOrder = append(eventsOrder, key)
	}

	key := "ID: " + fmt.Sprint(getGID()) + fmt.Sprint(user.name) + ", withdraw"
	events[key] = [2]time.Time{funcTime, funcEndTime}
	eventsOrder = append(eventsOrder, key)

	globalLock.Unlock()
}

func runSim() {
	startTime := time.Now()

	users := []*User{
		{"Brandon", 40},
		{"Austin", 40},
	}

	bank := Bank{
		balance: 0,
	}

	for i := 0; i < 3; i++ {
		for _, user := range users {
			wg.Add(1)
			go user.deposit(&bank, 20)
			time.Sleep(500 * time.Millisecond)
		}
	}

	for i := 0; i < 3; i++ {
		for _, user := range users {
			wg.Add(1)
			go user.withdraw(&bank, 20)
			time.Sleep(500 * time.Millisecond)
		}
	}

	wg.Wait()
	receiveTime := time.Now()

	events["Main"] = [2]time.Time{startTime, receiveTime}
	eventsOrder = append(eventsOrder, "Main")

	fmt.Println("length of events", len(events))

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
	animateChannel(win)
}

func main() {
	pixelgl.Run(runSim)
}
