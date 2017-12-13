package main

import (
	"github.com/ByteArena/box2d"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
)

const (
	BLOCK_CONTAINER_STATIC BlockContainerType = iota
	BLOCK_CONTAINER_DYNAMIC
)

// BlockContainerType determines whether container's box2d body is static or dynamic.
type BlockContainerType int

// Block container is which can have blocks.
type BlockContainer interface {
	// Set block at coordinate which is in given block object.
	Set(block *Block)
	// Get block at given coordinate.
	At(x, y int64) *Block

	ForEach(func(*Block))
}

// Block container object is a set of variables about rendering, physics.
type BlockContainerObject struct {
	container     BlockContainer
	containerType BlockContainerType

	vbo           uint32
	vao           uint32
	verticesCount int32

	body *box2d.B2Body
}

// Bind object's vertex array object
func (object *BlockContainerObject) Bind() {
	gl.BindVertexArray(object.vao)
}

// Draw object's vertex array object
func (object *BlockContainerObject) Draw() {
	object.Bind()
	gl.DrawArrays(gl.TRIANGLES, 0, object.verticesCount)
}

func (object *BlockContainerObject) Rebake(game *Game) {
	// Replace vertex buffer object
	gl.BindBuffer(gl.VERTEX_ARRAY, object.vbo)

	mesh, verticesCount, totalBytes := containerMesh(game, object.container)
	object.verticesCount = verticesCount
	gl.BufferSubData(gl.VERTEX_ARRAY, 0, totalBytes, gl.Ptr(mesh))

	// Destory all fixtures and recreate fixtures
	for fixtureIter := object.body.GetFixtureList(); fixtureIter != nil; fixtureIter = fixtureIter.GetNext() {
		fixtureIter.Destroy()
	}
	containerFixtures(game, object.container, object.body)
}

// Create a vertex array object which contains position, tex coordinate and texture layer and fixture definitions for box2d engine.
// Function caller have to create body (static or dynamic body) and create fixtures returned by function.
//
// return (vao, vbo uint32, verticesCount int32, fixtureDefs []*box2d.B2FixtureDef)
func NewBlockContainerObject(game *Game, blockContainer BlockContainer, containerType BlockContainerType) *BlockContainerObject {
	object := &BlockContainerObject{containerType: containerType, container: blockContainer}

	mesh, verticesCount, totalBytes := containerMesh(game, blockContainer)

	vbo := uint32(0)
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, totalBytes, gl.Ptr(mesh), gl.DYNAMIC_DRAW)

	vao := uint32(0)
	gl.GenVertexArrays(1, &(vao))
	gl.BindVertexArray(vao)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 6*4, gl.PtrOffset(5*4))

	object.vbo = vbo
	object.vao = vao
	object.verticesCount = verticesCount

	bodyDef := box2d.NewB2BodyDef()
	switch containerType {
	case BLOCK_CONTAINER_STATIC:
		bodyDef.Type = box2d.B2BodyType.B2_staticBody
	case BLOCK_CONTAINER_DYNAMIC:
		bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	}
	body := game.World.B2World.CreateBody(bodyDef)

	containerFixtures(game, blockContainer, body)

	object.body = body

	return object
}

// Create mesh vertices for a block container.
// (mesh []float32, verticesCount int, bytes int)
func containerMesh(game *Game, container BlockContainer) ([]float32, int32, int) {
	mesh := make([]float32, 0)
	verticesCount := int32(0)
	totalBytes := 0

	container.ForEach(func(block *Block) {
		blockMesh, count, bytes := blockMesh(game, block)

		mesh = append(mesh, blockMesh...)
		verticesCount += count
		totalBytes += bytes
	})

	return mesh, verticesCount, totalBytes
}

func containerFixtures(game *Game, container BlockContainer, body *box2d.B2Body) {
	container.ForEach(func(block *Block) {
		if block.BlockType == "" {
			return
		}

		fixDef := blockFixtureDef(game, block)
		fixture := body.CreateFixtureFromDef(fixDef)
		block.B2Fixture = fixture
	})
}

// Create mesh vertices for a block.
// (mesh []float32, verticesCount int, bytes int)
func blockMesh(game *Game, block *Block) ([]float32, int32, int) {
	layerIndex, exist := game.GameRenderer.GetTextureIndex(block.BlockType)
	vertices := 6
	floatsPerVertices := 6
	bytesPerFloat := 4
	bytes := vertices * floatsPerVertices * bytesPerFloat
	if !exist {
		return []float32{
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0,
		}, int32(vertices), bytes
	}
	layer := float32(layerIndex)
	x := float32(block.X)
	y := float32(block.Y)
	return []float32{
		-0.5 + x, 0.5 + y, 0, 0, 0, layer,
		-0.5 + x, -0.5 + y, 0, 0, 1, layer,
		0.5 + x, -0.5 + y, 0, 1, 1, layer,

		-0.5 + x, 0.5 + y, 0, 0, 0, layer,
		0.5 + x, -0.5 + y, 0, 1, 1, layer,
		0.5 + x, 0.5 + y, 0, 1, 0, layer,
	}, int32(vertices), bytes
}

func blockFixtureDef(game *Game, block *Block) *box2d.B2FixtureDef {
	descriptor := game.BlockTypeDescriptors[block.BlockType]

	vertices := make([]box2d.B2Vec2, len(descriptor.CollisionVertices))
	for i, vertex := range descriptor.CollisionVertices {
		vertices[i].X = vertex.X + float64(block.X)
		vertices[i].Y = vertex.Y + float64(block.Y)
	}
	shape := box2d.MakeB2PolygonShape()
	shape.Set(vertices, len(vertices))

	fixDef := box2d.MakeB2FixtureDef()
	fixDef.Shape = &shape
	fixDef.Restitution = descriptor.Restitution
	fixDef.Density = descriptor.Density
	fixDef.Friction = descriptor.Friction

	return &fixDef
}
