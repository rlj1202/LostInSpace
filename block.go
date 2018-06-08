package lostinspace

type Block struct {
	BlockType
	FrontFace int

	coord BlockCoord
}

func NewBlock(blockCoord BlockCoord, blockType BlockType, face int) *Block {
	block := &Block{
		BlockType: blockType,
		FrontFace: face,

		coord: blockCoord,
	}

	return block
}
