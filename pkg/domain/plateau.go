package domain

var _ LandingZone = Plateau{}

func NewPlateau(x, y uint64) Plateau {
	return Plateau{height: y, width: x}
}

type Plateau struct {
	height, width uint64
}

func (p Plateau) Height() uint64 {
	return p.height
}

func (p Plateau) Width() uint64 {
	return p.width
}
