package lostinspace

type Player struct {
	*Mesh
	Texture
	*Body
}

func NewPlayer(world *World, texture Texture) *Player {
	player := new(Player)
	player.Mesh = NewMesh(
		[]float32{
			-0.5, 0.5, 0,
			-0.5, -0.5, 0,
			0.5, -0.5, 0,
			0.5, 0.5, 0,
		},
		nil,
		[]float32{
			0, 0, 0,
			0, 1, 0,
			1, 1, 0,
			1, 0, 0,
		},
		[]uint16{0, 1, 2, 0, 2, 3},
	)
	player.Texture = texture
	player.Body = world.CreateBody(true)
	player.Body.AddCircleFixture(1.0, 0.2, 0.2, 0.49)
	player.Body.SetLinearDamping(2.0)
	player.Body.Bake()

	return player
}
