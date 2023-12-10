package main

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	_ "image/png"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Philosopher struct {
	name      string       // name of philosopher
	count     int          // count number of times they eat
	left      int          // fork number on the left
	right     int          // fork number on the right
	gopher    pixel.Sprite // animation
	spritePos pixel.Vec
	eating    bool
}

type Fork struct {
	index   int // name of philosopher
	owner   string
	forkPic pixel.Sprite // animation
	mat     pixel.Matrix
}

var philosophers []*Philosopher = make([]*Philosopher, 0)
var nameToIndex = make(map[string]int)

var forks []*Fork

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// func initializeDrawing(win *pixelgl.Window) {
// 	for i := 0; i < len(philosophers); i++ {
// 		sprite := philosophers[i].gopher
// 		mat := philosophers[i].mat
// 		sprite.Draw(win, mat)

// 		spriteFork := forks[i].forkPic
// 		matFork := forks[i].mat

// 		spriteFork.Draw(win, matFork)
// 	}

// 	for !win.Closed() {
// 		win.Update()
// 	}
// }

func initializePhilosophers(win *pixelgl.Window) {
	pic, err := loadPicture("hiking.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	spritePos := pixel.V(0, 0)

	philosophers = []*Philosopher{
		{"Michelle", 0, 0, 1, *sprite, spritePos, false},
		{"Bill", 0, 1, 2, *sprite, spritePos, false},
		{"Sonia", 0, 2, 3, *sprite, spritePos, false},
		{"Brooke", 0, 3, 4, *sprite, spritePos, false},
		{"Eric", 0, 4, 0, *sprite, spritePos, false},
	}
	i := 0

	for _, p := range philosophers {
		nameToIndex[p.name] = i
		i++
	}
	fmt.Println(nameToIndex)

	centerX := (win.Bounds().Center()).X
	centerY := (win.Bounds().Center()).Y
	radius := 300.0
	numSprites := len(philosophers)

	angleIncrement := (2 * math.Pi) / float64(numSprites)
	initialAngle := math.Pi / 2 // Start from the top

	for i := 0; i < numSprites; i++ {
		angle := initialAngle + float64(i)*angleIncrement
		spritePos := pixel.V(
			centerX+radius*math.Cos(angle),
			centerY+radius*math.Sin(angle),
		)
		mat := pixel.IM
		mat = mat.ScaledXY(spritePos, pixel.V(0.15, 0.15))
		philosophers[i].spritePos = spritePos
		sprite.Draw(win, mat)
	}
}

func initializeForks(win *pixelgl.Window) {
	pic, err := loadPicture("fork.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	mat := pixel.IM

	forks = []*Fork{
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
	}

	centerX := (win.Bounds().Center()).X
	centerY := (win.Bounds().Center()).Y
	radius := 300.0
	numSprites := len(forks)

	angleIncrement := (2 * math.Pi) / float64(numSprites)
	initialAngle := math.Pi/2 + math.Pi/4 // Start from the top

	for i := 0; i < numSprites; i++ {
		angle := initialAngle + float64(i)*angleIncrement
		spritePos := pixel.V(
			centerX+radius*math.Cos(angle),
			centerY+radius*math.Sin(angle),
		)
		mat := pixel.IM
		mat = mat.ScaledXY(spritePos, pixel.V(0.15, 0.15))
		forks[i].mat = mat
		sprite.Draw(win, mat)
	}
}

func drawNewFrame(win *pixelgl.Window) {
	forkPic, _ := loadPicture("fork.png")
	standing, _ := loadPicture("hiking.png")
	eating, _ := loadPicture("gamer.png")

	spriteFork := pixel.NewSprite(forkPic, forkPic.Bounds())
	spriteStand := pixel.NewSprite(standing, standing.Bounds())
	spriteEat := pixel.NewSprite(eating, eating.Bounds())

	// clear frame
	win.Clear(colornames.White)

	for i := 0; i < len(forks); i++ {
		mat := forks[i].mat
		spriteFork.Draw(win, mat)
	}
	// angle := 0.0

	for i := 0; i < len(philosophers); i++ {
		spritePos := philosophers[i].spritePos
		mat := pixel.IM
		if philosophers[i].eating {
			mat = mat.ScaledXY(spritePos, pixel.V(0.18, 0.18))
			spriteEat.Draw(win, mat)
		} else {
			mat = mat.ScaledXY(spritePos, pixel.V(0.15, 0.15))
			spriteStand.Draw(win, mat)
		}
	}
}

func updateEat(index int) {
	philosophers[index].eating = true
}

func drawText(textSeg string, win *pixelgl.Window) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 650), basicAtlas)
	basicTxt.Color = colornames.Black // Set the text color to Red

	fmt.Fprintln(basicTxt, textSeg)
	// basicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	// basicTxt := text.New(pixel.V(300, 800), basicAtlas)

	// fmt.Fprintln(basicTxt, textSeg)
	basicTxt.Draw(win, pixel.IM)
}

