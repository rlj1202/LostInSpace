package lostinspace

import "github.com/go-gl/gl/v4.1-compatibility/gl"

const (
	ATTRIB_POSITION  = 0
	ATTRIB_COLOR     = 1
	ATTRIB_TEX_COORD = 2
)

// vertex
// (x, y, z, x, y, z, ...)  position
// (r, g, b, r, g, b, ...)  color
// (u, v, l, u, v, l, ...)  texture coordinate
type Mesh struct {
	Positions []float32
	Colors    []float32
	TexCoords []float32
	Indices   []uint16

	elementsCount int32

	ibo uint32
	vao uint32

	positionBuffer uint32
	colorBuffer    uint32
	texCoordBuffer uint32
}

// This doesn't create opengl buffer objects.
// You must call Bake() function to create opengl objects.
func NewMesh(positions, colors, texCoords []float32, indices []uint16) *Mesh {
	mesh := &Mesh{
		Positions: positions,
		Colors:    colors,
		TexCoords: texCoords,
		Indices:   indices,
	}

	return mesh
}

func (mesh *Mesh) Draw() {
	gl.BindVertexArray(mesh.vao)
	gl.DrawElements(gl.TRIANGLES, mesh.elementsCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))
}

func (mesh *Mesh) Bake() {
	if mesh.vao == 0 { // uninitialized
		gl.GenVertexArrays(1, &mesh.vao)
		gl.BindVertexArray(mesh.vao)

		gl.GenBuffers(1, &mesh.positionBuffer)
		gl.GenBuffers(1, &mesh.colorBuffer)
		gl.GenBuffers(1, &mesh.texCoordBuffer)
		gl.GenBuffers(1, &mesh.ibo)
	}

	gl.BindVertexArray(mesh.vao)

	if mesh.Positions != nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.positionBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Positions)*4, gl.Ptr(mesh.Positions), gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(ATTRIB_POSITION)
		gl.VertexAttribPointer(ATTRIB_POSITION, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	} else {
		gl.DisableVertexAttribArray(ATTRIB_POSITION)
	}

	if mesh.Colors != nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.colorBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Colors)*4, gl.Ptr(mesh.Colors), gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(ATTRIB_COLOR)
		gl.VertexAttribPointer(ATTRIB_COLOR, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	} else {
		gl.DisableVertexAttribArray(ATTRIB_COLOR)
	}

	if mesh.TexCoords != nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.texCoordBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(mesh.TexCoords)*4, gl.Ptr(mesh.TexCoords), gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(ATTRIB_TEX_COORD)
		gl.VertexAttribPointer(ATTRIB_TEX_COORD, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	} else {
		gl.DisableVertexAttribArray(ATTRIB_TEX_COORD)
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*2, gl.Ptr(mesh.Indices), gl.DYNAMIC_DRAW)
	mesh.elementsCount = int32(len(mesh.Indices))
}

// TODO Deprecated
func (mesh *Mesh) Set(positions, colors, texCoords []float32, indices []uint16) {
	gl.BindVertexArray(mesh.vao)
	if positions != nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.positionBuffer)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(positions)*4, gl.Ptr(positions))
	}
	if colors != nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.colorBuffer)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(colors)*4, gl.Ptr(colors))
	}
	if texCoords != nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.texCoordBuffer)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(texCoords)*4, gl.Ptr(texCoords))
	}
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ibo)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, len(indices)*2, gl.Ptr(indices))
	mesh.elementsCount = int32(len(indices))
}

func (mesh *Mesh) Destroy() {
	if mesh.positionBuffer != 0 {
		gl.DeleteBuffers(1, &mesh.positionBuffer)
		mesh.positionBuffer = 0
	}
	if mesh.colorBuffer != 0 {
		gl.DeleteBuffers(1, &mesh.colorBuffer)
		mesh.colorBuffer = 0
	}
	if mesh.texCoordBuffer != 0 {
		gl.DeleteBuffers(1, &mesh.texCoordBuffer)
		mesh.texCoordBuffer = 0
	}
	if mesh.ibo != 0 {
		gl.DeleteBuffers(1, &mesh.ibo)
		mesh.ibo = 0
	}
	if mesh.vao != 0 {
		gl.DeleteVertexArrays(1, &mesh.vao)
		mesh.vao = 0
	}
	mesh.elementsCount = 0
}
