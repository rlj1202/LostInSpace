package lostinspace

import (
	"time"

	"github.com/ByteArena/box2d"
)

// World is where bodies are interact.
// Such as chunks, entities, etc.
type World struct {
	b2world *box2d.B2World
}

// Physics body which can collide.
type Body struct {
	b2body *box2d.B2Body
}

type Fixture struct {
	b2fixture *box2d.B2Fixture
}

func NewWorld() *World {
	world := new(World)
	b2world := box2d.MakeB2World(box2d.MakeB2Vec2(0, 0))
	world.b2world = &b2world

	return world
}

func (world *World) Update(dt time.Duration) {
	world.b2world.Step(dt.Seconds(), 8, 3)
}

func (world *World) CreateBody(movable bool) *Body {
	bodyDef := box2d.NewB2BodyDef()
	if movable {
		bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	} else {
		bodyDef.Type = box2d.B2BodyType.B2_staticBody
	}

	b2body := world.b2world.CreateBody(bodyDef)

	body := new(Body)
	body.b2body = b2body

	return body
}

func (body *Body) CreateCircleFixture(density, friction, restitution, radius float64) *Fixture {
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(radius)

	fixDef := box2d.MakeB2FixtureDef()
	fixDef.Density = density
	fixDef.Friction = friction
	fixDef.Restitution = restitution
	fixDef.Shape = &shape

	b2fix := body.b2body.CreateFixtureFromDef(&fixDef)

	fixture := new(Fixture)
	fixture.b2fixture = b2fix

	return fixture
}

func (body *Body) CreatePolygonFixture(density, friction, restitution float64, vertices []Vec2) *Fixture {
	b2vecs := make([]box2d.B2Vec2, len(vertices))
	for i, vec := range vertices {
		b2vecs[i] = box2d.B2Vec2(vec)
	}

	shape := box2d.MakeB2PolygonShape()
	shape.Set(b2vecs, len(b2vecs))

	fixDef := box2d.MakeB2FixtureDef()
	fixDef.Density = density
	fixDef.Friction = friction
	fixDef.Restitution = restitution
	fixDef.Shape = &shape

	b2fix := body.b2body.CreateFixtureFromDef(&fixDef)

	fixture := new(Fixture)
	fixture.b2fixture = b2fix

	return fixture
}

// Destroy all fixtures
func (body *Body) Clear() {
	fixList := make([]*box2d.B2Fixture, 0)
	for fix := body.b2body.GetFixtureList(); fix != nil; fix = fix.GetNext() {
		fixList = append(fixList, fix)
	}

	for _, fix := range fixList {
		body.b2body.DestroyFixture(fix)
	}
}

func (body *Body) GetPosition() (float64, float64) {
	pos := body.b2body.GetPosition()
	return pos.X, pos.Y
}

func (body *Body) GetAngle() float64 {
	return body.b2body.GetAngle()
}

func (body *Body) SetPosition(x, y float64) {
	body.b2body.SetTransform(box2d.MakeB2Vec2(x, y), 0)
}

func (body *Body) SetLinearDamping(damp float64) {
	body.b2body.SetLinearDamping(damp)
}

func (body *Body) ApplyForceToCenter(force Vec2) {
	body.b2body.ApplyForceToCenter(box2d.B2Vec2(force), true)
}

// Destroy body from world.
func (body *Body) Destroy() {
	body.Clear()
	b2world := body.b2body.GetWorld()
	b2world.DestroyBody(body.b2body)
	body.b2body = nil
}