func (p *Philosopher) Think() {
	fmt.Println(p.name, "is thinking.")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println(p.name, "is done thinking.")
}

func (p *Philosopher) Eat(win *pixelgl.Window) {
	p.count++
	textStr := (p.name + " is eating round:" + fmt.Sprint(p.count))
	updateEat(nameToIndex[p.name])
	drawNewFrame(win)
	drawText(textStr, win)
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println(p.name, "is done eating.")
}
func (p *Philosopher) Dine(table []sync.Mutex, win *pixelgl.Window) {
	for {
		p.Think()
		table[p.left].Lock()
		fmt.Printf("%s picks up fork %d.\n", p.name, p.left)
		runtime.Gosched() // hack to yield to the next goroutine
		table[p.right].Lock()
		fmt.Printf("%s picks up %d.\n", p.name, p.right)
		p.Eat(win)
		table[p.right].Unlock()
		fmt.Printf("%s puts down fork %d.\n", p.name, p.right)
		table[p.left].Unlock()
		fmt.Printf("%s puts down fork %d.\n", p.name, p.left)
	}
}

func run() {
	// var wg sync.WaitGroup
	// wg.Add(1)

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Aliceblue)

	// table := make([]sync.Mutex, len(philosophers))

	initializeForks(win)
	initializePhilosophers(win)
	//drawTable(win)
	drawText("dafuq", win)

	i := 0
	updateEat(i)
	drawNewFrame(win)
	for !win.Closed() {
		// Existing code...
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			updateEat(i)
			drawNewFrame(win)
			i += 1
		}
		win.Update()
	}
	// for _, philosopher := range philosophers {
	// 	go func(p *Philosopher) {
	// 		p.Dine(table, win)
	// 		// win.Update()
	// 	}(philosopher)
	// }
	// wg.Wait()

	// i := 0
	// for !win.Closed() {
	// 	updateEat(i)
	// 	drawNewFrame(win)
	// 	i += 1
	// 	time.Sleep(5 * time.Second)
	// 	win.Update()
	// }

}

func drawTable(win *pixelgl.Window) {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Data for the table (example)
	tableData := [][]string{
		{"Fork 1", "Philosopher 1"},
		{"Fork 2", "Philosopher 2"},
		{"Fork 3", "Philosopher 3"},
		{"Fork 4", "Philosopher 4"},
	}

	// Calculate text position for the table (top right corner)
	tableStartX := win.Bounds().Max.X - 200 // Adjust the position as needed
	tableStartY := win.Bounds().Max.Y - 30  // Adjust the position as needed

	// Draw the table
	for i, row := range tableData {
		for j := range row {
			textPos := pixel.V(tableStartX+float64(j*50), tableStartY-float64(i*20))
			basicTxt := text.New(textPos, basicAtlas)
			basicTxt.Color = colornames.Black
			fmt.Fprintln(basicTxt, tableData[i][0])

			textPos = pixel.V(tableStartX+float64(j*50)+20, tableStartY-float64(i*20))
			basicTxt = text.New(textPos, basicAtlas)
			fmt.Fprintln(basicTxt, tableData[i][1])
			// basicTxt.Draw(win, pixel.IM.Scaled(textPos, 4))

		}
	}

}

// func channel(win *pixelgl.Window, message string) {
// 	// Draw a rectangle
// 	imd := imdraw.New(nil)
// 	imd.Color = colornames.Blueviolet
// 	imd.Push(pixel.V(0, 0), pixel.V(400, 200))
// 	imd.Rectangle(0)
// 	imd.Draw(win)

// 	// Create a basic font
// 	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
// 	basicTxt := text.New(pixel.V(0, 0), basicAtlas)
// 	basicTxt.Color = colornames.Black
// 	basicTxt.WriteString(message)

// 	// initial pos
// 	textPos := pixel.V(0, 100)

// 	// Move the message within the rectangle
// 	for !win.Closed() {
// 		win.Clear(colornames.White) // Clear the window with white color

// 		// Update the text position (move to the right)
// 		textPos.X += 2 // Adjust speed as needed

// 		// Draw the rectangle
// 		imd.Draw(win)

// 		// Draw the text at the updated position
// 		basicTxt.Draw(win, pixel.IM.Moved(textPos))
// 		win.Update()

// 		if textPos.X > 400 {
// 			win.SetClosed(true)
// 		}

//		}
//	}

func main() {
	pixelgl.Run(run)
}
