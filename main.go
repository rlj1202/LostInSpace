package main

import (
	"fmt"
	_ "image/png"
	"os"
	"runtime"

	"time"

	"github.com/ByteArena/box2d"
	"github.com/rlj1202/LostInSpace/event"

	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 800
	height = 600
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	window.SetKeyCallback(keyInput)
	window.SetMouseButtonCallback(mouseInput)
	window.SetScrollCallback(scrollInput)
	window.SetCursorEnterCallback(cursorEnterInput)
	window.SetCursorPosCallback(cursorPosInput)

	stoneBlockTextureFile, err := os.Open("stoneTile_1.png")
	if err != nil {
		panic(err)
	}
	stoneBlockType := BlockTypeDescriptor{
		BlockType: "stone",
		Name:      "Regular Old Fancy Stone",
		CollisionVertices: []Vec2{
			{-0.5, -0.5}, {0.5, -0.5}, {0.5, 0.5}, {-0.5, 0.5}},
		Restitution: 1.0,
		Density:     0.5,
		Friction:    0.2,
		TextureFile: stoneBlockTextureFile}
	test1BlockTextureFile, err := os.Open("TestTile1.png")

	if err != nil {
		panic(err)
	}
	test1BlockType := BlockTypeDescriptor{
		BlockType: "test1",
		Name:      "TestTile1",
		CollisionVertices: []Vec2{
			{-0.5, -0.5}, {0.5, -0.5}, {0.5, 0.5}, {-0.5, 0.5}},
		Restitution: 1.0,
		Density:     0.0,
		Friction:    0.2,
		TextureFile: test1BlockTextureFile}

	test2BlockTextureFile, err := os.Open("TestTile2.png")
	if err != nil {
		panic(err)
	}
	test2BlockType := BlockTypeDescriptor{
		BlockType: "test2",
		Name:      "TestTile2",
		CollisionVertices: []Vec2{
			{-0.5, -0.5}, {0.5, -0.5}, {0.5, 0.5}, {-0.5, 0.5}},
		Restitution: 1.0,
		Density:     0.0,
		Friction:    0.2,
		TextureFile: test2BlockTextureFile}

	blocksPerPixel := 30.0 / 800.0
	game := new(Game)
	game.Init(width, height, []*BlockTypeDescriptor{&stoneBlockType, &test1BlockType, &test2BlockType})
	game.Camera.Width = blocksPerPixel * width
	game.Camera.Height = game.Camera.Width * float64(height) / float64(width)
	game.Terrain.Seed = 2

	for x := int64(0); x < 8; x++ {
		for y := int64(0); y < 8; y++ {
			chunk := game.Terrain.GenerateChunk(Coord{x, y})
			game.LoadChunk(chunk)
		}
	}

	t := time.Duration(0)
	dt := time.Duration(1.0 / 60.0 * float32(time.Second))

	currentTime := time.Now()
	accumulator := time.Duration(0)
	for !window.ShouldClose() {
		glfw.PollEvents()

		newTime := time.Now()
		frameTime := newTime.Sub(currentTime)

		if frameTime > 0 {
			window.SetTitle(fmt.Sprintf("LostInSpace - %v", time.Second/frameTime))
		}

		currentTime = newTime

		accumulator += frameTime
		for frameTime > 0 {
			deltaTime := dt
			if frameTime < deltaTime {
				deltaTime = frameTime
			}

			update(game, window, deltaTime)
			game.Update(deltaTime)

			frameTime -= deltaTime

			accumulator -= dt
			t += dt
		}

		game.Render()

		window.SwapBuffers()
	}
}

func update(game *Game, window *glfw.Window, deltaTime time.Duration) {
	for e := event.Pop(); e != nil; e = event.Pop() {
		switch e.(type) {
		case event.MouseEvent:
			mouseEvent := e.(event.MouseEvent)
			action := mouseEvent.Action
			button := mouseEvent.Button
			xpos := mouseEvent.XPos
			ypos := mouseEvent.YPos
			if action == glfw.Press {
				if button == glfw.MouseButtonLeft {
					x, y := game.GameRenderer.ToWorldCoord(float32(xpos), float32(ypos))
					fmt.Printf("LMB pressed: world coordinate (%v, %v)\n", x, y)
				}
				if button == glfw.MouseButtonRight {
					fmt.Printf("RMB pressed: %v, %v\n", xpos, ypos)
				}
				if button == glfw.MouseButtonMiddle {
					fmt.Printf("MMB pressed: %v, %v\n", xpos, ypos)
				}
			}
		case event.KeyboardEvent:
			keyboardEvent := e.(event.KeyboardEvent)
			action := keyboardEvent.Action
			if action == glfw.Press {
				fmt.Println("Key pressed")
			} else if action == glfw.Release {
				fmt.Println("Key released")
			} else if action == glfw.Repeat {
				fmt.Println("Key repeated")
			}
		case event.ScrollEvent:
			scrollEvent := e.(event.ScrollEvent)
			fmt.Printf("Scroll x:%v, y:%v\n", scrollEvent.XOff, scrollEvent.YOff)
		case event.CursorEnterEvent:
			enterEvent := e.(event.CursorEnterEvent)
			fmt.Printf("Cursor entered: %v\n", enterEvent.Entered)
		case event.CursorPosEvent:
			posEvent := e.(event.CursorPosEvent)
			fmt.Printf("Cursor pos: %v, %v\n", posEvent.XPos, posEvent.YPos)
		}
	}
	acc := 1000.0
	force := 40.0
	if window.GetKey(glfw.KeyW) != glfw.Release {
		game.Camera.Velocity[1] += acc * deltaTime.Seconds()
		game.Camera.B2Body.ApplyForceToCenter(box2d.MakeB2Vec2(0, force), true)
	}
	if window.GetKey(glfw.KeyS) != glfw.Release {
		game.Camera.Velocity[1] -= acc * deltaTime.Seconds()
		game.Camera.B2Body.ApplyForceToCenter(box2d.MakeB2Vec2(0, -force), true)
	}
	if window.GetKey(glfw.KeyA) != glfw.Release {
		game.Camera.Velocity[0] -= acc * deltaTime.Seconds()
		game.Camera.B2Body.ApplyForceToCenter(box2d.MakeB2Vec2(-force, 0), true)
	}
	if window.GetKey(glfw.KeyD) != glfw.Release {
		game.Camera.Velocity[0] += acc * deltaTime.Seconds()
		game.Camera.B2Body.ApplyForceToCenter(box2d.MakeB2Vec2(force, 0), true)
	}
	game.Camera.Velocity[0] *= 0.9
	game.Camera.Velocity[1] *= 0.9
	game.Camera.Position[0] += game.Camera.Velocity[0] * deltaTime.Seconds()
	game.Camera.Position[1] += game.Camera.Velocity[1] * deltaTime.Seconds()
}

func keyInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	event.Push(event.KeyboardEvent{Key: key, Scancode: scancode, Action: action, Mods: mods})
}

func mouseInput(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	x, y := w.GetCursorPos()
	event.Push(event.MouseEvent{XPos: x, YPos: y, Button: button, Action: action, Mod: mod})
}

func scrollInput(w *glfw.Window, xoff float64, yoff float64) {
	event.Push(event.ScrollEvent{XOff: xoff, YOff: yoff})
}

func cursorEnterInput(w *glfw.Window, entered bool) {
	event.Push(event.CursorEnterEvent{Entered: entered})
}

func cursorPosInput(w *glfw.Window, xpos float64, ypos float64) {
	event.Push(event.CursorPosEvent{XPos: xpos, YPos: ypos})
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "LostInSpace", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	return window
}
