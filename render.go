package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	vertexShaderSource = `
		#version 410

		uniform mat4 projectionMat;
		uniform mat4 cameraMat;
		uniform mat4 translateMat;
		uniform mat4 rotateMat;

		layout(location = 0) in vec3 vp;
		layout(location = 1) in vec2 vpTexCoord;
		layout(location = 2) in float layer;

		out vec3 texCoord;

		void main() {
			texCoord = vec3(vpTexCoord, layer);
			
			gl_Position = projectionMat * cameraMat * translateMat * rotateMat * vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410

		uniform sampler2DArray tex;

		in vec3 texCoord;

		out vec4 fragColor;

		void main() {
			fragColor = texture(tex, texCoord);
		}
	` + "\x00"
)

var (
	// temporal vao
	quad = []float32{
		-0.5, 0.5, 0, 0, 0, 0,
		-0.5, -0.5, 0, 0, 1, 0,
		0.5, -0.5, 0, 1, 1, 0,

		-0.5, 0.5, 0, 0, 0, 0,
		0.5, -0.5, 0, 1, 1, 0,
		0.5, 0.5, 0, 1, 0, 0,
	}
)

type GameRenderer struct {
	*Game

	shaderProgram *Shader

	projectionMat mgl32.Mat4
	cameraMat     mgl32.Mat4
	translateMat  mgl32.Mat4
	rotateMat     mgl32.Mat4

	blockTexture uint32
	vao          uint32

	blockTextureLayerIndex map[BlockType]int32
}

type Shader struct {
	program uint32

	vertexShader   uint32
	fargmentShader uint32
}

type TerrainRenderer struct { // TODO
	*Shader

	blockTexture uint32
}

func (renderer *GameRenderer) GetTextureIndex(blockType BlockType) (int32, bool) {
	index, exist := renderer.blockTextureLayerIndex[blockType]
	return index, exist
}

func (renderer *GameRenderer) ToWorldCoord(x, y float32) (float32, float32) {
	width, height := glfw.GetCurrentContext().GetSize()
	y = float32(height) - y
	x = x*2.0/float32(width) - 1.0
	y = y*2.0/float32(height) - 1.0
	screenCoord := mgl32.NewVecNFromData([]float32{x, y, 0, 1}).Vec4()
	m := renderer.projectionMat.Mul4(renderer.cameraMat)
	mInv := m.Inv()
	worldCoord := mInv.Mul4x1(screenCoord)
	return worldCoord.X(), worldCoord.Y()
}

func (renderer *GameRenderer) ToScreenCoord(x, y float32) (float32, float32) {
	worldCoord := mgl32.NewVecNFromData([]float32{x, y, 0, 1}).Vec4()
	m := renderer.projectionMat.Mul4(renderer.cameraMat)
	screenCoord := m.Mul4x1(worldCoord)
	return screenCoord.X(), screenCoord.Y()
}

func (renderer *GameRenderer) Init(width, height uint64, blockTypeDescriptors []*BlockTypeDescriptor) {
	initOpenGL()
	program, err := NewShaderProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}
	renderer.shaderProgram = program
	renderer.blockTextureLayerIndex = make(map[BlockType]int32)

	program.UniformInt("tex", 0)

	renderer.vao = makeVao(renderer.shaderProgram.program, quad)
	blockTexture, err := renderer.loadBlockTextures(blockTypeDescriptors)
	if err != nil {
		panic(err)
	}
	renderer.blockTexture = blockTexture
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, renderer.blockTexture)
}

func (renderer *GameRenderer) RenderGame(dt time.Duration) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	program := renderer.shaderProgram
	program.Bind()

	cameraWidthHalf := float32(renderer.Camera.Width) / 2.0 / float32(renderer.Camera.Zoom)
	cameraHeightHalf := float32(renderer.Camera.Height) / 2.0 / float32(renderer.Camera.Zoom)
	renderer.projectionMat = mgl32.Ortho(
		-cameraWidthHalf, cameraWidthHalf,
		-cameraHeightHalf, cameraHeightHalf,
		float32(5.0), float32(-5.0))

	renderer.cameraMat = mgl32.Translate3D(
		float32(-renderer.Game.Camera.GetPosition().X),
		float32(-renderer.Game.Camera.GetPosition().Y),
		0,
	)
	program.UniformMatrix4("projectionMat", renderer.projectionMat)
	program.UniformMatrix4("cameraMat", renderer.cameraMat)
	renderer.renderTerrain(renderer.Game.Terrain)
	renderer.renderObjects()
}

