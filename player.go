package main

import (
	"os"

	"github.com/ByteArena/box2d"
)

// Player
type Player struct {
	Health int
	*Inventory

	vbo           uint32
	vao           uint32
	verticesCount int32

	*box2d.B2Body
}

type Inventory struct {
	Items []*Item
}

type Item struct {
	ItemType
}

type ItemType string

type ItemTypeDescriptor struct {
	ItemType
	Name        string
	TextureFile *os.File
}

func NewPlayer(game *Game) *Player {
	bodyDef := box2d.NewB2BodyDef()
	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody

	body := game.World.B2World.CreateBody(bodyDef)

	fixDef := box2d.MakeB2FixtureDef()
	circleShape := box2d.NewB2CircleShape()
	circleShape.SetRadius(0.5)
	fixDef.Shape = circleShape
	fixDef.Density = 1.0
	fixDef.Friction = 0.2
	fixDef.Restitution = 0.2

	body.CreateFixtureFromDef(&fixDef)
	body.SetLinearDamping(2.0)

	return &Player{
		Health:    100,
		Inventory: NewInventory(),
		B2Body:    body,
	}
}

func NewInventory() *Inventory {
	return &Inventory{
		Items: make([]*Item, 100),
	}
}
