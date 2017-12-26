package lostinspace

type Block struct {
	BlockType
	FrontFace int

	coord BlockCoord
}

func NewBlock(blockCoord BlockCoord, blockType BlockType) *Block {
	block := &Block{
		BlockType: blockType,
		coord:     blockCoord,
	}

	return block
}
