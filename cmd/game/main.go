package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/rlj1202/LostInSpace"
)

func main() {
	runtime.LockOSThread()

	defer glfw.Terminate()

	window := lostinspace.NewWindow(800, 600, "LostInSpace", false)
	game := lostinspace.NewGame(window, blockTypeDic())

	curTime := time.Now()
	for !window.ShouldClose() {
		glfw.PollEvents()

		newTime := time.Now()
		frameTime := newTime.Sub(curTime)

		if frameTime > 0 {
			window.SetTitle(fmt.Sprintf("LostInSpace - %v", time.Second/frameTime))
		}

		curTime = newTime

		for frameTime > 0 {
			deltaTime := time.Duration(1.0 / 60.0 * float32(time.Second))
			if frameTime < deltaTime {
				deltaTime = frameTime
			}

			game.Update(deltaTime)

			frameTime -= deltaTime
		}

		game.Render()

		window.Update()
	}
	game.Destroy()
}

func blockTypeDic() *lostinspace.BlockTypeDictionary {
	stoneTypeTexFile, err := os.Open("stonetile_1.png")
	if err != nil {
		panic(err)
	}
	stoneTypeDescriptor := lostinspace.BlockTypeDescriptor{
		BlockType:   "stone",
		Name:        "Regular old fancy stone",
		Density:     0.5,
		Friction:    0.2,
		Restitution: 0.1,
		CollisionVertices: []lostinspace.Vec2{
			{-0.5, 0.5},
			{-0.5, -0.5},
			{0.5, -0.5},
			{0.5, 0.5},
		},
		TextureFile: stoneTypeTexFile,
	}
	test1TypeTexFile, err := os.Open("testtile_1.png")
	if err != nil {
		panic(err)
	}
	test1TypeDescriptor := lostinspace.BlockTypeDescriptor{
		BlockType:   "test1",
		Name:        "Test tile 1",
		Density:     1.0,
		Friction:    0.2,
		Restitution: 0.5,
		CollisionVertices: []lostinspace.Vec2{
			{-0.5, 0.5},
			{-0.5, -0.5},
			{0.5, -0.5},
			{0.5, 0.5},
		},
		TextureFile: test1TypeTexFile,
	}
	test2TypeTexFile, err := os.Open("testtile_2.png")
	if err != nil {
		panic(err)
	}
	test2TypeDescriptor := lostinspace.BlockTypeDescriptor{
		BlockType:   "test2",
		Name:        "Test tile 2",
		Density:     1.0,
		Friction:    0.2,
		Restitution: 0.5,
		CollisionVertices: []lostinspace.Vec2{
			{-0.5, 0.5},
			{-0.5, -0.5},
			{0.5, -0.5},
			{0.5, 0.5},
		},
		TextureFile: test2TypeTexFile,
	}

	dic := lostinspace.NewBlockTypeDictionary(
		[]*lostinspace.BlockTypeDescriptor{
			&stoneTypeDescriptor,
			&test1TypeDescriptor,
			&test2TypeDescriptor,
		})

	return dic
}
