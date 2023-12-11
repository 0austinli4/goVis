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
	name      string // name of philosopher
	count     int    // count number of times they eat
	left      int    // fork number on the left
	right     int    // fork number on the right
	spritePos pixel.Vec
	textPos   pixel.Vec
	eating    bool
}

type Fork struct {
	owner string // name of philosopher
	mat   pixel.Matrix
	pos   pixel.Vec
}

var philosophers []*Philosopher = make([]*Philosopher, 0)
var nameToIndex = make(map[string]int)
var textSeg []string = make([]string, 0)
var table []sync.Mutex = make([]sync.Mutex, 5)
var forks []*Fork

// load in png pictures
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

// initialize philsopohers animation (5 gopghers in circle)
func initializePhilosophers(win *pixelgl.Window) {
	pic, err := loadPicture("hiking.png")
	if err != nil {
		panic(err)
	}
	sprite := pixel.NewSprite(pic, pic.Bounds())
	spritePos := pixel.V(0, 0)
	//textPos := pixel.V(0, 0)
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	philosophers = []*Philosopher{
		{"Michelle", 0, 0, 1, spritePos, spritePos, false},
		{"Bill", 0, 1, 2, spritePos, spritePos, false},
		{"Sonia", 0, 2, 3, spritePos, spritePos, false},
		{"Brooke", 0, 3, 4, spritePos, spritePos, false},
		{"Eric", 0, 4, 0, spritePos, spritePos, false},
	}
	// assign name to index for animation
	i := 0
	for _, p := range philosophers {
		nameToIndex[p.name] = i
		i++
	}
	// declare initial circle
	centerX := (win.Bounds().Center()).X
	centerY := (win.Bounds().Center()).Y
	radius := 300.0
	numSprites := len(philosophers)
	angleIncrement := (2 * math.Pi) / float64(numSprites)
	initialAngle := math.Pi / 2 // Start from the top

	for i := 0; i < numSprites; i++ {
		angle := initialAngle + float64(i)*angleIncrement
		spritePos = pixel.V(
			centerX+radius*math.Cos(angle),
			centerY+radius*math.Sin(angle),
		)
		mat := pixel.IM
		mat = mat.ScaledXY(spritePos, pixel.V(0.15, 0.15))
		philosophers[i].spritePos = spritePos
		sprite.Draw(win, mat)

		radiusText := 220.0
		textPos := pixel.V(
			centerX-80+radiusText*math.Cos(angle),
			centerY+radiusText*math.Sin(angle),
		)
		philosophers[i].textPos = textPos
		basicTxt := text.New(textPos, basicAtlas)
		basicTxt.Color = colornames.Black
		fmt.Fprintln(basicTxt, philosophers[i].name)
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 1.5))
	}
	// text position for CLICK TO START
	textPos := pixel.V(
		win.Bounds().Max.X/2-200,
		win.Bounds().Max.Y-float64(100),
	)
	basicTxt := text.New(textPos, basicAtlas)
	basicTxt.Color = colornames.Black
	fmt.Fprintln(basicTxt, "CLICK TO START")
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
}

