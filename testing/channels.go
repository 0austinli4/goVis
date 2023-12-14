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

	offsetY := (win.Bounds().Max.Y / float64(rows)) * float64(index)
	if message == "Main" {
		offsetY = 50
	}
	if offsetY < 50 {
		offsetY = 50 + 20*float64(index+1)
	}
	// fmt.Println("CURRENT OFFSET", offsetY)

	// beginning of rectangle is startPoint
	startRectX := (float64(startPoint)/float64(scale))*win.Bounds().Max.X + 20
	length := (float64(endPoint-startPoint)/float64(scale))*win.Bounds().Max.X - 50
	// limit the bounds
	if startRectX > win.Bounds().Max.X-50 {
		startRectX = win.Bounds().Max.X - 50
	}

	if length < 20 {
		length = 20
	}

	// Draw a rectangle
	eventChannelsS = append(eventChannelsS, pixel.V(startRectX, offsetY))
	eventChannelsF = append(eventChannelsF, pixel.V(startRectX+length, offsetY+10))

	// draw name of event
	namesPos := pixel.V(startRectX, offsetY)
	names = append(names, namesPos)

	// text pos
	textPos := pixel.V(startRectX, offsetY)
	eventsText[message] = textPos
}

func animateAll(win *pixelgl.Window) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	win.Update()

	for !win.Pressed(pixelgl.MouseButtonLeft) {
		batchRect := pixel.NewBatch(&pixel.TrianglesData{}, nil)
		batchRect.Clear()
		win.Clear(colornames.Aliceblue)

		// draw channels
		i := 0

		for key := range events {

			if key == "Main" {
				xPos := eventsText[key]
				xPos.X += 2 // Adjust speed as needed
				eventsText[key] = xPos
			}

			imd := imdraw.New(nil)
			imd.Color = colornames.Cadetblue
			imd.Push(eventChannelsS[i], eventChannelsF[i])
			imd.Rectangle(0)
			imd.Draw(batchRect)

			if eventsText["Main"].X > eventChannelsS[i].X {
				imd := imdraw.New(nil)
				lineUp(imd, win, eventChannelsS[i], i)
			}

			if eventsText["Main"].X > eventChannelsF[i].X {
				imd := imdraw.New(nil)
				lineDown(imd, win, eventChannelsF[i])
				// imd.Color = colornames.Cadetblue
				// imd.Push(eventChannelsS[i], eventChannelsF[i])
				// imd.Rectangle(0)
				// imd.Draw(batchRect)
			}
			i += 1
		}
		batchRect.Draw(win)

		// draw text
		i = 0
		for key := range events {
			basicTxt := text.New(eventsText[key], basicAtlas)
			if key == "Main" {
				nameText := text.New(pixel.V(10, 50), basicAtlas)
				nameText.Color = colornames.Black
				fmt.Fprintln(nameText, key)
				nameText.Draw(win, pixel.IM)
			}
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

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
