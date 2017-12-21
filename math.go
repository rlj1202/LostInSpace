package lostinspace

import (
	"math"
)

type Vec2 struct {
	X, Y float64
}
type Vec3 struct {
	X, Y, Z float64
}
type Vec4 struct {
	X, Y, Z, W float64
}
type Mat4 [16]float64

func Distance(a, b *Vec2) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

func DistanceP(a, b *Vec2) float64 {
	return math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2)
}

func (d *Vec2) Add(a, b *Vec2) {
	d.X = a.X + b.X
	d.Y = a.Y + b.Y
}

func (d *Vec2) Sub(a, b *Vec2) {
	d.X = a.X - b.X
	d.Y = a.Y - b.Y
}

func (d *Vec2) Mul(a *Vec2, b float64) {
	d.X = a.X * b
	d.Y = a.Y * b
}

func (d *Vec2) Div(a *Vec2, b float64) {
	d.X = a.X / b
	d.Y = a.Y / b
}

func (d *Vec2) DotProduct(a, b *Vec2) {
	d.X = a.X * b.X
	d.Y = a.Y * b.Y
}

func CrossProduct(a, b *Vec2) float64 {
	return a.X*b.Y - b.X*a.Y
}

func PerlinNoiseImproved(perm [256]int, x, y, z float64) float64 {
	p := append(perm[:], perm[:]...)

	xi := int(x) & 0xff // Coordinate is repeated
	yi := int(y) & 0xff
	zi := int(z) & 0xff

	xf := x - math.Floor(x)
	yf := y - math.Floor(y)
	zf := z - math.Floor(z)

	u := fade(xf)
	v := fade(yf)
	w := fade(zf)

	aaa := p[p[p[xi]+yi]+zi]
	aba := p[p[p[xi]+yi+1]+zi]
	aab := p[p[p[xi]+yi]+zi+1]
	abb := p[p[p[xi]+yi+1]+zi+1]
	baa := p[p[p[xi+1]+yi]+zi]
	bba := p[p[p[xi+1]+yi+1]+zi]
	bab := p[p[p[xi+1]+yi]+zi+1]
	bbb := p[p[p[xi+1]+yi+1]+zi+1]

	x1 := lerp(grad(aaa, xf, yf, zf), grad(baa, xf-1, yf, zf), u)
	x2 := lerp(grad(aba, xf, yf-1, zf), grad(bba, xf-1, yf-1, zf), u)
	y1 := lerp(x1, x2, v)

	x1 = lerp(grad(aab, xf, yf, zf-1), grad(bab, xf-1, yf, zf-1), u)
	x2 = lerp(grad(abb, xf, yf-1, zf-1), grad(bbb, xf-1, yf-1, zf-1), u)
	y2 := lerp(x1, x2, v)

	return (lerp(y1, y2, w) + 1.0) / 2.0
}

func grad(hash int, x, y, z float64) float64 {
	switch hash & 0xf {
	case 0x0:
		return x + y
	case 0x1:
		return -x + y
	case 0x2:
		return x - y
	case 0x3:
		return -x - y
	case 0x4:
		return x + z
	case 0x5:
		return -x + z
	case 0x6:
		return x - z
	case 0x7:
		return -x - z
	case 0x8:
		return y + z
	case 0x9:
		return -y + z
	case 0xa:
		return y - z
	case 0xb:
		return -y - z
	case 0xc:
		return y + x
	case 0xd:
		return -y + z
	case 0xe:
		return y - x
	case 0xf:
		return -y - z
	default:
		return 0
	}
}

func lerp(a, b, x float64) float64 {
	return a + x*(b-a)
}

// result = 6t^5 - 15t^4 + 10t^3
func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}
