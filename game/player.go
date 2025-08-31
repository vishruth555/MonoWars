package game

import (
	"fmt"
	"time"

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

const (
	friction           = 0.2
	MaxSpeed           = 0.2
	accelerationFactor = 0.75
)

func (p *Player) setAcceleration(dx, dy float32) {
	p.ay = dy * accelerationFactor
	p.ax = dx * accelerationFactor
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

func (player *Player) Move(g *Game) bool {
	dt := float32(g.tickRate) / float32(time.Second)
	player.ApplyVelocity(dt)
	player.ApplyFriction(dt)

	if player.Vely == 0 && player.Velx == 0 {
		return false
	}

	newX := player.X + player.Velx
	newY := player.Y + player.Vely

	x1 := int(newX)
	y1 := int(newY)

	x2 := int(newX + 0.8)
	y2 := int(newY + 0.8)

	if g.TileMap[y1][x1] != 0 && g.TileMap[y1][x1] != player.Id {
		if g.TileMap[y2][x2] != 0 && g.TileMap[y2][x2] != player.Id {
			player.X = newX
			player.Y = newY
			return true
		}
	}

	// ceilX := int(math.Ceil(float64(newX)))
	// ceilY := int(math.Ceil(float64(newY)))

	// floorX := int(math.Floor(float64(newX)))
	// floorY := int(math.Floor(float64(newY)))

	// if g.TileMap[floorY][floorX] != 0 && g.TileMap[floorY][floorX] != player.Id {
	// 	if g.TileMap[ceilY][ceilX] != 0 && g.TileMap[ceilY][ceilX] != player.Id {
	// 		player.X = float32(math.Round(float64(newX*100))) / 100
	// 		player.Y = float32(math.Round(float64(newY*100))) / 100
	// 		return true
	// 	}
	// }
	player.Stop()
	return false
}
