package main

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gleroi/codingame/puzzle/mars_lander/ep2/client"
	"golang.org/x/image/colornames"
)

var examples = []struct {
	surface string
	lander  string
}{
	{
		surface: `7
	0 100
	1000 500
	1500 1500
	3000 1000
	4000 150
	5500 150
	6999 800`,
		lander: `2500 2700 0 0 550 0 0`,
	},

	{surface: `10
0 100
1000 500
1500 100
3000 100
3500 500
3700 200
5000 1500
5800 300
6000 1000
6999 2000`,
		lander: `6500 2800 -100 0 600 90 0`,
	},
}

func run(srvIn io.Reader, srvOut io.Writer) {
	cfg := pixelgl.WindowConfig{
		Title:     "Pixel Rocks!",
		Bounds:    pixel.R(0, 0, 1024, 768),
		Resizable: true,
		VSync:     true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	id := 0

	debug("write surface\n")
	fmt.Fprintln(srvOut, examples[id].surface)
	debug("write lander\n")
	fmt.Fprintln(srvOut, examples[id].lander)

	surface := client.ReadSurface(strings.NewReader(examples[id].surface), false)
	lander := client.ReadLander(strings.NewReader(examples[id].lander), false)

	landers := make([]client.Lander, 0, 1024)
	done := false
	for !win.Closed() && !done {
		start := time.Now()

		win.Clear(colornames.Black)
		camera := pixel.IM.ScaledXY(pixel.V(0, 0),
			pixel.V(win.Bounds().Size().X/7000, win.Bounds().Size().Y/3000))
		win.SetMatrix(camera)

		imd := imdraw.New(nil)
		drawSurface(imd, surface)

		drawLander(imd, lander, colornames.White, 50)

		for _, l := range landers {
			drawLander(imd, l, colornames.Green, 5)
		}

		imd.Draw(win)
		win.Update()

		var angle, power float64

		debug("read input lander\n")
		fmt.Fscan(srvIn, &angle, &power)

		if power > lander.Power {
			power = math.Min(lander.Power+1, 4)
		} else if power < lander.Power {
			power = math.Max(0, lander.Power-1)
		}

		if lander.Fuel < 0 {
			power = 0
		}

		if angle-lander.Rotation > 15 {
			angle = lander.Rotation + 15
		} else if angle-lander.Rotation < -15 {
			angle = lander.Rotation - 15
		}

		landers = append(landers, lander.Round())
		lander = lander.Next(1, float64(power), float64(angle))

		debug("write output lander\n")
		landerR := lander.Round()
		fmt.Fprintln(srvOut, landerR.Pos.X, landerR.Pos.Y, landerR.Speed.X, landerR.Speed.Y,
			landerR.Fuel, landerR.Rotation, landerR.Power)

		delay := time.Now().Sub(start)
		time.Sleep((1*time.Second - delay) / (5 * time.Millisecond))

		done = (lander.Pos.X < 0 || lander.Pos.X > 7000 || lander.Pos.Y < 0 || lander.Pos.Y > 3000)
		debug("done: %t\n", done)
	}
	debug("Closed!\n")
}

func drawSurface(imd *imdraw.IMDraw, surface *client.Surface) {
	imd.Color = colornames.Brown
	for _, pt := range surface.Points {
		imd.Push(pixel.V(pt.X, pt.Y))
	}
	imd.Line(2)
}

func drawLander(imd *imdraw.IMDraw, lander client.Lander, color color.RGBA, size float64) {
	imd.Color = color

	center := pixel.V(lander.Pos.X, lander.Pos.Y)
	imd.SetMatrix(pixel.IM.Rotated(center, client.Radian(lander.Rotation)-math.Pi/2))
	imd.Push(center.Add(pixel.V(0, -size/3)), center.Add(pixel.V(0, size/3)))
	imd.Ellipse(pixel.V(size/2, size/4), 10)

	imd.Color = colornames.Red
	thrust := center.Add(pixel.V(0, -size/2))
	imd.Push(thrust, thrust.Add(pixel.V(0, -50)))
	imd.Line(10)

	imd.SetMatrix(pixel.IM)
}

func main() {
	cltIn, srvOut := io.Pipe()
	srvIn, cltOut := io.Pipe()

	go client.Client(cltIn, cltOut)

	pixelgl.Run(func() {
		run(srvIn, srvOut)
	})
}

func debug(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "server: "+format, v...)
}