func (renderer *GameRenderer) renderTerrain(terrain *Terrain) {
	program := renderer.shaderProgram
	for _, chunk := range terrain.Chunks {
		renderer.translateMat = mgl32.Translate3D(float32(chunk.X*16), float32(chunk.Y*16), 0)
		renderer.rotateMat = mgl32.HomogRotate3DZ(0)
		program.UniformMatrix4("translateMat", renderer.translateMat)
		program.UniformMatrix4("rotateMat", renderer.rotateMat)

		chunk.BlockContainerObject.Draw()
		/*
			for _, block := range chunk.Blocks {
				blockX := float32(int64(block.X) + chunk.X*16)
				blockY := float32(int64(block.Y) + chunk.Y*16)

				renderer.translateMat = mgl32.Translate3D(blockX, blockY, 0)

				//scale := float32(1.0)
				//scaleMat := mgl32.Scale3D(scale, scale, scale)
				//transformationMat = scaleMat.Mul4(transformationMat)

				//minX, minY := toCoord(0, 0, projectionMat, transformationMat)
				//maxX, maxY := toCoord(width, height, projectionMat, transformationMat)

				//if blockX < minX || maxX < blockX || blockY < minY || maxY < blockY {
				//	continue
				//}

				if block.BlockType == "" {
					continue
				}
				layer := renderer.blockTextureLayerIndex[block.BlockType]
				gl.ProgramUniform1i(renderer.program, renderer.layerLoc, layer)

				gl.ProgramUniformMatrix4fv(renderer.program, renderer.translateMatLoc, 1, false, &(renderer.translateMat[0]))

				gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))
			}
		*/
	}
}

func (renderer *GameRenderer) renderObjects() {
	//gl.BindVertexArray(renderer.vao)
	program := renderer.shaderProgram

	entities := renderer.Game.Entities
	for _, entity := range entities {
		body := entity.BlockContainerObject.body

		rotation := body.GetAngle()
		position := body.GetPosition()

		renderer.translateMat = mgl32.Translate3D(
			float32(position.X),
			float32(position.Y),
			float32(0),
		)
		renderer.rotateMat = mgl32.HomogRotate3DZ(float32(rotation))

		program.UniformMatrix4("translateMat", renderer.translateMat)
		program.UniformMatrix4("rotateMat", renderer.rotateMat)

		entity.Draw()
	}
}

func NewShaderProgram(vertexShaderRaw, fragmentShaderRaw string) (*Shader, error) {
	vertexShader, err := compileShader(vertexShaderRaw, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fragmentShader, err := compileShader(fragmentShaderRaw, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return &Shader{program, vertexShader, fragmentShader}, nil
}

func (shader *Shader) Bind() {
	gl.UseProgram(shader.program)
}

func (shader *Shader) GetUniformLoc(name string) int32 {
	return gl.GetUniformLocation(shader.program, gl.Str(name+"\x00"))
}

func (shader *Shader) UniformInt(name string, value int32) {
	loc := shader.GetUniformLoc(name)
	gl.ProgramUniform1i(shader.program, loc, value)
}

func (shader *Shader) UniformMatrix4(name string, mat mgl32.Mat4) {
	loc := shader.GetUniformLoc(name)
	gl.ProgramUniformMatrix4fv(shader.program, loc, 1, false, &(mat[0]))
}

func compileShader(raw string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	source, free := gl.Strs(raw)
	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(prog uint32, points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 6*4, gl.PtrOffset(5*4))
	return vao
}

// Instead of using texture atlas, i decided to use texture array because i think it is more easy to use and intuitive.
// So, before game starts, load block textures and store layer index number to each block type. And use that layer index number when draw a block.
func (renderer *GameRenderer) loadBlockTextures(blockTypeDescriptors []*BlockTypeDescriptor) (uint32, error) {
	textureID := uint32(0)
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, textureID)

	mipmapLevel := int32(0)

	texWidth := int32(16)
	texHeight := int32(16)

	gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, 1, gl.RGBA8, texWidth, texHeight, int32(len(blockTypeDescriptors)))
	for i, blockTypeDescriptor := range blockTypeDescriptors {
		imgFile := blockTypeDescriptor.TextureFile
		img, _, err := image.Decode(imgFile)
		if err != nil {
			return 0, err
		}

		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			return 0, fmt.Errorf("Unsupported stride.")
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, mipmapLevel, 0, 0, int32(i), texWidth, texHeight, 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

		renderer.blockTextureLayerIndex[blockTypeDescriptor.BlockType] = int32(i)
	}

	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	gl.BindTexture(gl.TEXTURE_2D_ARRAY, 0)

	return textureID, nil
}

func loadTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("Unsupported stride.")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	textureID := uint32(0)
	gl.GenTextures(1, &textureID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0,
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return textureID, nil
}

// initOpenGL initializes OpenGL and returns an initialized program.
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl version", version)

	gl.ClearColor(0.0, 0.0, 0.1, 1.0)

	// Enable Anti-aliasing
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.POLYGON_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Enable(gl.MULTISAMPLE)
	gl.Hint(gl.POLYGON_SMOOTH_HINT, gl.NICEST)
	gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)
}
