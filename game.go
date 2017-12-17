package main

import (
	"time"

	"github.com/ByteArena/box2d"
)

type Game struct {
	*World
	*Camera
	*Terrain
	*Player
	*GameRenderer
	Entities             []*Entity
	BlockTypeDescriptors map[BlockType]*BlockTypeDescriptor
}

type World struct {
	*box2d.B2World
}

type Camera struct {
	// Deprecated
	Velocity [2]float64
	// Deprecated
	Position [2]float64

	Width  float64
	Height float64
	Zoom   float64

	*box2d.B2Body
}

func (game *Game) Init(width, height uint64, descriptors []*BlockTypeDescriptor) {
	game.World = new(World)
	b2World := box2d.MakeB2World(box2d.MakeB2Vec2(0, 0))
	game.World.B2World = &b2World

	game.Player = NewPlayer(game)

	game.Camera = &Camera{
		Velocity: [2]float64{0, 0},
		Position: [2]float64{0, 0},
		Width:    0,
		Height:   0,
		Zoom:     1.0,
	}
	game.Camera.B2Body = game.Player.B2Body

	game.Terrain = new(Terrain)
	game.Terrain.Chunks = make(map[Coord]*Chunk)

	game.GameRenderer = new(GameRenderer)
	game.GameRenderer.Game = game
	game.GameRenderer.Init(width, height, descriptors)

	game.Entities = make([]*Entity, 0, 1)
	game.BlockTypeDescriptors = make(map[BlockType]*BlockTypeDescriptor)

	for _, descriptor := range descriptors {
		game.LoadBlockType(descriptor)
	}

	entity := NewEntity()
	entity.Set(&Block{
		BlockType: "stone",
		Coord: Coord{
			X: 0,
			Y: 0,
		},
		FrontFace: 0,
	})
	entity.Set(&Block{
		BlockType: "stone",
		Coord: Coord{
			X: 1,
			Y: -1,
		},
		FrontFace: 0,
	})
	entity.Set(&Block{
		BlockType: "stone",
		Coord: Coord{
			X: 2,
			Y: 0,
		},
	})
	entity.bake(game)
	entity.BlockContainerObject.body.SetTransform(box2d.MakeB2Vec2(5.0, 5.0), 0)
	entity.BlockContainerObject.body.SetLinearVelocity(box2d.MakeB2Vec2(4.0, 3.0))
	entity.BlockContainerObject.body.SetAngularVelocity(1.0)
	game.Entities = append(game.Entities, entity)
}

func (game *Game) Update(dt time.Duration) {
	game.World.B2World.Step(dt.Seconds(), 8, 3)
}

func (game *Game) Render() {
	game.GameRenderer.RenderGame(0)
}

func (game *Game) LoadBlockType(descriptor *BlockTypeDescriptor) {
	game.BlockTypeDescriptors[descriptor.BlockType] = descriptor
}

func (game *Game) LoadChunk(chunk *Chunk) {
	coord := chunk.Coord
	game.Terrain.Chunks[coord] = chunk

	chunk.bake(game)
}

func (game *Game) UnloadChunk(coord Coord) {
}
