package lostinspace

type Block struct {
	BlockType
	BlockCoord
}

func NewBlock(blockCoord BlockCoord, blockType BlockType) *Block {
	block := &Block{
		BlockType:  blockType,
		BlockCoord: blockCoord,
	}

	return block
}
