package lostinspace

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	vs = `
		#version 410

		uniform int translateMode;
		
		uniform mat4 translate;
		uniform mat4 rotate;
		uniform mat4 camera;
		uniform mat4 projection;

		layout(location = 0) in vec3 position;
		layout(location = 1) in vec3 color;
		layout(location = 2) in vec3 texCoord;

		out vec3 fragColor;
		out vec3 fragTexCoord;

		void main() {
			fragColor = color;

			mat4 m = projection * camera * translate * rotate;
			if (translateMode == 0) {
				fragTexCoord = texCoord;
				gl_Position = m * vec4(position, 1);
			} else if (translateMode == 1) {
				fragTexCoord = (m * vec4(texCoord, 1)).xyz;
				gl_Position = vec4(position, 1);
			}
		}
	`
	fs = `
		#version 410

		uniform int texMode;

		uniform sampler2D tex2D;
		uniform sampler2DArray tex2DArray;

		in vec3 fragColor;
		in vec3 fragTexCoord;

		out vec4 finalColor;

		void main() {
			vec4 diffuseColor = vec4(fragColor, 1);
			vec4 texColor;
			if (texMode == 0) {
				texColor = texture(tex2D, fragTexCoord.xy);
			} else if (texMode == 1) {
				texColor = texture(tex2DArray, fragTexCoord);
			}
			//finalColor = mix(texColor, diffuseColor, 0.5);
			finalColor = texColor;
		}
	`
)

type Game struct {
	window  *Window
	world   *World
	terrain *Terrain
	dic     *BlockTypeDictionary
	player  *Player
	camera  *Camera

	seed *Seed

	loadSectorQueue chan WorldSectorCoord
	bakeSectorQueue chan WorldSectorCoord

	chunksToDraw map[WorldChunkCoord]*Chunk
	chunksRadius int
	quit         chan bool

	// maximum number of loaded sectors
	maxSectors    int
	sectorsToLoad []WorldSectorCoord

	shader *ShaderProgram
	entity *BlockEntity

	msaaTex uint32
	msaaFbo uint32

	bgTex   Texture
	quad    *Mesh
	bgRatio float32
	bgHu    float32
	bgHv    float32
}

