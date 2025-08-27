package game

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Game struct {
	Id      int         `json:"id"`
	TileMap [16][12]int `json:"tileMap"` // 0: empty, 1: white, 2: black , defined as [y][x]
	Players []*Player   `json:"players"`
	Bullets []*Bullet   `json:"bullets"`
}

var (
	roomID int
	mu     sync.Mutex
)

var mapLayout1 = [16][12]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 1, 1, 1, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 1, 1, 1, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 1, 1, 1, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 0},
	{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 0},
	{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func NewGame(player1Conn, player2Conn *websocket.Conn) *Game {
	mu.Lock()
	defer mu.Unlock()

	player1 := NewPlayer(1, player1Conn)
	player2 := NewPlayer(2, player2Conn)

	roomID++
	return &Game{
		Id:      roomID,
		TileMap: mapLayout1,
		Players: []*Player{player1, player2},
		Bullets: []*Bullet{},
	}
}

func (g *Game) Start() {
	SendGameInit(g)
	Player1Chan := make(chan Input)
	Player2Chan := make(chan Input)

	go g.Players[0].Listen(Player1Chan)
	go g.Players[1].Listen(Player2Chan)

	// tick := time.NewTicker(1 * time.Second)
	tick := time.NewTicker(40 * time.Millisecond)
	count := 0

	for {
		select {

		case player1Inputs := <-Player1Chan:
			fmt.Println("player 1 sent: ", player1Inputs, " at ", time.Now().Format("15:04:05.000"))

			player1 := g.Players[0]
			isBullet := player1.HandleInput(player1Inputs)
			if isBullet {
				bullet := NewBullet(player1.Id, player1.X, player1.Y, player1.Dir)
				g.Bullets = append(g.Bullets, bullet)
			}

		case player2Inputs := <-Player2Chan:
			fmt.Println("player 2 sent: ", player2Inputs, " at ", time.Now().Format("15:04:05.000"))

			player2 := g.Players[1]
			isBullet := player2.HandleInput(player2Inputs)
			if isBullet {
				bullet := NewBullet(player2.Id, player2.X, player2.Y, player2.Dir)
				g.Bullets = append(g.Bullets, bullet)
			}

		case <-tick.C:
			count++
			if SendTick(g, count) {
				fmt.Println("Game over!")
				SendGameEnd(g)
				return
			}
		}
	}
}

func (g *Game) MovePlayer(index int) bool {
	player := g.Players[index]
	if player.Vely == 0 && player.Velx == 0 {
		return false
	}
	// player.X += player.Velx
	// player.Y += player.Vely
	newX := player.X + player.Velx
	newY := player.Y + player.Vely

	ceilX := int(math.Ceil(float64(newX)))
	ceilY := int(math.Ceil(float64(newY)))

	floorX := int(math.Floor(float64(newX)))
	floorY := int(math.Floor(float64(newY)))

	if g.TileMap[floorY][floorX] != 0 && g.TileMap[floorY][floorX] != player.Id {
		if g.TileMap[ceilY][ceilX] != 0 && g.TileMap[ceilY][ceilX] != player.Id {
			player.X = float32(math.Round(float64(newX*100))) / 100
			player.Y = float32(math.Round(float64(newY*100))) / 100
			return true
		}
	}
	return false
}

func (g *Game) MoveBullet(index int) bool {
	bullet := g.Bullets[index]
	var xPos, yPos int
	switch bullet.Dir {
	case Left:
		bullet.X -= MaxBulletSpeed
		xPos = int(math.Floor(float64(bullet.X)))
		yPos = int(bullet.Y)
	case Right:
		bullet.X += MaxBulletSpeed
		xPos = int(math.Ceil(float64(bullet.X)))
		yPos = int(bullet.Y)
	case Up:
		bullet.Y -= MaxBulletSpeed
		yPos = int(math.Floor(float64(bullet.Y)))
		xPos = int(bullet.X)
	case Down:
		bullet.Y += MaxBulletSpeed
		yPos = int(math.Ceil(float64(bullet.Y)))
		xPos = int(bullet.X)
	}

	if g.TileMap[yPos][xPos] == 0 {
		g.Bullets[index].isActive = false
		return false
	}
	if g.TileMap[yPos][xPos] == bullet.Id {
		g.SwapTile(xPos, yPos)
	}
	return true
}

func (g *Game) SwapTile(xPos, yPos int) {
	if g.TileMap[yPos][xPos] == 1 {
		g.TileMap[yPos][xPos] = 2
	} else {
		g.TileMap[yPos][xPos] = 1
	}
}

func (g *Game) GetState() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Game) PrintState() {
	state, _ := g.GetState()
	fmt.Println(string(state))
}
