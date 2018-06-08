package lostinspace

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// BlockObject
type BlockObject struct { // TODO
	storage BlockStorage
	world   *World

	mesh *Mesh

	mainBody  *Body
	subBodies []*Body
}

func NewBlockObject(storage BlockStorage, world *World, bodyType BodyType) *BlockObject {
	obj := new(BlockObject)

	obj.storage = storage
	obj.world = world

	obj.mesh = NewMesh(nil, nil, nil, nil)

	obj.mainBody = world.CreateBody(bodyType)
	obj.subBodies = make([]*Body, 0)

	return obj
}

// Build method create mesh vertices and physical bodies and fixtures.
// This method can be called in any goroutine.
func (obj *BlockObject) Build() {
}

// Bake method calls opengl methods and box2d methods
// so this method must be called in main thread, not in goroutine.
func (obj *BlockObject) Bake() {
	obj.mesh.Bake()
	obj.mainBody.Bake()
	for _, subBody := range obj.subBodies {
		subBody.Bake()
	}
}

// Destroy method removes opengl objects and box2d bodies
// so this method must be called in main thead, not in goroutine.
func (obj *BlockObject) Destroy() {
	obj.mesh.Destroy()
	obj.mainBody.Destroy()
	for _, subBody := range obj.subBodies {
		subBody.Destroy()
	}
}

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
			coordVecs := []Vec2{
				{-0.5, -0.5},
				{-0.5, 0.5},
				{0.5, 0.5},
				{0.5, -0.5},
			}
			for _, coordVec := range coordVecs {
				rotatedVec := mgl32.Rotate2D(float32(block.FrontFace) * math.Pi / 2.0).Mul2x1(mgl32.Vec2{float32(coordVec.X), float32(coordVec.Y)})
				rotatedVec[0] += 0.5
				rotatedVec[1] += 0.5

				coords = append(coords,
					rotatedVec[0], rotatedVec[1], layer,
				)
			}
			/*
				coords = append(coords,
					0, 0, layer,
					0, 1, layer,
					1, 1, layer,
					1, 0, layer,
				)
			*/
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
			rotatedVertex := mgl32.Rotate2D(float32(block.FrontFace) * math.Pi / 2.0).Mul2x1(mgl32.Vec2{float32(vertex.X), float32(vertex.Y)})

			vertices[i] = Vec2{
				float64(rotatedVertex[0]) + float64(block.coord.X),
				float64(rotatedVertex[1]) + float64(block.coord.Y),
			}
		}

		/*
			if !block.Fixed {
				singleBody := body.world.CreateBody(true)
				singleBody.AddPolygonFixture(des.Density, des.Friction, des.Restitution, vertices)
			}
		*/

		body.AddPolygonFixture(des.Density, des.Friction, des.Restitution, vertices)
	})
}
