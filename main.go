package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var a fyne.App = app.New()
var w fyne.Window = a.NewWindow("Hello")
var containers []*fyne.Container = make([]*fyne.Container, 0)
var startTime time.Time
var currTime = float32(0.0)

var i int = 0
var idName string

type Groups struct {
	mu       sync.Mutex
	counters map[string]int
}

func (c *Groups) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
	idName = "name" + fmt.Sprint(i) + fmt.Sprint(getGID())
	channelTransition(idName)
}

func main() {
	//idName := ""
	w.SetPadded(false)
	w.Resize(fyne.NewSize(1000, 1000))

	//wg.Add(2)
	// go inc()
	// go inc()
	// wg.Wait()
	c := Groups{
		counters: map[string]int{"a": 0, "b": 0},
	}

	var wg sync.WaitGroup

	doIncrement := func(name string, n int) {
		for i := 0; i < n; i++ {
			c.inc(name)
		}
		wg.Done()
	}

	wg.Add(6)
	startTime = time.Now() // Get the current time when the function starts
	fmt.Println(startTime)
	go doIncrement("a", 1)
	go doIncrement("a", 1)
	go doIncrement("b", 1)
	go doIncrement("a", 1)
	go doIncrement("a", 1)
	go doIncrement("b", 1)

	wg.Wait()
	fmt.Println(c.counters)

	// idName = "main" + fmt.Sprint(i) + fmt.Sprint(getGID())
	// channelTransition(idName)

	stack := container.New(layout.NewGridLayoutWithRows(6))

	for _, obj := range containers {
		stack.Add(obj)
	}
	w.SetContent(stack)
	w.ShowAndRun()
}

// func inc() {
// 	i = i + 1
// 	idName = "ping" + fmt.Sprint(i) + fmt.Sprint(getGID())
// 	channelTransition(idName)
// 	wg.Done()
// }

func channelTransition(goid string) {
	elapsed := time.Since(startTime)
	blockOffset := (float32(elapsed.Nanoseconds()) / float32(time.Second.Nanoseconds())) * 1000000 // Fraction of elapsed time in seconds
	// print(blockOffset)
	// define colors
	// red := color.NRGBA{R: 0xff, A: 0xff}
	randColor1 := color.RGBA{
		R: uint8(rand.Intn(256)), // Random Red component
		G: uint8(rand.Intn(256)), // Random Green component
		B: uint8(rand.Intn(256)), // Random Blue component
		A: 0xff,                  // Alpha value (fully opaque)
	}
	// create moving piece ("Data")
	// Create channel
	randColor2 := color.RGBA{
		R: uint8(rand.Intn(256)), // Random Red component
		G: uint8(rand.Intn(256)), // Random Green component
		B: uint8(rand.Intn(256)), // Random Blue component
		A: 0xff,                  // Alpha value (fully opaque)
	}
	channel := canvas.NewRectangle(randColor1)
	channel.Move(fyne.NewPos(blockOffset, 0))

	square := canvas.NewRectangle(randColor2)
	square.Resize(fyne.NewSize(30, 30))
	channelName := canvas.NewText(goid, color.Black)
	channelName.TextSize = 20
	data := container.NewWithoutLayout(square, channelName)
	data.Move(fyne.NewPos(blockOffset, 0))

	channel.Resize(fyne.NewSize(blockOffset, 20)) // Set the size of the rectangle

	//pauseDuration := time.Duration(currTime * float32(time.Second))

	container := container.NewWithoutLayout(channel, data)

	// delayAnimation := func() {
	// 	time.AfterFunc(pauseDuration, func() {
	// 		// Second animation after the delay
	// 		move := canvas.NewPositionAnimation(fyne.NewPos(container.Position().X+blockOffset, container.Position().Y), fyne.NewPos(container.Position().X+blockOffset*2, container.Position().Y), time.Second*4, data.Move)
	// 		move.AutoReverse = false
	// 		move.Start()
	// 		canvas.Refresh(data)
	// 	})
	// }
	// delayAnimation()
	move := canvas.NewPositionAnimation(fyne.NewPos(container.Position().X+blockOffset, container.Position().Y), fyne.NewPos(container.Position().X+blockOffset*2, container.Position().Y), time.Second*4, data.Move)
	move.AutoReverse = false
	move.Start()
	canvas.Refresh(data)

	currTime = currTime + blockOffset/100
	containers = append(containers, container)

}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
