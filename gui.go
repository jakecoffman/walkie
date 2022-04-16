package walkie

import (
	"fmt"
	"os"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/inkyblackness/imgui-go"
	"github.com/jakecoffman/walkie/gui"
)

type Gui struct {
	*imgui.Context
	platform *gui.GLFW
	renderer *gui.OpenGL3

	game *Game
}

func NewGui(game *Game) *Gui {
	g := &Gui{
		Context: imgui.CreateContext(nil),
		game:    game,
	}
	io := imgui.CurrentIO()

	platform, err := gui.NewGLFW(io, game.window.Window)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	g.platform = platform

	renderer, err := gui.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	g.renderer = renderer

	imgui.CurrentIO().SetClipboard(clipboard{platform: g.platform})

	return g
}

func (gui *Gui) Destroy() {
	gui.renderer.Dispose()
	gui.Context.Destroy()
}

// Platform covers mouse/keyboard/gamepad inputs, cursor shape, timing, windowing.
type Platform interface {
	// ShouldStop is regularly called as the abort condition for the program loop.
	ShouldStop() bool
	// ProcessEvents is called once per render loop to dispatch any pending events.
	ProcessEvents()
	// DisplaySize returns the dimension of the display.
	DisplaySize() [2]float32
	// FramebufferSize returns the dimension of the framebuffer.
	FramebufferSize() [2]float32
	// NewFrame marks the begin of a render pass. It must update the imgui IO state according to user input (mouse, keyboard, ...)
	NewFrame()
	// PostRender marks the completion of one render pass. Typically this causes the display buffer to be swapped.
	PostRender()
	// ClipboardText returns the current text of the clipboard, if available.
	ClipboardText() (string, error)
	// SetClipboardText sets the text as the current text of the clipboard.
	SetClipboardText(text string)
}

type clipboard struct {
	platform Platform
}

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

func (gui *Gui) Render() {
	p := gui.platform

	p.NewFrame()
	imgui.NewFrame()

	{
		imgui.Begin("Menu")

		if imgui.Button("Resume game") {
			gui.game.unpause()
		}

		if imgui.Button("Reset objects") {
			gui.game.reset()
		}

		imgui.Checkbox("Chase Banana", &gui.game.chaseBananaMode)
		imgui.Checkbox("Random Bombs", &gui.game.randomBombMode)
		if imgui.Checkbox("Vsync", &gui.game.vsync) {
			if gui.game.vsync {
				glfw.SwapInterval(1)
			} else {
				glfw.SwapInterval(0)
			}
		}

		if imgui.ButtonV("Quit", imgui.Vec2{200, 20}) {
			gui.game.window.SetShouldClose(true)
		}

		// TODO add text of FPS based on IO.Framerate()

		imgui.End()
	}

	imgui.Render()
	gui.renderer.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
}
