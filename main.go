package main

import (
	"bytes"
	"fmt"
	"image/color"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")
	idName := ""
	w.SetPadded(false)
	w.Resize(fyne.NewSize(1000, 1000))

	messages := make(chan string)
	content := container.NewStack()

	go func() {
		messages <- "ping"
		idName = "ping" + fmt.Sprint(getGID())
		channelTransition(content, w, idName, 0, 100)
	}()

	msg := <-messages
	fmt.Println(msg)
	// main method channel
	idName = "main" + fmt.Sprint(getGID())
	channelTransition(content, w, idName, 0, 500)

	w.SetContent(content)

	w.ShowAndRun()
}

func channelTransition(content *fyne.Container, w fyne.Window, goid string, channelX float32, channelY float32) {
	// define colors
	red := color.NRGBA{R: 0xff, A: 0xff}
	blue := color.NRGBA{B: 0xff, A: 0xff}
	// create moving piece ("Data")
	data := canvas.NewRectangle(color.Black)
	data.Resize(fyne.NewSize(50, 50))

	// Create channel
	channel := canvas.NewRectangle(blue)
	channel.Resize(fyne.NewSize(800, 100))        // Set the size of the rectangle
	channel.Move(fyne.NewPos(channelX, channelY)) // Adjust the x, y coordinates as needed

	// name of goRoutine
	channelName := canvas.NewText(goid, color.Black)
	channelName.Move(fyne.NewPos(channelX, channelY+100))

	// a container has channel, data, and channel name
	container := container.NewWithoutLayout(channel, data, channelName)

	canvas.NewColorRGBAAnimation(red, blue, time.Second*2, func(c color.Color) {
		move := canvas.NewPositionAnimation(fyne.NewPos(channelX, channelY), fyne.NewPos(channelX+800, channelY), time.Second*10, data.Move)
		move.AutoReverse = false
		move.Start()
		canvas.Refresh(data)
	}).Start()

	content.Add(container)
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
