package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var a fyne.App = app.New()
var w fyne.Window = a.NewWindow("Circle")

var forks []fyne.CanvasObject = make([]fyne.CanvasObject, 0)

var philosophers []fyne.CanvasObject = make([]fyne.CanvasObject, 0)

var nameContent []fyne.CanvasObject = make([]fyne.CanvasObject, 0)

var philosopherNames = make([]string, 0)
var nameToIndex = make(map[string]int, 0)

func makePhilosophserTable() {
	circle := canvas.NewCircle(color.White)
	circle.StrokeColor = color.Gray{0x99}
	circle.StrokeWidth = 3
	circle.Resize(fyne.NewSize(600, 600))
	circle.Move(fyne.NewPos(20, 20))

	fillPhilosophers(5)
	fillForks(5)

	philosophersContainer := container.NewWithoutLayout(philosophers...)
	forksContainer := container.NewWithoutLayout(forks...)
	nameContentContainer := container.NewWithoutLayout(nameContent...)

	containerFinal := container.NewWithoutLayout(circle, philosophersContainer, forksContainer, nameContentContainer)
	w.SetContent(containerFinal)
}

func fillForks(n int) {
	const radius = 300
	const centerX, centerY = 300, 300

	angle := (2 * math.Pi) / float64(n) // Angle for each philosopher
	offset := angle / 2                 // Half the angle between philosophers

	for i := 0; i < n; i++ {
		// Calculate position for each icon
		x := centerX + radius*math.Cos(float64(i)*angle+offset)
		y := centerY + radius*math.Sin(float64(i)*angle+offset)

		icon := canvas.NewImageFromResource(theme.DocumentCreateIcon())
		icon.Resize(fyne.NewSize(40, 40))
		icon.Move(fyne.NewPos(float32(x), float32(y)))
		forks = append(forks, icon)
	}
}

func fillPhilosophers(n int) {
	const radius = 300
	const centerX, centerY = 300, 300

	// Calculate angle step between each icon
	angle := (2 * math.Pi) / float64(n)

	for i := 0; i < n; i++ {
		// Calculate position for each icon
		x := centerX + radius*math.Cos(float64(i)*angle)
		y := centerY + radius*math.Sin(float64(i)*angle)

		name := philosopherNames[i]
		nameText := widget.NewLabel(name)

		icon := canvas.NewImageFromResource(theme.AccountIcon())
		icon.Resize(fyne.NewSize(40, 40))
		icon.Move(fyne.NewPos(float32(x), float32(y)))

		nameText.Resize(fyne.NewSize(40, 40))
		nameText.Move(fyne.NewPos(float32(x), float32(y)-20))

		philosophers = append(philosophers, icon)
		nameContent = append(nameContent, nameText)
	}

}

type Philosopher struct {
	name  string // name of philosopher
	count int    // count number of times they eat
	left  int    // fork number on the left
	right int    // fork number on the right
}

func (p *Philosopher) Think() {
	fmt.Println(p.name, "is thinking.")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println(p.name, "is done thinking.")
}

func (p *Philosopher) Eat() {
	p.count++
	fmt.Printf("%s is eating round: %d\n", p.name, p.count)
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println(p.name, "is done eating.")
}

func (p *Philosopher) Dine(table []sync.Mutex) {
	for {
		p.Think()
		table[p.left].Lock()
		fmt.Printf("%s picks up fork %d.\n", p.name, p.left)
		runtime.Gosched() // hack to yield to the next goroutine
		table[p.right].Lock()
		fmt.Printf("%s picks up %d.\n", p.name, p.right)
		p.Eat()
		addAnimationEvent(p, "eat")
		table[p.right].Unlock()
		fmt.Printf("%s puts down fork %d.\n", p.name, p.right)
		table[p.left].Unlock()
		fmt.Printf("%s puts down fork %d.\n", p.name, p.left)
	}
}

func addAnimationEvent(p *Philosopher, action string) {
	index := nameToIndex[p.name]

	canvas.NewAnimation(red, blue, time.Second*2, func(c color.Color) {
		obj.FillColor = c
		canvas.Refresh(obj)
	}).Start()

}

func main() {
	// create a wait group so main won't end
	//var wg sync.WaitGroup
	// wg.Add(1)

	philosophers := []*Philosopher{
		&Philosopher{"Michelle", 0, 0, 1},
		&Philosopher{"Bill", 0, 1, 2},
		&Philosopher{"Sonia", 0, 2, 3},
		&Philosopher{"Brooke", 0, 3, 4},
		&Philosopher{"Eric", 0, 4, 0},
	}

	for index, philosopher := range philosophers {
		philosopherNames = append(philosopherNames, philosopher.name)
		nameToIndex[philosopher.name] = index
	}

	makePhilosophserTable()

	//table := make([]sync.Mutex, len(philosophers))

	// for _, philosopher := range philosophers {
	// 	go func(p *Philosopher) {
	// 		p.Dine(table)
	// 	}(philosopher)
	// }

	// always waits
	// wg.Wait()

	w.Resize(fyne.NewSize(1000, 1000))
	w.ShowAndRun()
}
