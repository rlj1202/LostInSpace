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
	Entities             []*box2d.B2Body
	BlockTypeDescriptors map[BlockType]*BlockTypeDescriptor // TODO
}

type World struct {
	box2d.B2World
	TerrainBodies map[Coord]*TerrainBody
}

type TerrainBody struct {
	*box2d.B2Body
	BlockFixtures [CHUNK_WIDTH * CHUNK_HEIGHT]*box2d.B2Fixture
}

type Player struct {
}

type Camera struct {
	Velocity [2]float64
	Position [2]float64

	Width  float64
	Height float64
	Zoom   float64
}

func (game *Game) Init(width, height uint64, blockTypeDescriptors []*BlockTypeDescriptor) {
	game.World = new(World)
	game.World.Init()
	game.Camera = &Camera{[2]float64{0, 0}, [2]float64{0, 0}, 0, 0, 1.0}
	game.Terrain = new(Terrain)
	game.Terrain.Chunks = make(map[Coord]*Chunk)
	game.Player = new(Player)
	game.GameRenderer = new(GameRenderer)
	game.GameRenderer.Game = game
	game.GameRenderer.Init(width, height, blockTypeDescriptors)
	game.Entities = make([]*box2d.B2Body, 0, 1)
	game.BlockTypeDescriptors = make(map[BlockType]*BlockTypeDescriptor)

	for _, blockTypeDescriptor := range blockTypeDescriptors {
		game.LoadBlockType(blockTypeDescriptor)
	}

	// Test physical object code below
	body := box2d.MakeB2BodyDef()
	body.Position.Set(5.0, 5.0)
	body.Type = box2d.B2BodyType.B2_dynamicBody

	square := game.B2World.CreateBody(&body)

	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(0.5, 0.5)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 20.0
	square.CreateFixtureFromDef(&fd)
	square.SetLinearVelocity(box2d.MakeB2Vec2(4.0, 1.0))

	game.Entities = append(game.Entities, square)
}

func (game *Game) Update(dt time.Duration) {
	game.B2World.Step(dt.Seconds(), 8, 3)
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

	terrainBodyDef := box2d.MakeB2BodyDef()
	terrainBodyDef.Position.Set(float64(coord.X*CHUNK_WIDTH), float64(coord.Y*CHUNK_HEIGHT))
	terrainBodyDef.Type = box2d.B2BodyType.B2_staticBody

	terrainBody := game.World.CreateBody(&terrainBodyDef)
	game.World.TerrainBodies[coord] = &TerrainBody{B2Body: terrainBody}

	blockShape := box2d.MakeB2PolygonShape()

	blockFixtureDef := box2d.MakeB2FixtureDef()
	blockFixtureDef.Shape = &blockShape
	blockFixtureDef.Restitution = 1.0
	for y := uint8(0); y < CHUNK_WIDTH; y++ {
		for x := uint8(0); x < CHUNK_HEIGHT; x++ {
			block := chunk.GetBlock(x, y)

			if block.BlockType == "" {
				continue
			}

			blockTypeDescriptor := game.BlockTypeDescriptors[block.BlockType]
			vertices := make([]box2d.B2Vec2, len(blockTypeDescriptor.CollisionVertices))
			for i, vertex := range blockTypeDescriptor.CollisionVertices {
				vertices[i].X = vertex.X + float64(x)
				vertices[i].Y = vertex.Y + float64(y)
			}

			//blockShape.SetAsBoxFromCenterAndAngle(0.5, 0.5, box2d.MakeB2Vec2(float64(x), float64(y)), 0)
			blockShape.Set(vertices, len(vertices))
			blockFixture := terrainBody.CreateFixtureFromDef(&blockFixtureDef)
			game.World.TerrainBodies[coord].BlockFixtures[x+y*CHUNK_WIDTH] = blockFixture
		}
	}
}

func (game *Game) UnloadChunk(coord Coord) {
}

func (world *World) Init() {
	world.B2World = box2d.MakeB2World(box2d.MakeB2Vec2(0, 0))
	world.TerrainBodies = make(map[Coord]*TerrainBody)
}
