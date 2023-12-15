package main

import (
	"bytes"
	"fmt"
	_ "image/png"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var events map[string][2]time.Time = make(map[string][2]time.Time)
var eventsOrder []string = make([]string, 0)

var eventChannelsS []pixel.Vec = make([]pixel.Vec, 0)
var eventChannelsF []pixel.Vec = make([]pixel.Vec, 0)

var eventLocksS []pixel.Vec = make([]pixel.Vec, 0)
var eventLocksF []pixel.Vec = make([]pixel.Vec, 0)

var eventsText map[string]pixel.Vec = make(map[string]pixel.Vec, 0)
var names []pixel.Vec = make([]pixel.Vec, 0)

var namesOfRoutines []string = make([]string, 0)

func channel(win *pixelgl.Window, message string, rows int, index int, startPoint time.Duration,
	endPoint time.Duration, scale time.Duration) {

	offsetY := (win.Bounds().Max.Y/float64(rows))*float64(index) + 70

	if message == "Main" {
		offsetY = 50
	}

	// beginning of rectangle is startPoint
	startRectX := (float64(startPoint)/float64(scale))*win.Bounds().Max.X + 20
	length := (float64(endPoint-startPoint)/float64(scale))*win.Bounds().Max.X - 50

	// limit the bounds
	if startRectX > win.Bounds().Max.X-50 {
		startRectX = win.Bounds().Max.X - 50
	}

	if length < 10 && strings.Contains(message, "waiting") {
		return
	} else if length < 10 {
		length = 10
	}

	// Draw a rectangle
	if strings.Contains(message, "waiting") {
		eventLocksS = append(eventLocksS, pixel.V(startRectX, offsetY))
		eventLocksF = append(eventLocksF, pixel.V(startRectX+length, offsetY+10))
	} else {
		eventChannelsS = append(eventChannelsS, pixel.V(startRectX, offsetY))
		eventChannelsF = append(eventChannelsF, pixel.V(startRectX+length, offsetY+10))
	}

	// draw name of event
	namesPos := pixel.V(startRectX, offsetY)
	names = append(names, namesPos)

	// list of names of goid
	namesOfRoutines = append(namesOfRoutines, message)

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
		lockCounter := 0
		channelCounter := 0

		for i := range namesOfRoutines {
			key := namesOfRoutines[i]

			if key == "Main" && eventsText["Main"].X < win.Bounds().Max.X-30 {
				xPos := eventsText[key]
				xPos.X += 2 // Adjust speed as needed
				eventsText[key] = xPos
			}

			if strings.Contains(key, "waiting") {
				if i > 0 && eventsText["Main"].X > eventLocksF[lockCounter].X {
					imd := imdraw.New(nil)
					imd.Color = colornames.Lawngreen
					imd.Push(eventLocksS[lockCounter], eventLocksF[lockCounter])
					imd.Rectangle(0)
					imd.Draw(batchRect)

					stop := pixel.V(eventLocksF[lockCounter].X, eventChannelsS[channelCounter].Y-20)
					start := pixel.V(eventLocksF[lockCounter].X, eventChannelsS[channelCounter-1].Y)
					arrow(imd, win, start, stop)

					textPos := pixel.V(eventLocksF[lockCounter].X+20, eventLocksF[lockCounter].Y+20)
					nameText := text.New(textPos, basicAtlas)
					nameText.Color = colornames.Purple
					fmt.Fprintln(nameText, "Giving lock")
					nameText.Draw(win, pixel.IM)
				}
				lockCounter += 1
			} else {
				imd := imdraw.New(nil)
				if strings.Contains(key, "deposit") {
					imd.Color = colornames.Cadetblue
				} else {
					imd.Color = colornames.Orange
				}

				imd.Push(eventChannelsS[channelCounter], eventChannelsF[channelCounter])
				imd.Rectangle(0)
				imd.Draw(batchRect)
				channelCounter += 1
			}
		}
		batchRect.Draw(win)

		// draw text
		for i := range namesOfRoutines {
			key := namesOfRoutines[i]

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

	for i := range eventsOrder {
		key := eventsOrder[i]

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

func arrow(imd *imdraw.IMDraw, win *pixelgl.Window, p2 pixel.Vec, p1 pixel.Vec) {
	// Calculate arrow direction and length
	dir := p2.Sub(p1)

	// Draw a line from p1 to p2
	imd.Color = colornames.Black
	imd.Push(p1, p2)
	imd.Line(3)

	// Calculate the arrowhead points
	angle := math.Pi / 6 // Angle of the arrowhead
	arrowLen := 8.0      // Length of the arrowhead

	arrowDir := dir.Unit()
	arrowP1 := p1.Sub(arrowDir.Scaled(arrowLen))                                         // Tip of the arrow
	arrowP2 := p1.Sub(arrowDir.Rotated(math.Pi / 2).Scaled(arrowLen * math.Tan(angle)))  // Left side of the arrowhead
	arrowP3 := p1.Sub(arrowDir.Rotated(-math.Pi / 2).Scaled(arrowLen * math.Tan(angle))) // Right side of the arrowhead

	// Draw the arrowhead by connecting the three points
	imd.Push(arrowP1, arrowP2)
	imd.Line(3)
	imd.Push(arrowP1, arrowP3)
	imd.Line(3)

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