func NewGame(window *Window, dic *BlockTypeDictionary) *Game {
	for _, desc := range dic.data {
		log.Printf("%s\n", desc)
	}

	playerTexFile, err := os.Open("player.png")
	if err != nil {
		panic(err)
	}
	playerTex := NewTexture2D(playerTexFile)

	game := new(Game)
	game.window = window
	game.world = NewWorld()
	game.terrain = NewTerrain()
	game.dic = dic
	game.dic.arrayTex.Bind(1)
	game.player = NewPlayer(game.world, playerTex)
	game.player.Mesh.Bake()
	width, height := glfw.GetCurrentContext().GetSize()
	game.camera = NewCamera(20, 20*float64(height)/float64(width))
	game.camera.SetTarget(game.player.Body)

	game.loadSectorQueue = make(chan WorldSectorCoord, 9)
	game.bakeSectorQueue = make(chan WorldSectorCoord, 9)

	seed := NewSeed(2)
	game.seed = seed
	game.loadSectorQueue <- WorldSectorCoord{0, 0}
	game.loadSectorQueue <- WorldSectorCoord{-1, 0}
	game.loadSectorQueue <- WorldSectorCoord{0, -1}
	game.loadSectorQueue <- WorldSectorCoord{-1, -1}

	game.quit = make(chan bool)
	game.chunksRadius = 9
	go func() {
		for {
			select {
			case <-game.quit:
				return
			default:
				posX, posY := game.camera.target.GetPosition()
				radius := game.chunksRadius

				chunks := make(map[WorldChunkCoord]*Chunk)
				for y := 0; y < radius*2; y++ {
					for x := 0; x < radius*2; x++ {
						worldChunkCoord := WorldChunkCoord{
							int64(posX)/CHUNK_WIDTH + int64(x-radius),
							int64(posY)/CHUNK_HEIGHT + int64(y-radius),
						}
						chunk := game.terrain.GetChunk(worldChunkCoord)
						if chunk == nil {
							continue
						}
						chunks[worldChunkCoord] = chunk
					}
				}
				game.chunksToDraw = chunks

				time.Sleep(time.Second * 2)
			}
		}
	}()
	game.maxSectors = 9
	go func() {
		var curSector, prevSector WorldSectorCoord

		for {
			select {
			case <-game.quit:
				return
			default:
				x, y := game.player.GetPosition()
				worldBlockCoord := WorldBlockCoord{int64(math.Floor(x)), int64(math.Floor(y))}
				curSector, _, _ = worldBlockCoord.Parse()

				if curSector != prevSector {
					game.loadSectorQueue <- curSector
					game.loadSectorQueue <- curSector.Left()
					game.loadSectorQueue <- curSector.Left().Up()
					game.loadSectorQueue <- curSector.Left().Down()
					game.loadSectorQueue <- curSector.Right()
					game.loadSectorQueue <- curSector.Right().Up()
					game.loadSectorQueue <- curSector.Right().Down()
					game.loadSectorQueue <- curSector.Up()
					game.loadSectorQueue <- curSector.Down()

					prevSector = curSector
				}

				time.Sleep(time.Second * 2)
			}
		}
	}()

	go func() { // sector loader: generate sector or load from file.
		for {
			select {
			case <-game.quit:
				return
			case coord := <-game.loadSectorQueue:
				sector := game.terrain.GetSector(coord)
				if sector != nil {
					continue
				}

				// Load sector from file.
				sector = LoadSector(coord)
				if sector == nil {
					// Generate sector if there is no file.
					sector = GenerateSector(game.seed, coord)
				}

				game.terrain.SetSector(sector)

				for _, chunk := range sector.Chunks {
					worldChunkCoord := CombineWorldChunkCoord(coord, chunk.coord)

					game.terrain.BakeChunk(game.world, game.dic, worldChunkCoord)
				}

				game.bakeSectorQueue <- coord
			}
		}
	}()

	gl.GenTextures(1, &game.msaaTex)
	gl.BindTexture(gl.TEXTURE_2D_MULTISAMPLE, game.msaaTex)
	gl.TexImage2DMultisample(gl.TEXTURE_2D_MULTISAMPLE, 4, gl.RGBA8, int32(width), int32(height), false)

	gl.GenFramebuffers(1, &game.msaaFbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, game.msaaFbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D_MULTISAMPLE, game.msaaTex, 0)

	bgTexFile, err := os.Open("bg_starfield.png")
	if err != nil {
		panic(err)
	}
	game.bgTex = NewTexture2D(bgTexFile)
	imageW, imageH := game.bgTex.GetSize()
	frameW, frameH := glfw.GetCurrentContext().GetSize()
	hu, hv := float32(frameW)/float32(imageW)/2.0, float32(frameH)/float32(imageH)/2.0
	game.bgHu = hu
	game.bgHv = hv
	game.bgRatio = float32(imageH) / float32(imageW)
	game.quad = NewMesh(
		[]float32{
			-1, 1, 0,
			-1, -1, 0,
			1, -1, 0,
			1, 1, 0,
		},
		nil,
		[]float32{
			-hu, -hv, 0,
			-hu, hv, 0,
			hu, hv, 0,
			hu, -hv, 0,
		},
		[]uint16{0, 1, 2, 0, 2, 3},
	)
	game.quad.Bake()

	game.shader = NewShaderProgram(vs, fs)
	game.shader.UniformInt("tex2D", 0)
	game.shader.UniformInt("tex2DArray", 1)
	game.shader.UniformMat4("projection", game.camera.GetProjectionMat())
	game.shader.UniformMat4("camera", mgl32.Ident4())
	game.shader.UniformMat4("translate", mgl32.Ident4())
	game.shader.UniformMat4("rotate", mgl32.Ident4())
	game.shader.UniformMat4("texTranslate", mgl32.Ident4())
	game.shader.Bind()

	RegisterEventListener(game)

	game.entity = NewBlockEntity(game.world)
	game.entity.Set(NewBlock(BlockCoord{0, 0}, "stone"))
	game.entity.Set(NewBlock(BlockCoord{1, 1}, "stone"))
	game.entity.Set(NewBlock(BlockCoord{2, 0}, "stone"))
	game.entity.Bake(game.world, dic)

	return game
}

