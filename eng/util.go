package eng

import (
	"github.com/go-gl/mathgl/mgl32"
)

func Vec2(x, y int) mgl32.Vec2 {
	return mgl32.Vec2{float32(x), float32(y)}
}

func V(vector Vector) mgl32.Vec2 {
	return mgl32.Vec2{vector.X, vector.Y}
}
