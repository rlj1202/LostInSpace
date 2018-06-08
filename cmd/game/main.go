package main

import (
	"fmt"
	"image"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/rlj1202/LostInSpace"
)

func main() {
	runtime.LockOSThread()

	defer glfw.Terminate()

	icons := icons()

	window := lostinspace.NewWindow(800, 600, "LostInSpace", icons, true)
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

func icons() []image.Image {
	file16, err := os.Open("icon_16_16.png")
	if err != nil {
		panic(err)
	}
	defer file16.Close()

	file32, err := os.Open("icon_32_32.png")
	if err != nil {
		panic(err)
	}
	defer file32.Close()

	file64, err := os.Open("icon_64_64.png")
	if err != nil {
		panic(err)
	}
	defer file64.Close()

	icon16, _, err := image.Decode(file16)
	if err != nil {
		panic(err)
	}
	icon32, _, err := image.Decode(file32)
	if err != nil {
		panic(err)
	}
	icon64, _, err := image.Decode(file64)
	if err != nil {
		panic(err)
	}

	return []image.Image{icon16, icon32, icon64}
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
		Restitution: 0.01,
		CollisionVertices: []lostinspace.Vec2{
			{-0.5, 0.5},
			{-0.5, -0.5},
			{0.5, -0.5},
			{0.5, 0.5},
		},
		Fixed:       true,
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
		Restitution: 0.01,
		CollisionVertices: []lostinspace.Vec2{
			{-0.5, 0.5},
			{-0.5, -0.5},
			{0.5, -0.5},
			{0.5, 0.5},
		},
		Fixed:       true,
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
		Restitution: 0.01,
		CollisionVertices: []lostinspace.Vec2{
			{-0.5, 0.5},
			{-0.5, -0.5},
			{0.5, -0.5},
			{0.5, 0.5},
		},
		Fixed:       true,
		TextureFile: test2TypeTexFile,
	}
	doorTypeTexFile, err := os.Open("door_0.png")
	if err != nil {
		panic(err)
	}
	doorTypeDescriptor := lostinspace.BlockTypeDescriptor{
		BlockType:   "door0",
		Name:        "Test Door",
		Density:     1.0,
		Friction:    0.2,
		Restitution: 0.01,
		CollisionVertices: []lostinspace.Vec2{
			{-0.5, 0.25},
			{-0.5, -0.25},
			{0.5, -0.25},
			{0.5, 0.25},
		},
		Fixed:       false,
		TextureFile: doorTypeTexFile,
	}

	dic := lostinspace.NewBlockTypeDictionary(
		[]*lostinspace.BlockTypeDescriptor{
			&stoneTypeDescriptor,
			&test1TypeDescriptor,
			&test2TypeDescriptor,
			&doorTypeDescriptor,
		})

	return dic
}
