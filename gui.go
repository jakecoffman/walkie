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
	glfw *gui.GLFW
	gl3  *gui.OpenGL3

	game *Game

	vsync bool
}

func NewGui(game *Game) *Gui {
	g := &Gui{
		Context: imgui.CreateContext(nil),
		game:    game,
		vsync:   true,
	}
	io := imgui.CurrentIO()

	platform, err := gui.NewGLFW(io, game.window.Window)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	g.glfw = platform

	renderer, err := gui.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	g.gl3 = renderer

	imgui.CurrentIO().SetClipboard(clipboard{platform: g.glfw})

	return g
}

func (gui *Gui) Destroy() {
	gui.gl3.Dispose()
	gui.Context.Destroy()
}

type clipboard struct {
	platform *gui.GLFW
}

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

func (gui *Gui) Render() {
	p := gui.glfw

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

		if imgui.Checkbox("Vsync", &gui.vsync) {
			if gui.vsync {
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
	gui.gl3.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
}