func (game *Game) Update(dt time.Duration) {
	PollEvents()

	keyA := GetKeyActionState(KEY_A)
	keyD := GetKeyActionState(KEY_D)
	keyW := GetKeyActionState(KEY_W)
	keyS := GetKeyActionState(KEY_S)
	if keyA == ACTION_PRESS || keyA == ACTION_REPEAT {
		game.player.ApplyForceToCenter(Vec2{-40, 0})
	}
	if keyD == ACTION_PRESS || keyD == ACTION_REPEAT {
		game.player.ApplyForceToCenter(Vec2{40, 0})
	}
	if keyW == ACTION_PRESS || keyW == ACTION_REPEAT {
		game.player.ApplyForceToCenter(Vec2{0, 40})
	}
	if keyS == ACTION_PRESS || keyS == ACTION_REPEAT {
		game.player.ApplyForceToCenter(Vec2{0, -40})
	}

	select {
	case coord := <-game.bakeSectorQueue:
		sector := game.terrain.GetSector(coord)

		start := time.Now()
		for _, chunk := range sector.Chunks {
			chunk.mesh.Bake()
			chunk.body.Bake()
		}
		elapsed := time.Since(start)

		log.Printf("Bake sector: %v, %v seconds\n", coord, elapsed.Seconds())
	default:
		break
	}

	game.world.Update(dt)
}

func (game *Game) Destroy() {
	close(game.quit)
}

func (game *Game) Render() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// render to the fbo for multisampling
	gl.BindFramebuffer(gl.FRAMEBUFFER, game.msaaFbo)

	game.shader.UniformMat4("rotate", mgl32.Ident4())

	// render background
	game.renderBackground()

	// init uniforms
	game.shader.UniformInt("texMode", 0)
	game.shader.UniformInt("translateMode", 0)
	game.shader.UniformMat4("projection", game.camera.GetProjectionMat())
	game.shader.UniformMat4("camera", game.camera.GetCameraMat())

	// render player
	game.renderPlayer()

	// render chunks
	game.renderChunks()

	// render an entity
	x, y := game.entity.Body.GetPosition()
	angle := game.entity.Body.GetAngle()
	game.shader.UniformMat4("translate", mgl32.Translate3D(
		float32(x),
		float32(y),
		0,
	))
	game.shader.UniformMat4("rotate", mgl32.HomogRotate3DZ(float32(angle)))
	game.entity.Mesh.Draw()

	// completed rendering
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// resolve result to screen
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, game.msaaFbo)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	width, height := glfw.GetCurrentContext().GetSize()
	gl.BlitFramebuffer(
		0, 0, int32(width), int32(height),
		0, 0, int32(width), int32(height), gl.COLOR_BUFFER_BIT, gl.NEAREST)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (game *Game) renderPlayer() {
	game.player.Texture.Bind(0)
	x, y := game.player.GetPosition()
	game.shader.UniformMat4("translate", mgl32.Translate3D(
		float32(x),
		float32(y),
		0,
	))
	game.player.Mesh.Draw()
}

func (game *Game) renderChunks() {
	game.shader.UniformInt("texMode", 1)
	for worldChunkCoord, chunk := range game.chunksToDraw {
		game.shader.UniformMat4("translate", mgl32.Translate3D( // TODO 이 행렬을 캐쉬해두면 속도가 빨라질듯
			float32(worldChunkCoord.X*CHUNK_WIDTH),
			float32(worldChunkCoord.Y*CHUNK_HEIGHT),
			0,
		))
		chunk.mesh.Draw()
	}
}

func (game *Game) renderBackground() {
	zoom := float32(game.camera.GetZoom())
	game.shader.UniformMat4("projection", mgl32.Scale3D(1.0/zoom, 1.0/zoom, 1))
	a, b := game.camera.target.GetPosition()
	game.shader.UniformMat4("camera", mgl32.Translate3D(
		float32(a)/1500.0/game.bgHu*zoom,
		float32(-b)/1500.0/game.bgHu*zoom/game.bgRatio,
		0,
	))
	game.shader.UniformMat4("translate", mgl32.Ident4())
	game.shader.UniformInt("texMode", 0)
	game.shader.UniformInt("translateMode", 1)
	game.bgTex.Bind(0)
	game.quad.Draw()
}

