package main

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

var a fyne.App = app.New()
var w fyne.Window = a.NewWindow("Hello")

// var containers []contaer
var events []event = make([]event, 0)
var containers []*fyne.Container = make([]*fyne.Container, 0)
var moveList []*fyne.Animation = make([]*fyne.Animation, 0)
var animationList []func() = make([]func(), 0)
var startTime time.Time
var endTime time.Time
var creationTime time.Time
var ogTime time.Time

// type Groups struct {
// 	mu       sync.Mutex
// 	counters map[string]int
// }

type event struct {
	ID           uint64
	Name         string
	CreationTime time.Time
	StartTime    time.Time
	EndTime      time.Time
}

// func (c *Groups) inc(name string) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	startTime = time.Now()
// 	eventID := getGID()
// 	c.counters[name]++
// 	eventName := "Inc " + name + fmt.Sprint(eventID)
// 	endTime = time.Now()
// 	newEvent := event{
// 		ID:           eventID,
// 		Name:         eventName,
// 		CreationTime: creationTime,
// 		StartTime:    startTime,
// 		EndTime:      endTime,
// 	}
// 	events = append(events, newEvent)

// }

// func sum(s []int, c chan int) {
// 	startTime = time.Now()
// 	sum := 0
// 	for _, v := range s {
// 		sum += v
// 	}
// 	id := getGID()
// 	endTime = time.Now()
// 	addNewEvent(id, "Calculated"+fmt.Sprint(sum)+", id: "+fmt.Sprint(id), creationTime, startTime, endTime)
// 	c <- sum // send sum to c
// }

func initializeWindow() {
	w.SetPadded(false)
	w.Resize(fyne.NewSize(800, 800))
	ogTime = time.Now()
}

func main() {
	initializeWindow()
	say("world")
	say("hello")
	// s := []int{7, 2, 8, -9, 4, 0}

	// c := make(chan int)
	// go sum(s[:len(s)/2], c)
	// go sum(s[len(s)/2:], c)
	// x, y := <-c, <-c // receive from c
	// fmt.Println(x, y, x+y)
	funcEndTime := time.Now()
	showandRun(funcEndTime)
}

func showandRun(funcEndTime time.Time) {
	addNewEvent(getGID(), "main", ogTime, ogTime, funcEndTime)

	totalDuration := funcEndTime.Sub(ogTime)
	scale := float32(totalDuration.Nanoseconds())

	stack := container.New(layout.NewGridLayoutWithRows(len(events)))

	for _, event := range events {
		createContainer(event, scale)
	}
	for _, obj := range containers {
		stack.Add(obj)
	}
	w.SetContent(stack)

	for _, animation := range animationList {
		animation()
	}
	for _, move := range moveList {
		move.Start()
	}
	// visualizeLocks()
	w.ShowAndRun()
}

func say(s string) {
	startTime = time.Now()
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
	id := getGID()
	endTime = time.Now()
	addNewEvent(id, "s: "+fmt.Sprint(id), creationTime, startTime, endTime)
}

func addNewEvent(eventID uint64, eventName string, creationTime time.Time, startTime time.Time, endTime time.Time) {
	newEvent := event{
		ID:           eventID,
		Name:         eventName,
		CreationTime: creationTime,
		StartTime:    startTime,
		EndTime:      endTime,
	}

	events = append(events, newEvent)
}

// func main() {
// 	var wg sync.WaitGroup

// 	//idName := ""
// 	w.SetPadded(false)
// 	w.Resize(fyne.NewSize(3000, 2000))

// 	c := Groups{
// 		counters: map[string]int{"a": 0, "b": 0},
// 	}

// 	doIncrement := func(name string, n int) {
// 		for i := 0; i < n; i++ {
// 			c.inc(name)
// 		}
// 		wg.Done()
// 	}
// 	ogTime = time.Now()
// 	wg.Add(6)
// 	// idName = "MAIN" + fmt.Sprint(i) + fmt.Sprint(getGID())
// 	// channelTransition(idName)

// 	creationTime = time.Now() // Get the current time when the function starts
// 	go doIncrement("a", 1)
// 	//creationTime = time.Now() // Get the current time when the function starts
// 	go doIncrement("a", 1)
// 	//creationTime = time.Now() // Get the current time when the function starts
// 	go doIncrement("b", 1)
// 	//creationTime = time.Now() // Get the current time when the function starts
// 	go doIncrement("a", 1)
// 	//creationTime = time.Now() // Get the current time when the function starts
// 	go doIncrement("a", 1)
// 	//creationTime = time.Now() // Get the current time when the function starts
// 	go doIncrement("b", 1)

