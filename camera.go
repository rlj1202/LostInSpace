package lostinspace

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	zoom float64

	hwidth  float64
	hheight float64

	hZoomWidth  float64
	hZoomHeight float64

	target *Body

	projectionMat mgl32.Mat4
	cameraMat     mgl32.Mat4
}

func NewCamera(width, height float64) *Camera {
	camera := new(Camera)
	camera.zoom = 1.0
	camera.hwidth = width / 2.0
	camera.hheight = height / 2.0

	camera.createProjectionMat()

	return camera
}

func (camera *Camera) SetTarget(body *Body) {
	camera.target = body
}

func (camera *Camera) GetSize() (width float64, height float64) {
	width = camera.hwidth * 2
	height = camera.hheight * 2
	return
}

func (camera *Camera) GetProjectionMat() mgl32.Mat4 {
	return camera.projectionMat
}

func (camera *Camera) GetCameraMat() mgl32.Mat4 {
	if camera.target == nil {
		return mgl32.Ident4()
	}
	x, y := camera.target.GetPosition()
	return mgl32.Translate3D(float32(-x), float32(-y), 0)
}

func (camera *Camera) GetZoom() float64 {
	return camera.zoom
}

func (camera *Camera) SetZoom(zoom float64) {
	camera.zoom = math.Max(zoom, 0.1)
	camera.createProjectionMat()
}

func (camera *Camera) GetAABB() *AABB {
	if camera.target == nil {
		return nil
	}
	x, y := camera.target.GetPosition()

	return &AABB{
		Center:  Vec2{x, y},
		HWidth:  camera.hZoomWidth,
		HHeight: camera.hZoomHeight,
	}
}

func (camera *Camera) createProjectionMat() {
	camera.hZoomWidth = camera.hwidth / camera.zoom
	camera.hZoomHeight = camera.hheight / camera.zoom
	hw := float32(camera.hZoomWidth)
	hh := float32(camera.hZoomHeight)
	camera.projectionMat = mgl32.Ortho2D(-hw, hw, -hh, hh)
}
