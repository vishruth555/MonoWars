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

const ()

func (p *Player) SetVelocity(dx, dy float32) {
	p.Velx = dx
	p.Vely = dy
}
func (p *Player) SetDirection() {
	switch {
	case p.Velx < 0:
		p.Dir = Left
	case p.Velx > 0:
		p.Dir = Right
	case p.Vely < 0:
		p.Dir = Up
	case p.Vely > 0:
		p.Dir = Down
	}
}

func (p *Player) HandleInput(input Input) bool {
	switch input.Type {
	case Move:
		p.SetVelocity(input.Dx, input.Dy)
		p.SetDirection()

	case Shoot:
		p.Bullets--
		return true
	}
	return false
}
