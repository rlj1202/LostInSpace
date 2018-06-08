package lostinspace

import (
	"fmt"
	"os"
)

// An block type to represents empty space. or void, whatever.
const BLOCK_TYPE_VOID BlockType = ""

// BlockType represents type of block.
// Identical for each block types.
type BlockType string

// Struct to store all informations about all block types
// such as physics properties, texture file, texture array indices etc.
type BlockTypeDictionary struct {
	data     map[BlockType]*BlockTypeDescriptor
	arrayTex *Texture2DArray
}

// Struct to store informations about block type.
type BlockTypeDescriptor struct {
	// BlockType is id to distinguish block type.
	// It has to be unique.
	BlockType

	// Name will be shown in game.
	Name string

	Density     float64
	Friction    float64
	Restitution float64
	// Vertices which represent collision polygon.
	//
	// (x, y, x, y, x, y, ...)
	CollisionVertices []Vec2
	// Fixed property represents whether a block can move or can't.
	// Non-fixed block will be create as seperated body from blockcontainer
	// and will have a joint (prismatic joint for example) to stick together.
	Fixed bool

	TextureFile *os.File

	layerIndex int
}

func NewBlockTypeDictionary(descriptors []*BlockTypeDescriptor) *BlockTypeDictionary {
	dic := new(BlockTypeDictionary)
	dic.data = make(map[BlockType]*BlockTypeDescriptor)

	imgFiles := make([]*os.File, len(descriptors))
	for i, descriptor := range descriptors {
		dic.data[descriptor.BlockType] = descriptor
		imgFiles[i] = descriptor.TextureFile
		descriptor.layerIndex = i
	}
	dic.arrayTex = NewTexture2DArray(16, 16, imgFiles)

	return dic
}

// Get descriptor of given block type.
func (dic *BlockTypeDictionary) Get(blockType BlockType) *BlockTypeDescriptor {
	return dic.data[blockType]
}

func (desc *BlockTypeDescriptor) String() string {
	return fmt.Sprintf(
		`BlockTypeDescriptor{
			BlockType: "%s",
			Name: "%s",
			Density: %f,
			Friction: %f,
			Restitution: %f,
			Fixed: %t,
		}`,
		desc.BlockType, desc.Name, desc.Density, desc.Friction, desc.Restitution, desc.Fixed,
	)
}
