package main

import (
	"bytes"
	"fmt"
	_ "image/png"
	"runtime"
	"strconv"
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

// var (
// 	// Create a global mutex
// 	mutex sync.Mutex
// )

func channel(win *pixelgl.Window, message string, rows int, index int, startPoint time.Duration, endPoint time.Duration, scale time.Duration) {
	// length of the rectangle is determined by endPoint  - startPoint
	// offsetY := (win.Bounds().Max.Y / float64(rows)) * float64(index)
	fmt.Println("CURRENT INDEX", index)
	offsetY := float64(50.0 * index)

	// beginning of rectangle is startPoint
	// startRect := 0.0
	startRect := (float64(startPoint) / float64(scale)) * win.Bounds().Max.X
	length := (float64(endPoint-startPoint)/float64(scale))*win.Bounds().Max.X - 1

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
		win.Clear(colornames.Aliceblue)
		// Draw the rectangle
		i := 0
		for key := range events {

			if key == "Main" {
				xPos := eventsText[key]
				xPos.X += 2 // Adjust speed as needed
				eventsText[key] = xPos
			}

			fmt.Println(eventChannelsS[i].X, eventChannelsF[i].X)

			if eventsText["Main"].X > eventChannelsS[i].X {
				imd := imdraw.New(nil)
				lineUp(imd, win, eventChannelsS[i], i)
				imd.Color = colornames.Cadetblue
				imd.Push(eventChannelsS[i], eventChannelsF[i])
				imd.Rectangle(0)
				imd.Draw(win)
			}

			if eventsText["Main"].X > eventChannelsF[i].X {
				imd := imdraw.New(nil)
				lineDown(imd, win, eventChannelsF[i])
				imd.Color = colornames.Cadetblue
				imd.Push(eventChannelsS[i], eventChannelsF[i])
				imd.Rectangle(0)
				imd.Draw(win)
			}

			basicTxt := text.New(eventsText[key], basicAtlas)
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
		if key == "Main" {
			continue
		}
		startTime := events[key][0]
		startPoint := startTime.Sub(universalStartTime)
		endTime := events[key][1]
		endPoint := endTime.Sub(universalStartTime)
		channel(win, key, rows, i, startPoint, endPoint, scale)
		i += 1
	}
	startTime := events["Main"][0]
	startPoint := startTime.Sub(universalStartTime)
	endTime := events["Main"][1]
	endPoint := endTime.Sub(universalStartTime)
	channel(win, "Main", rows, i, startPoint, endPoint, scale)

	animateAll(win)
}

func lineDown(imd *imdraw.IMDraw, win *pixelgl.Window, pos pixel.Vec) {
	imd.Color = colornames.Black
	imd.Push(pos, pixel.V(pos.X, eventChannelsS[len(eventChannelsS)-1].Y))
	imd.Line(2)
	imd.Draw(win)
}

func lineUp(imd *imdraw.IMDraw, win *pixelgl.Window, pos pixel.Vec, index int) {
	imd.Color = colornames.Black
	imd.Push(pixel.V(pos.X, eventChannelsS[len(eventChannelsS)-1].Y), pixel.V(pos.X, eventChannelsF[index].Y))
	imd.Line(2)
	imd.Draw(win)
}

// func worker(done chan bool, win *pixelgl.Window) {
// 	startTime := time.Now()
// 	fmt.Print("working...")
// 	time.Sleep(time.Second)
// 	fmt.Println("done")
// 	done <- true
// 	endTime := time.Now()
// 	events["Program 1"+fmt.Sprint(getGID)] = [2]time.Time{startTime, endTime}
// }

// func runChannel(win *pixelgl.Window) {
// 	startTime := time.Now()
// 	done := make(chan bool, 1)
// 	go worker(done, win)
// 	<-done
// 	receiveTime := time.Now()
// 	events["Main"+fmt.Sprint(getGID)] = [2]time.Time{startTime, receiveTime}
// 	animateChannel(win)
// }

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
	fmt.Println(msg)

	go func() {
		funcTime := time.Now()
		messages <- "ping"
		funcEndTime := time.Now()
		events["Program 2"] = [2]time.Time{funcTime, funcEndTime}
	}()
	msg = <-messages
	fmt.Println(msg)

	go func() {
		funcTime := time.Now()
		messages <- "ping"
		funcEndTime := time.Now()
		events["Program 3"] = [2]time.Time{funcTime, funcEndTime}
	}()
	msg = <-messages
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
