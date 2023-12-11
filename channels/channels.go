package main

import (
	"bytes"
	"fmt"
	_ "image/png"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var events map[string][2]time.Time = make(map[string][2]time.Time)

var eventChannelsS []pixel.Vec = make([]pixel.Vec, 0)
var eventChannelsF []pixel.Vec = make([]pixel.Vec, 0)

var eventsText map[string]pixel.Vec = make(map[string]pixel.Vec, 0)
var names []pixel.Vec = make([]pixel.Vec, 0)

var (
	// Create a global mutex
	mutex sync.Mutex
)

func channel(win *pixelgl.Window, message string, rows int, index int, startPoint time.Duration, endPoint time.Duration, scale time.Duration) {
	// length of the rectangle is determined by endPoint  - startPoint
	offsetY := (win.Bounds().Max.Y / float64(rows)) * float64(index)

	// beginning of rectangle is startPoint
	startRect := (float64(startPoint) / float64(scale)) * win.Bounds().Max.X
	length := (float64(endPoint-startPoint) / float64(scale)) * win.Bounds().Max.X

	// Draw a rectangle
	eventChannelsS = append(eventChannelsS, pixel.V(startRect, win.Bounds().Max.Y-float64(index)*offsetY-20))
	eventChannelsF = append(eventChannelsF, pixel.V(startRect+length, win.Bounds().Max.Y-float64(index)*offsetY-40))

	// draw name of event
	namesPos := pixel.V(startRect, win.Bounds().Max.Y-float64(index)*offsetY-40)
	names = append(names, namesPos)

	// Create a basic font
	textPos := pixel.V(startRect, win.Bounds().Max.Y-float64(index)*offsetY-40)
	eventsText[message] = textPos
}

func animateAll(win *pixelgl.Window) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	win.Update()

	for !win.Closed() {
		win.Clear(colornames.White)
		// Draw the rectangle
		i := 0
		for key := range events {
			imd := imdraw.New(nil)
			imd.Color = colornames.Green
			imd.Push(eventChannelsS[i], eventChannelsF[i])
			imd.Rectangle(0)
			imd.Draw(win)

			xPos := eventsText[key]
			xPos.X += 4 // Adjust speed as needed
			eventsText[key] = xPos

			if eventsText["Main"].X > eventChannelsF[i].X {
				lineDown(win, eventChannelsF[i])
				eventChannelsF[i].X = 0
			}

			basicTxt := text.New(names[i], basicAtlas)
			basicTxt.Color = colornames.Black
			fmt.Fprintln(basicTxt, key)
			basicTxt.Draw(win, pixel.IM)

			basicTxt = text.New(eventsText[key], basicAtlas)
			basicTxt.Color = colornames.Black
			fmt.Fprintln(basicTxt, key)
			basicTxt.Draw(win, pixel.IM)
			i += 1
		}
		// Draw the text at the updated position
		win.Update()
	}
}

func animateChannel(win *pixelgl.Window) {
	universalStartTime := events["Main"][0]
	universalEndTime := events["Main"][1]

	scale := universalEndTime.Sub(universalStartTime)
	rows := len(events)

	i := 0

	for key := range events {
		startTime := events[key][0]
		startPoint := startTime.Sub(universalStartTime)
		endTime := events[key][1]
		endPoint := endTime.Sub(universalStartTime)
		fmt.Println("start point", startTime, "endPoint", endTime, "Scale: ", scale)
		fmt.Println("printing channel")
		channel(win, key, rows, i, startPoint, endPoint, scale)
		i += 1
	}
	animateAll(win)
}

func lineDown(win *pixelgl.Window, pos pixel.Vec) {
	fmt.Println(pos)
	imd := imdraw.New(nil)
	imd.Color = colornames.Black
	imd.Clear()
	imd.Push(pos, pixel.V(pos.X, 0)) // Draw a line from the current position down to Y=0
	imd.Line(2)
	imd.Draw(win)
}

func runChannel(win *pixelgl.Window) {
	startTime := time.Now()
	messages := make(chan string)

	go func() {
		funcTime := time.Now()
		messages <- "ping"
		funcEndTime := time.Now()
		events["Program 1"] = [2]time.Time{funcTime, funcEndTime}
	}()

	msg := <-messages
	receiveTime := time.Now()
	fmt.Println(msg)
	events["Main"] = [2]time.Time{startTime, receiveTime}

	fmt.Println(events)

	animateChannel(win)
}

// func calcSquares(number int, squareop chan int) {
// 	funcTime := time.Now()
// 	sum := 0
// 	for number != 0 {
// 		digit := number % 10
// 		sum += digit * digit
// 		number /= 10
// 	}
// 	time.Sleep(5 * time.Second)
// 	squareop <- sum
// 	funcEndTime := time.Now()
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	events["calcSquares"+fmt.Sprint(getGID())] = [2]time.Time{funcTime, funcEndTime}
// }

// func calcCubes(number int, cubeop chan int) {
// 	funcTime := time.Now()
// 	time.Sleep(2 * time.Second)
// 	sum := 0
// 	for number != 0 {
// 		digit := number % 10
// 		sum += digit * digit * digit
// 		number /= 10
// 	}
// 	cubeop <- sum
// 	funcEndTime := time.Now()
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	events["calcCubes"+fmt.Sprint(getGID())] = [2]time.Time{funcTime, funcEndTime}
// }

// func runChannel(win *pixelgl.Window) {
// 	startTime := time.Now()

// 	number := 589
// 	sqrch := make(chan int)
// 	cubech := make(chan int)
// 	go calcSquares(number, sqrch)
// 	go calcCubes(number, cubech)
// 	squares, cubes := <-sqrch, <-cubech
// 	fmt.Println("Final output", squares+cubes)

// 	receiveTime := time.Now()

// 	events["Main"+fmt.Sprint(getGID())] = [2]time.Time{startTime, receiveTime}
// 	animateChannel(win)
// }

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

	// for !win.Closed() {
	// 	win.Update()
	// }
	runChannel(win)
}

func main() {
	pixelgl.Run(run)
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