// initialize forks
func initializeForks(win *pixelgl.Window) {
	pic, err := loadPicture("fork.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	mat := pixel.IM

	forks = []*Fork{
		{"", mat, pixel.V(0, 0)},
		{"", mat, pixel.V(0, 0)},
		{"", mat, pixel.V(0, 0)},
		{"", mat, pixel.V(0, 0)},
		{"", mat, pixel.V(0, 0)},
	}

	centerX := (win.Bounds().Center()).X
	centerY := (win.Bounds().Center()).Y
	radius := 300.0
	numSprites := len(forks)

	angleIncrement := (2 * math.Pi) / float64(numSprites)
	initialAngle := math.Pi/2 + (-0.5)*math.Pi/4 // Start from the top

	for i := 0; i < numSprites; i++ {
		angle := initialAngle + float64(i)*angleIncrement
		spritePos := pixel.V(
			centerX+radius*math.Cos(angle),
			centerY+radius*math.Sin(angle),
		)
		mat := pixel.IM
		mat = mat.ScaledXY(spritePos, pixel.V(0.15, 0.15))
		spritePos = pixel.V(
			centerX+radius*math.Cos(angle)-50,
			centerY+radius*math.Sin(angle)-50,
		)
		forks[i].pos = spritePos
		forks[i].mat = mat
		sprite.Draw(win, mat)
	}
}

// draw one new frame of animation
// updates: philsopohers, fork ownership, display of actions in text
func drawNewFrame(win *pixelgl.Window) {
	forkPic, _ := loadPicture("fork.png")
	standing, _ := loadPicture("hiking.png")
	eating, _ := loadPicture("gamer.png")

	spriteFork := pixel.NewSprite(forkPic, forkPic.Bounds())
	spriteStand := pixel.NewSprite(standing, standing.Bounds())
	spriteEat := pixel.NewSprite(eating, eating.Bounds())

	// clear frame
	win.Clear(colornames.Aliceblue)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// draw visual for forks
	for i := 0; i < len(forks); i++ {
		mat := forks[i].mat
		spriteFork.Draw(win, mat)

		basicTxt := text.New(forks[i].pos, basicAtlas)
		basicTxt.Color = colornames.Blueviolet

		fmt.Fprintln(basicTxt, fmt.Sprint(i))
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 1.2))
	}

	// draw philosophers and names
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
		textPos := philosophers[i].textPos
		basicTxt := text.New(textPos, basicAtlas)
		basicTxt.Color = colornames.Black
		fmt.Fprintln(basicTxt, philosophers[i].name+" (Dined: "+fmt.Sprint(philosophers[i].count)+")")
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 1.2))
	}

	// draw latest 8 actions
	basicTxt := text.New(pixel.V(20, 80), basicAtlas)
	basicTxt.Color = colornames.Black
	if len(textSeg) > 8 {
		// Remove the first element if the queue length exceeds 10
		textSeg = textSeg[len(textSeg)-8:]
	}
	fmt.Fprintln(basicTxt, "Order of Actions")
	for _, segment := range textSeg {
		fmt.Fprintln(basicTxt, segment)
	}
	basicTxt.Draw(win, pixel.IM)

	// draw forks and ownership
	for i, fork := range forks {
		textstr := "Fork " + fmt.Sprint(i) + ": " + fork.owner
		textPos := pixel.V(win.Bounds().Max.X-200, win.Bounds().Max.Y-100-float64(20*i))
		basicTxt = text.New(textPos, basicAtlas)
		basicTxt.Color = colornames.Black
		fmt.Fprintln(basicTxt, textstr)
		basicTxt.Draw(win, pixel.IM)
	}

	basicTxt.Draw(win, pixel.IM)
	win.Update()
	for !win.JustPressed(pixelgl.MouseButtonLeft) {
		win.Update()
	}
}

// declare philosophers to be eating
func updateEat(p *Philosopher) {
	philosophers[nameToIndex[p.name]].eating = !(philosophers[nameToIndex[p.name]].eating)
}

// think function: delays by a random time
func (p *Philosopher) Think(win *pixelgl.Window) {
	fmt.Println(p.name, "is thinking.")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println(p.name, "is done thinking.")
}

// eat function  -> take random time duration to finish eating
func (p *Philosopher) Eat(win *pixelgl.Window) {
	p.count++
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println(p.name, "is done eating.")
}

// Dine function : declares algorithm
func (p *Philosopher) Dine(win *pixelgl.Window) {
	for {
		p.Think(win)

		table[p.left].Lock()
		fmt.Printf("%s picks up fork %d.\n", p.name, p.left)
		forks[p.left].owner = p.name
		textSeg = append(textSeg, p.name+" picks up fork "+fmt.Sprint(p.left))

		runtime.Gosched() // hack to yield to the next goroutine

		table[p.right].Lock()
		forks[p.right].owner = p.name
		fmt.Printf("%s picks up %d.\n", p.name, p.right)
		textSeg = append(textSeg, p.name+" picks up fork "+fmt.Sprint(p.right))

		textSeg = append(textSeg, p.name+" is eating round:"+fmt.Sprint(p.count))
		p.Eat(win)
		updateEat(p)
		drawNewFrame(win)

		table[p.right].Unlock()
		forks[p.right].owner = ""
		fmt.Printf("%s puts down fork %d.\n", p.name, p.right)
		textSeg = append(textSeg, p.name+" puts down fork "+fmt.Sprint(p.right))

		table[p.left].Unlock()
		forks[p.left].owner = ""
		fmt.Printf("%s puts down fork %d.\n", p.name, p.left)
		textSeg = append(textSeg, p.name+" puts down fork "+fmt.Sprint(p.left))

		updateEat(p)
		drawNewFrame(win)

	}
}
