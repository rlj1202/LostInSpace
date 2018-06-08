package lostinspace

import (
	"time"

	"github.com/rlj1202/box2d"
)

const (
	STATIC BodyType = iota
	DYNAMIC
	KINEMATIC
)

type BodyType int

// World is where bodies are interact.
// Such as chunks, entities, etc.
type World struct {
	b2world *box2d.B2World
}

// Physics body which can collide.
type Body struct {
	world *World

	bodyDef *box2d.B2BodyDef
	fixDefs []*box2d.B2FixtureDef

	b2body *box2d.B2Body
}

type Joint struct { // TODO
	world *World

	jointDef box2d.B2JointDefInterface
	b2joint  box2d.B2JointInterface
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

func (world *World) CreateBody(bodyType BodyType) *Body {
	bodyDef := box2d.NewB2BodyDef()
	switch bodyType {
	case STATIC:
		bodyDef.Type = box2d.B2BodyType.B2_staticBody
	case DYNAMIC:
		bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	case KINEMATIC:
		bodyDef.Type = box2d.B2BodyType.B2_kinematicBody
	}

	body := new(Body)
	body.world = world
	body.bodyDef = bodyDef

	return body
}

func (world *World) CreatePrismaticJoint(
	bodyA, bodyB *Body,
	anchorA, anchorB, axisA Vec2,
	angle float64,
	limited bool, min, max float64,
	motored bool, motorSpeed, motorForce float64,
	collide bool,
) *Joint { // TODO wtf how can i create joints in bytearena/box2d? there is no fucking single examples about joints.
	def := box2d.MakeB2PrismaticJointDef()

	def.BodyA = bodyA.b2body
	def.BodyB = bodyB.b2body

	def.LocalAnchorA = toBox2dVec2(anchorA)
	def.LocalAnchorB = toBox2dVec2(anchorB)
	def.LocalAxisA = toBox2dVec2(axisA)

	def.ReferenceAngle = angle

	def.EnableLimit = limited
	def.LowerTranslation = min
	def.UpperTranslation = max

	def.EnableMotor = motored
	def.MotorSpeed = motorSpeed
	def.MaxMotorForce = motorForce

	def.CollideConnected = collide

	b2joint := world.b2world.CreateJoint(&def) // pass pointer to joint def

	joint := new(Joint)
	joint.jointDef = &def
	joint.b2joint = b2joint

	return joint
}

func (body *Body) AddCircleFixture(density, friction, restitution, radius float64) {
	shape := box2d.MakeB2CircleShape()
	shape.SetRadius(radius)

	fixDef := box2d.MakeB2FixtureDef()
	fixDef.Density = density
	fixDef.Friction = friction
	fixDef.Restitution = restitution
	fixDef.Shape = &shape

	body.fixDefs = append(body.fixDefs, &fixDef)
}

func (body *Body) AddPolygonFixture(density, friction, restitution float64, vertices []Vec2) {
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

	body.fixDefs = append(body.fixDefs, &fixDef)
}

func (body *Body) Bake() {
	if body.b2body == nil {
		body.b2body = body.world.b2world.CreateBody(body.bodyDef)
	}

	for _, fixDef := range body.fixDefs {
		body.b2body.CreateFixtureFromDef(fixDef)
	}
}

// Destroy all fixtures
func (body *Body) Clear() {
	if body.b2body != nil {
		fixList := make([]*box2d.B2Fixture, 0)
		for fix := body.b2body.GetFixtureList(); fix != nil; fix = fix.GetNext() {
			fixList = append(fixList, fix)
		}

		for _, fix := range fixList {
			body.b2body.DestroyFixture(fix)
		}
	}

	body.fixDefs = nil
}

func (body *Body) GetPosition() (float64, float64) {
	var pos box2d.B2Vec2
	if body.b2body == nil {
		pos = body.bodyDef.Position
	} else {
		pos = body.b2body.GetPosition()
	}
	return pos.X, pos.Y
}

func (body *Body) GetAngle() float64 {
	if body.b2body == nil {
		return body.bodyDef.Angle
	} else {
		return body.b2body.GetAngle()
	}
}

func (body *Body) SetPosition(x, y float64) {
	if body.b2body == nil {
		body.bodyDef.Position = box2d.MakeB2Vec2(x, y)
	} else {
		body.b2body.SetTransform(box2d.MakeB2Vec2(x, y), 0)
	}
}

func (body *Body) SetLinearDamping(damp float64) {
	if body.b2body == nil {
		body.bodyDef.LinearDamping = damp
	} else {
		body.b2body.SetLinearDamping(damp)
	}
}

func (body *Body) ApplyForceToCenter(force Vec2) {
	if body.b2body == nil {
		return
	}
	body.b2body.ApplyForceToCenter(box2d.B2Vec2(force), true)
}

// Destroy body from world.
func (body *Body) Destroy() {
	body.Clear()
	b2world := body.b2body.GetWorld()
	b2world.DestroyBody(body.b2body)
	body.b2body = nil
}

func toBox2dVec2(vec Vec2) box2d.B2Vec2 {
	return box2d.B2Vec2(vec)
}
