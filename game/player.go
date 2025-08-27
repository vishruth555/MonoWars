package game

import (
	"fmt"

	"github.com/gorilla/websocket"
)

const MaxBulletCount = 10

type Player struct {
	Id      int
	Conn    *websocket.Conn
	X       float32
	Y       float32
	Bullets int
	Dir     Direction
	Velx    float32
	Vely    float32
	ay      float32
	ax      float32
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Type int

const (
	Move Type = iota
	Shoot
)

type Input struct {
	Type Type    `json:"type"`
	Dx   float32 `json:"dx"`
	Dy   float32 `json:"dy"`
}

func NewPlayer(id int, conn *websocket.Conn) *Player {
	var x, y float32
	var dir Direction

	if id == 2 {
		x = 8.0
		y = 3.0
		dir = Down
	} else {
		x = 3.0
		y = 12.0
		dir = Up
	}

	return &Player{
		Id:      id,
		Conn:    conn,
		X:       x,
		Y:       y,
		Bullets: MaxBulletCount,
		Dir:     dir,
	}
}

func (p *Player) Listen(ws chan Input) {
	var input Input
	for {

		err := p.Conn.ReadJSON(&input)
		if err != nil {
			//TODO handle this by adding an error channel and closing the game
			fmt.Println("Read error: ", err)
			break
		}
		ws <- input
	}
}

// func (p *Player) SetVelocity(dx, dy float32) {
// 	p.Velx = dx
// 	p.Vely = dy
// }

func (p *Player) setAcceleration(dx, dy float32) {
	p.ay = dy
	p.ax = dx
}
func (p *Player) SetDirection() {
	switch {
	case p.ax < 0:
		p.Dir = Left
	case p.ax > 0:
		p.Dir = Right
	case p.ay < 0:
		p.Dir = Up
	case p.ay > 0:
		p.Dir = Down
	}
}

func (p *Player) HandleInput(input Input) bool {
	switch input.Type {
	case Move:
		p.setAcceleration(input.Dx, input.Dy)
		p.SetDirection()

	case Shoot:
		// p.Bullets--
		return true
	}
	return false
}

const (
	friction = 0.2
	MaxSpeed = 0.2
)

func (p *Player) ApplyVelocity(dt float32) {
	p.Velx += p.ax * dt
	if p.Velx > MaxSpeed {
		p.Velx = MaxSpeed
	} else if p.Velx < -MaxSpeed {
		p.Velx = -MaxSpeed
	}

	p.Vely += p.ay * dt
	if p.Vely > MaxSpeed {
		p.Vely = MaxSpeed
	} else if p.Vely < -MaxSpeed {
		p.Vely = -MaxSpeed
	}

}

func (p *Player) ApplyFriction(dt float32) {

	if p.Velx > 0 {
		p.Velx -= friction * dt
		if p.Velx < 0 {
			p.Velx = 0
		}
	} else if p.Velx < 0 {
		p.Velx += friction * dt
		if p.Velx > 0 {
			p.Velx = 0
		}
	}

	if p.Vely > 0 {
		p.Vely -= friction * dt
		if p.Vely < 0 {
			p.Vely = 0
		}
	} else if p.Vely < 0 {
		p.Vely += friction * dt
		if p.Vely > 0 {
			p.Vely = 0
		}
	}

}

func (p *Player) Stop() {
	p.Velx = 0
	p.Vely = 0
}
