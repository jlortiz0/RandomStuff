package main

import (
	"strings"
	"time"

	"github.com/adrg/sysfont"
	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var COLOR_BLACK = sdl.Color{A: 255}
var COLOR_WHITE = sdl.Color{R: 255, G: 255, B: 255, A: 255}

var window *sdl.Window
var display *sdl.Renderer
var font *ttf.Font

var fHeight int32

func main() {
	if !dialog.Message("This doesn't work.\nContinue?").YesNo() {
		return
	}
	err := sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()
	sdl.EventState(sdl.MOUSEMOTION, sdl.DISABLE)
	sdl.EventState(sdl.KEYUP, sdl.DISABLE)
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "2")

	initWindow()
	defer window.Destroy()
	defer display.Destroy()
	var fontName string
	for _, v := range sysfont.NewFinder(nil).List() {
		if v.Name == "Arial" {
			fontName = v.Filename
		} else if v.Name == "Ubuntu Mono" || strings.HasSuffix(v.Filename, "UbuntuMono-Regular.ttf") {
			fontName = v.Filename
			break
		}
	}
	err = ttf.Init()
	if err != nil {
		panic(err)
	}
	font, err = ttf.OpenFont(fontName, 24)
	if err != nil {
		panic(err)
	}
	fHeight = int32(font.Height()) + 10
	defer font.Close()

	prevDelay = time.Now()
	display.Clear()
	initExec()
	display.Present()
	time.Sleep(2 * time.Second)
	cleanupExec()
}

var prevDelay time.Time

func delay() {
	target := time.Since(prevDelay).Milliseconds()
	target = 33 - target
	if target < 0 {
		target = 0
	}
	sdl.Delay(uint32(target))
	prevDelay = time.Now()
}

func drawText(text string, dest *sdl.Surface, x, y int32, bg, fg sdl.Color) {
	txtSurf, err := font.RenderUTF8Shaded(text, fg, bg)
	if err != nil {
		panic(err)
	}
	txtSurf.Blit(nil, dest, &sdl.Rect{X: x, Y: y})
	txtSurf.Free()
}

func initWindow() {
	var err error
	window, display, err = sdl.CreateWindowAndRenderer(1024, 768, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	window.SetTitle("Mini RISC-V")
}
