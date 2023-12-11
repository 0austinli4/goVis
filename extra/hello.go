package main

import (
	"image"
	"math"
	"os"

	_ "image/png"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Philosopher struct {
	name   string       // name of philosopher
	count  int          // count number of times they eat
	left   int          // fork number on the left
	right  int          // fork number on the right
	gopher pixel.Sprite // animation
	mat    pixel.Matrix
}

type Fork struct {
	index   int // name of philosopher
	owner   string
	forkPic pixel.Sprite // animation
	mat     pixel.Matrix
}

var philosophers []*Philosopher = make([]*Philosopher, 0)
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

func initializeDrawing(win *pixelgl.Window) {
	for i := 0; i < len(philosophers); i++ {
		sprite := philosophers[i].gopher
		mat := philosophers[i].mat
		sprite.Draw(win, mat)

		spriteFork := forks[i].forkPic
		matFork := forks[i].mat

		spriteFork.Draw(win, matFork)
	}
}

func initializePhilosophers(win *pixelgl.Window) {
	pic, err := loadPicture("hiking.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	mat := pixel.IM

	philosophers := []*Philosopher{
		{"Michelle", 0, 0, 1, *sprite, mat},
		{"Bill", 0, 1, 2, *sprite, mat},
		{"Sonia", 0, 2, 3, *sprite, mat},
		{"Brooke", 0, 3, 4, *sprite, mat},
		{"Eric", 0, 4, 0, *sprite, mat},
	}

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
		philosophers[i].mat = mat
		sprite.Draw(win, mat)
	}
}

func initializeForks(win *pixelgl.Window) {
	pic, err := loadPicture("fork.png")
	if err != nil {
		panic(err)
	}

	forks = make([]*Fork, 5)

	sprite := pixel.NewSprite(pic, pic.Bounds())
	mat := pixel.IM
	forks := []*Fork{
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
		{0, "", *sprite, mat},
	}

	centerX := (win.Bounds().Center()).X
	centerY := (win.Bounds().Center()).Y
	radius := 300.0
	numSprites := len(philosophers)

	angleIncrement := (2 * math.Pi) / float64(numSprites)
	initialAngle := math.Pi / 2 // Start from the top

	for i := 0; i < 5; i++ {
		angle := initialAngle + float64(i)*angleIncrement
		spritePos := pixel.V(
			centerX+radius*math.Cos(angle),
			centerY+radius*math.Sin(angle),
		)
		mat = mat.ScaledXY(spritePos, pixel.V(0.15, 0.15))
		forks[i].mat = mat
		sprite.Draw(win, mat)
	}

	for !win.Closed() {
		win.Update()
	}

}

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
	initializeForks(win)
	initializePhilosophers(win)
	initializeDrawing(win)

	pic, err := loadPicture("fork.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	win.Clear(colornames.Greenyellow)
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	for !win.Closed() {
		win.Update()
	}

	// fmt.Println("initialized forks"
}

func main() {
	pixelgl.Run(run)
}

// func draw(win *pixelgl.Window, message string) {
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