// 	wg.Wait()

// 	mainEvent := event{
// 		ID:           getGID(),
// 		Name:         "main",
// 		CreationTime: ogTime,
// 		StartTime:    ogTime,
// 		EndTime:      time.Now(),
// 	}
// 	events = append(events, mainEvent)

// 	stack := container.New(layout.NewGridLayoutWithRows(len(events)))

// 	for _, event := range events {
// 		createContainer(event)
// 	}

// 	//print(containers)

//		for _, obj := range containers {
//			// fmt.Println(obj)
//			stack.Add(obj)
//		}
//		w.SetContent(stack)
//		for _, move := range moveList {
//			move.Start()
//		}
//		w.ShowAndRun()
//	}

func visualizeLocks() {
	lockImage := canvas.NewImageFromResource(theme.MediaRecordIcon())
	lockImage.FillMode = canvas.ImageFillOriginal

	lockContainer := container.New(layout.NewVBoxLayout(),
		lockImage,
	)
	w.SetContent(lockContainer)
	// Open and close the lock animation
	go animateLock(lockImage)

}
func animateLock(lock *canvas.Image) {
	// Open and close the lock animation
	for {
		time.Sleep(2 * time.Second)
		lock.Resource = theme.RadioButtonIcon()
		lock.Refresh()
		time.Sleep(2 * time.Second)
		lock.Resource = theme.MediaRecordIcon()
		lock.Refresh()
	}
}

func createContainer(testEvent event, scale float32) {
	red := color.NRGBA{R: 0xcc, A: 0xff}                    // Adjust the R value as needed (here, 0xcc)
	gray := color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff} // Adjust RGB values as needed (here, 0x33 for each)

	name := testEvent.Name       // Accessing Name field
	sTime := testEvent.StartTime // Accessing StartTime field
	eTime := testEvent.EndTime   // Accessing EndTime field

	fmt.Print(eTime, sTime)
	elapsed := eTime.Sub(sTime)
	soTimeDiff := sTime.Sub(ogTime)

	offset := float32(soTimeDiff.Nanoseconds()) / 1000
	length := (float32(elapsed.Nanoseconds()) / scale) * 800
	delay := (time.Duration(offset/100) * 1e9)

	fmt.Println("Offset of block: ", offset)
	fmt.Println("length of Block: ", length)
	fmt.Println("delay for this block", delay)

	channel := canvas.NewRectangle(red)
	channel.Move(fyne.NewPos(float32(offset), 0))
	// fmt.Println("length", )
	channel.Resize(fyne.NewSize(length, 30))

	square := canvas.NewRectangle(gray)
	square.Resize(fyne.NewSize(30, 30))

	channelName := canvas.NewText(name, color.Black)
	channelName.TextSize = 30

	data := container.NewWithoutLayout(square, channelName)
	container := container.NewWithoutLayout(channel, data)

	// w.SetContent(container)
	// w.ShowAndRun()
	//data.Move(fyne.NewPos(blockOffset, 0))
	//channel.Resize(fyne.NewSize(blockOffset, 30)) // Set the size of the rectangle
	containers = append(containers, container)

	startPos := fyne.NewPos(container.Position().X+offset, container.Position().Y)
	// fmt.Println(math.Min(float64(container.Position().X+offset+length), 2500))

	endPos := fyne.NewPos(float32(math.Min(float64(container.Position().X+offset+length), 600)), container.Position().Y)

	//canvas.Refresh()
	delayAnimation := func() {
		time.AfterFunc(delay, func() {
			move := canvas.NewPositionAnimation(startPos, endPos, time.Second*4, data.Move)
			move.Start()
		})
	}
	// delayAnimation()
	// move := canvas.NewPositionAnimation(fyne.NewPos(container.Position().X+blockOffset, container.Position().Y), fyne.NewPos(container.Position().X+blockOffset*2, container.Position().Y), time.Second*4, data.Move)
	// move.AutoReverse = false
	// move.Start()
	// canvas.Refresh(data)
	moveList = append(moveList, canvas.NewPositionAnimation(startPos, endPos, time.Second*6, data.Move))
	animationList = append(animationList, delayAnimation)
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
