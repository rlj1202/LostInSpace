package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	vertexShaderSource = `
		#version 410

		uniform mat4 projectionMat;
		uniform mat4 cameraMat;
		uniform mat4 translateMat;

		in vec3 vp;
		in vec2 vpTexCoord;

		out vec2 fragTexCoord;

		void main() {
			fragTexCoord = vpTexCoord;
			gl_Position = projectionMat * cameraMat * translateMat * vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410

		uniform sampler2DArray tex;
		uniform int layer;

		in vec2 fragTexCoord;

		out vec4 frag_colour;

		void main() {
			frag_colour = texture(tex, vec3(fragTexCoord, layer));
		}
	` + "\x00"
)

var (
	quad = []float32{
		-0.5, 0.5, 0, 0, 0,
		-0.5, -0.5, 0, 0, 1,
		0.5, -0.5, 0, 1, 1,

		-0.5, 0.5, 0, 0, 0,
		0.5, -0.5, 0, 1, 1,
		0.5, 0.5, 0, 1, 0,
	}
)

type GameRenderer struct {
	*Game

	program uint32

	projectionMat mgl32.Mat4
	cameraMat     mgl32.Mat4
	translateMat  mgl32.Mat4

	projectionMatLoc int32
	cameraMatLoc     int32
	translateMatLoc  int32
	textureLoc       int32
	layerLoc         int32

	blockTexture uint32
	vao          uint32

	blockTextureLayerIndex map[BlockType]int32
}

func (renderer *GameRenderer) Init(width, height uint64, blockTypeDescriptors []*BlockTypeDescriptor) {
	renderer.program = initOpenGL()
	renderer.blockTextureLayerIndex = make(map[BlockType]int32)

	renderer.projectionMatLoc = gl.GetUniformLocation(renderer.program, gl.Str("projectionMat\x00"))
	renderer.cameraMatLoc = gl.GetUniformLocation(renderer.program, gl.Str("cameraMat\x00"))
	renderer.translateMatLoc = gl.GetUniformLocation(renderer.program, gl.Str("translateMat\x00"))
	renderer.textureLoc = gl.GetUniformLocation(renderer.program, gl.Str("tex\x00"))
	renderer.layerLoc = gl.GetUniformLocation(renderer.program, gl.Str("layer\x00"))

	gl.ProgramUniform1i(renderer.program, renderer.textureLoc, 0)

	renderer.vao = makeVao(renderer.program, quad)
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
	gl.UseProgram(renderer.program)

	cameraWidthHalf := float32(renderer.Camera.Width) / 2.0
	cameraHeightHalf := float32(renderer.Camera.Height) / 2.0
	renderer.projectionMat = mgl32.Ortho(
		-cameraWidthHalf, cameraWidthHalf,
		-cameraHeightHalf, cameraHeightHalf,
		float32(5.0), float32(-5.0))

	renderer.cameraMat = mgl32.Translate3D(float32(-renderer.Game.Camera.Position[0]), float32(-renderer.Game.Camera.Position[1]), 0)
	gl.ProgramUniformMatrix4fv(renderer.program, renderer.projectionMatLoc, 1, false, &renderer.projectionMat[0])
	gl.ProgramUniformMatrix4fv(renderer.program, renderer.cameraMatLoc, 1, false, &(renderer.cameraMat[0]))

	renderer.renderTerrain(renderer.Game.Terrain)
	renderer.renderObjects()
}

func (renderer *GameRenderer) renderTerrain(terrain *Terrain) {
	gl.BindVertexArray(renderer.vao)

	for _, chunk := range terrain.Chunks {
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
	}
}

func (renderer *GameRenderer) renderObjects() {
	entities := renderer.Game.Entities
	for _, entity := range entities {
		position := entity.GetPosition()

		gl.ProgramUniform1i(renderer.program, renderer.layerLoc, 0)

		renderer.translateMat = mgl32.Translate3D(float32(position.X), float32(position.Y), 0)
		gl.ProgramUniformMatrix4fv(renderer.program, renderer.translateMatLoc, 1, false, &(renderer.translateMat[0]))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quad)/3))
	}
}

func toGlPos(x, y float32, projectionMat, transformationMat mgl32.Mat4) (float32, float32) {
	coord := projectionMat.Mul4(transformationMat).Mul4x1(mgl32.NewVecNFromData([]float32{x, y, 0, 1}).Vec4())
	return coord.X(), coord.Y()
}

func toCoord(x, y float32, projectionMat, transformationMat mgl32.Mat4) (float32, float32) {
	coord := projectionMat.Mul4(transformationMat).Inv().Mul4x1(mgl32.NewVecNFromData([]float32{x, y, 0, 1}).Vec4())
	return coord.X(), coord.Y()
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
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

	vpAttrib := uint32(gl.GetAttribLocation(prog, gl.Str("vp\x00")))
	gl.EnableVertexAttribArray(vpAttrib)
	gl.VertexAttribPointer(vpAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	vpTexCoordAttrib := uint32(gl.GetAttribLocation(prog, gl.Str("vpTexCoord\x00")))
	gl.EnableVertexAttribArray(vpTexCoordAttrib)
	gl.VertexAttribPointer(vpTexCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

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
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

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
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	return textureID, nil
}

// initOpenGL initializes OpenGL and returns an initialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	// Enable Anti-aliasing
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.POLYGON_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.Enable(gl.MULTISAMPLE)
	gl.Hint(gl.POLYGON_SMOOTH_HINT, gl.NICEST)
	gl.Hint(gl.LINE_SMOOTH_HINT, gl.NICEST)

	return prog
}
