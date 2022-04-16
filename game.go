package walkie

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/jakecoffman/walkie/eng"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	state      int
	Keys       [1024]bool
	vsync      bool
	fullscreen bool
	window     *eng.OpenGlWindow
	gui        *Gui

	projection mgl32.Mat4

	*eng.ResourceManager

	ParticleGenerator *eng.ParticleGenerator
	SpriteRenderer    *eng.SpriteRenderer
	TextRenderer      *eng.TextRenderer

	chaseBananaMode bool
	randomBombMode  bool

	level string
}

const (
	worldWidth  = 1920
	worldHeight = 1080
)

// Game state
const (
	stateActive = iota
	statePause
)

func (g *Game) New(openGlWindow *eng.OpenGlWindow) {
	g.vsync = true
	g.window = openGlWindow
	g.gui = NewGui(g)
	g.Keys = [1024]bool{}

	g.ResourceManager = eng.NewResourceManager()

	g.LoadShader("assets/shaders/main.vs.glsl", "assets/shaders/main.fs.glsl", "sprite")
	g.LoadShader("assets/shaders/particle.vs.glsl", "assets/shaders/particle.fs.glsl", "particle")
	g.LoadShader("assets/shaders/text.vs.glsl", "assets/shaders/text.fs.glsl", "text")

	//center := Vector{worldWidth / 2, worldHeight / 2}

	g.projection = mgl32.Ortho(0, worldWidth, worldHeight, 0, -1, 1)
	g.Shader("sprite").Use().SetInt("sprite", 0).SetMat4("projection", g.projection)
	g.Shader("particle").Use().SetInt("sprite", 0).SetMat4("projection", g.projection)
	g.SpriteRenderer = eng.NewSpriteRenderer(g.Shader("sprite"))
	g.TextRenderer = eng.NewTextRenderer(g.Shader("text"), float32(openGlWindow.Width), float32(openGlWindow.Height), "assets/fonts/Roboto-Light.ttf", 24)
	g.TextRenderer.SetColor(1, 1, 1, 1)

	// Load all textures by name
	_ = filepath.Walk("assets/textures", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}
		log.Println("Loading", info.Name())
		g.LoadTexture(fmt.Sprintf("assets/textures/%v", info.Name()), strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())))
		return nil
	})

	g.reset()
}

func (g *Game) Update(dt float64) {
	if g.state == statePause {
		return
	}
}

func (g *Game) Render(alpha float64) {
	if g.window.UpdateViewport {
		g.window.UpdateViewport = false
		g.window.ViewportWidth, g.window.ViewPortHeight = g.window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(g.window.ViewportWidth), int32(g.window.ViewPortHeight))
		g.TextRenderer.Use().SetMat4("projection", mgl32.Ortho2D(0, float32(g.window.Width), float32(g.window.Height), 0))
		log.Printf("update viewport %#v\n", g.window)
	}

	g.SpriteRenderer.DrawSprite(
		g.Texture("background"),
		mgl32.Vec2{worldWidth / 2, worldHeight / 2},
		mgl32.Vec2{worldWidth, worldHeight}, 0, eng.White,
	)

	g.gui.Render()
}

func (g *Game) Close() {
	g.gui.Destroy()
	g.Clear()
}

func (g *Game) pause() {
	g.state = statePause
}

func (g *Game) unpause() {
	g.state = stateActive
}

func (g *Game) reset() {

}

func (g *Game) MouseToSpace(x, y, ww, wh float32) eng.Vector {
	model := mgl32.Translate3D(0, 0, 0)
	obj, err := mgl32.UnProject(mgl32.Vec3{x, wh - y, 0}, model, g.projection, 0, 0, int(ww), int(wh))
	if err != nil {
		panic(err)
	}

	return eng.Vec(obj.X(), obj.Y())
}
