package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gopxl/pixel/pixelgl"
)

var wg sync.WaitGroup

func runChannelHelloWorld(win *pixelgl.Window) {
	startTime := time.Now()
	messages := make(chan string)

	funcTime := time.Now()
	go func() {

		messages <- "ping"
	}()
	funcEndTime := time.Now()
	events["Program 1"] = [2]time.Time{funcTime, funcEndTime}

	msg := <-messages
	fmt.Println(msg)

	funcTime = time.Now()
	go func() {
		messages <- "ping"
	}()
	funcEndTime = time.Now()
	events["Program 2"] = [2]time.Time{funcTime, funcEndTime}
	msg = <-messages
	fmt.Println(msg)

	funcTime = time.Now()
	go func() {
		messages <- "ping"
	}()
	msg = <-messages
	funcEndTime = time.Now()
	events["Program 3"] = [2]time.Time{funcTime, funcEndTime}

	receiveTime := time.Now()
	fmt.Println(msg)
	events["Main"] = [2]time.Time{startTime, receiveTime}

	fmt.Println(events)

	animateChannel(win)
}

func ponger(c chan string) {
	for i := 0; i < 2; i++ {
		funcTime := time.Now()
		c <- "pong"
		funcEndTime := time.Now()
		events["Ponger"+fmt.Sprintln(i)+fmt.Sprintln(getGID())] = [2]time.Time{funcTime, funcEndTime}
	}

}

func pinger(c chan string) {
	for i := 0; i < 2; i++ {
		funcTime := time.Now()
		c <- "ping"
		funcEndTime := time.Now()
		events["Pinger "+fmt.Sprintln(i)+fmt.Sprintln(getGID())] = [2]time.Time{funcTime, funcEndTime}
	}
}

func printer(c chan string) {
	for i := 0; i < 4; i++ {
		funcTime := time.Now()
		msg := <-c
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
		funcEndTime := time.Now()
		events["Printer "+fmt.Sprintln(getGID())] = [2]time.Time{funcTime, funcEndTime}
	}
	wg.Done()
}

func runPingPong(win *pixelgl.Window) {
	startTime := time.Now()
	var c chan string = make(chan string)
	wg.Add(1)
	go pinger(c)
	go ponger(c)
	go printer(c)
	wg.Wait()

	// var input string
	// fmt.Scanln(&input)
	fmt.Println("attempintg to pirnt main")
	receiveTime := time.Now()
	events["Main"] = [2]time.Time{startTime, receiveTime}
	animateChannel(win)
}

// func run() {
// 	cfg := pixelgl.WindowConfig{
// 		Title:  "Pixel Rocks!",
// 		Bounds: pixel.R(0, 0, 1024, 768),
// 		VSync:  true,
// 	}

// 	win, err := pixelgl.NewWindow(cfg)
// 	if err != nil {
// 		panic(err)
// 	}
// 	win.Clear((colornames.White))

// 	runChannelHelloWorld(win)
// 	runPingPong(win)
// 	runPingPong(win)
// }

// func main() {
// 	pixelgl.Run(run)
// }
