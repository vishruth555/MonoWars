package game

import "math"

const MaxBulletSpeed = 0.5 //max is 1

type Bullet struct {
	Id  int
	X   float32
	Y   float32
	Dir Direction
}

func NewBullet(id int, x float32, y float32, dir Direction) *Bullet {

	return &Bullet{
		Id:  id,
		X:   float32(math.Round(float64(x))),
		Y:   float32(math.Round(float64(y))),
		Dir: dir,
	}
}
