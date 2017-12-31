package lostinspace

// return (positions, nil, coords, indices)
func BakeBlockStorageMesh(mesh *Mesh, storage BlockStorage, dic *BlockTypeDictionary) {
	positions := make([]float32, 0)
	//colors := make([]float32, 0)
	coords := make([]float32, 0)
	indices := make([]uint16, 0)

	indexOffset := uint16(0)
	storage.ForEach(func(block *Block) {
		if block.BlockType == BLOCK_TYPE_VOID {
			positions = append(positions,
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			)
			//colors = append(colors,
			//	0, 0, 0,
			//	0, 0, 0,
			//	0, 0, 0,
			//	0, 0, 0,
			//)
			coords = append(coords,
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			)
		} else {
			x := float32(block.coord.X)
			y := float32(block.coord.Y)
			descriptor := dic.Get(block.BlockType)
			layer := float32(0)
			if descriptor != nil {
				layer = float32(descriptor.layerIndex)
			}

			positions = append(positions,
				-0.5+x, 0.5+y, 0,
				-0.5+x, -0.5+y, 0,
				0.5+x, -0.5+y, 0,
				0.5+x, 0.5+y, 0,
			)
			//colors = append(colors,
			//	1, 1, 1,
			//	1, 1, 1,
			//	1, 1, 1,
			//	1, 1, 1,
			//)
			coords = append(coords,
				0, 0, layer,
				0, 1, layer,
				1, 1, layer,
				1, 0, layer,
			)
		}

		indices = append(indices,
			0+indexOffset, 1+indexOffset, 2+indexOffset,
			0+indexOffset, 2+indexOffset, 3+indexOffset,
		)

		indexOffset += 4
	})

	mesh.Positions = positions
	mesh.Colors = nil
	mesh.TexCoords = coords
	mesh.Indices = indices
}

func BakeBlockStorageBody(body *Body, storage BlockStorage, dic *BlockTypeDictionary) {
	storage.ForEach(func(block *Block) {
		if block.BlockType == "" {
			return
		}

		des := dic.Get(block.BlockType)

		vertices := make([]Vec2, len(des.CollisionVertices))
		for i, vertex := range des.CollisionVertices {
			vertices[i] = Vec2{
				vertex.X + float64(block.coord.X),
				vertex.Y + float64(block.coord.Y)}
		}

		body.AddPolygonFixture(des.Density, des.Friction, des.Restitution, vertices)
	})
}
