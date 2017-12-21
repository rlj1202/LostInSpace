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
	elementsCount int32

	ibo uint32
	vao uint32

	positionBuffer uint32
	colorBuffer    uint32
	texCoordBuffer uint32
}

func NewMesh(positions, colors, texCoords []float32, indices []uint16) *Mesh {
	mesh := &Mesh{}

	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	if positions != nil {
		gl.GenBuffers(1, &mesh.positionBuffer)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.positionBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(positions)*4, gl.Ptr(positions), gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(ATTRIB_POSITION)
		gl.VertexAttribPointer(ATTRIB_POSITION, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	} else {
		gl.DisableVertexAttribArray(ATTRIB_POSITION)
	}

	if colors != nil {
		gl.GenBuffers(1, &mesh.colorBuffer)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.colorBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(ATTRIB_COLOR)
		gl.VertexAttribPointer(ATTRIB_COLOR, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	} else {
		gl.DisableVertexAttribArray(ATTRIB_COLOR)
	}

	if texCoords != nil {
		gl.GenBuffers(1, &mesh.texCoordBuffer)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.texCoordBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(texCoords)*4, gl.Ptr(texCoords), gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(ATTRIB_TEX_COORD)
		gl.VertexAttribPointer(ATTRIB_TEX_COORD, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	} else {
		gl.DisableVertexAttribArray(ATTRIB_TEX_COORD)
	}

	gl.GenBuffers(1, &mesh.ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*2, gl.Ptr(indices), gl.DYNAMIC_DRAW)
	mesh.elementsCount = int32(len(indices))

	return mesh
}

func (mesh *Mesh) Bind() {
	gl.BindVertexArray(mesh.vao)
}

func (mesh *Mesh) Draw() {
	mesh.Bind()
	gl.DrawElements(gl.TRIANGLES, mesh.elementsCount, gl.UNSIGNED_SHORT, gl.PtrOffset(0))
}

func (mesh *Mesh) Set(positions, colors, texCoords []float32, indices []uint16) {
	mesh.Bind()
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