// Load chunk to given coordinates.
// If there is no chunk generated, generate it.
// If there is generated chunk, load it from file and bake it.
/*
func (game *Game) loadChunk(seed *Seed, worldChunkCoord WorldChunkCoord) {
	sectorCoord, chunkCoord := worldChunkCoord.Parse()
	log.Printf("Load chunk: %s -> %s, %s\n", worldChunkCoord, sectorCoord, chunkCoord)

	chunk := game.terrain.GetChunk(worldChunkCoord)
	if chunk != nil {
		chunk.Destroy()
	}

	chunk = GenerateChunk(seed, worldChunkCoord)
	game.terrain.SetChunk(worldChunkCoord, chunk)
	game.terrain.BakeChunk(game.world, game.dic, worldChunkCoord)
}
*/

/*
func (game *Game) loadSector(seed *Seed, worldSectorCoord WorldSectorCoord) {
	sector := game.terrain.GetSector(worldSectorCoord)
	if sector != nil {
		return
	}

	// try to load from file. if there is no file, generate it.
	sector = LoadSector(worldSectorCoord)
	if sector == nil {
		sector = GenerateSector(seed, worldSectorCoord)
	}

	game.terrain.SetSector(sector)
	for _, chunk := range sector.Chunks {
		game.terrain.BakeChunk(game.world, game.dic, CombineWorldChunkCoord(worldSectorCoord, chunk.coord))
	}
}
*/

func (game *Game) unloadSector(worldSectorCoord WorldSectorCoord) {
	log.Printf("unload sector: %s\n", worldSectorCoord)

	sector := game.terrain.GetSector(worldSectorCoord)
	if sector != nil {
		sector.Destroy()
		game.terrain.DeleteSector(worldSectorCoord)
	}
}

func (game Game) OnEvent(event Event) {
	switch event.(type) {
	case MouseEvent:
		mouseEvent := event.(MouseEvent)
		action := mouseEvent.Action
		if action != ACTION_PRESS {
			break
		}
		button := mouseEvent.Button
		xpos := float32(mouseEvent.XPos)
		ypos := float32(mouseEvent.YPos)

		width, height := game.window.GetSize()
		worldPos, err := mgl32.UnProject(
			mgl32.Vec3{xpos, float32(height) - ypos, 0},
			game.camera.GetCameraMat(),
			game.camera.GetProjectionMat(),
			0, 0,
			width, height,
		)
		if err == nil {
			log.Printf("worldPos: %v\n", worldPos)
			worldCoord := WorldBlockCoord{
				int64(math.Floor(float64(worldPos.X()) + 0.5)),
				int64(math.Floor(float64(worldPos.Y()) + 0.5)),
			}
			sectorCoord, chunkCoord, blockCoord := worldCoord.Parse()

			worldChunkCoord := CombineWorldChunkCoord(sectorCoord, chunkCoord)
			chunk := game.terrain.GetChunk(worldChunkCoord)
			if chunk == nil {
				break
			}

			if button == MOUSE_BUTTON_LEFT {
				chunk.Set(NewBlock(blockCoord, ""))
				chunk.body.Clear()
				game.terrain.BakeChunk(game.world, game.dic, worldChunkCoord)
				chunk.mesh.Bake()
				chunk.body.Bake()
			} else if button == MOUSE_BUTTON_RIGHT {
				chunk.Set(NewBlock(blockCoord, "stone"))
				chunk.body.Clear()
				game.terrain.BakeChunk(game.world, game.dic, worldChunkCoord)
				chunk.mesh.Bake()
				chunk.body.Bake()
			}
		}
	case ScrollEvent:
		scrollEvent := event.(ScrollEvent)
		yoff := scrollEvent.YOff

		camera := game.camera
		camera.SetZoom(camera.GetZoom() + yoff/20.0)
	}
}
