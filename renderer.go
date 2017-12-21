package lostinspace

type RendererTest interface { // TODO
	Bind()
	Render()
}

type Renderer struct {
	shaderProgram    *ShaderProgram
	meshesPerTexture map[Texture][]*Mesh
}

func NewRenderer(vertexShaderRaw, fragmentShaderRaw string) *Renderer {
	renderer := new(Renderer)
	renderer.shaderProgram = NewShaderProgram(vertexShaderRaw, fragmentShaderRaw)
	renderer.meshesPerTexture = make(map[Texture][]*Mesh)

	return renderer
}

func (renderer *Renderer) Bind() {
	renderer.shaderProgram.Bind()
}

func (renderer *Renderer) Add(texture Texture, mesh *Mesh) { // TODO
	meshes, exist := renderer.meshesPerTexture[texture]
	if !exist {
		meshes = make([]*Mesh, 0, 10)
	}

	renderer.meshesPerTexture[texture] = append(meshes, mesh)
}

func (renderer *Renderer) Render() { // TODO
	for texture, meshes := range renderer.meshesPerTexture {
		texture.Bind(0)

		for _, mesh := range meshes {
			mesh.Draw()
		}
	}
}
