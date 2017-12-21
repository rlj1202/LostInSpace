package lostinspace

type BlockStorage interface {
	Set(*Block)
	At(BlockCoord) *Block
	ForEach(func(*Block))
}
