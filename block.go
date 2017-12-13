package main

import (
	"os"

	"github.com/ByteArena/box2d"
)

// Identifier to distinguish block type. It is string type.
type BlockType string

// Container for block informations such as block type, coordinates, etc.
type Block struct {
	BlockType
	Coord
	// Represents where block is facing.
	// 0 = right, 1 = up, 2 = left, 3 = bottom (counter-clock wise)
	FrontFace uint8

	*box2d.B2Fixture
}

type BlockTypeDescriptor struct {
	BlockType
	// Block name which will be appeared to user.
	Name string
	// Each vertex must be in range between -0.5 and 0.5.
	CollisionVertices []Vec2
	Restitution       float64
	Density           float64
	Friction          float64
	TextureFile       *os.File
}
