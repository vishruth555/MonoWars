package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Game struct {
	Id       int         `json:"id"`
	TileMap  [16][12]int `json:"tileMap"` // 0: empty, 1: white, 2: black , defined as [y][x]
	Players  []*Player   `json:"players"`
	Bullets  []*Bullet   `json:"bullets"`
	tickRate time.Duration
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
		Id:       roomID,
		TileMap:  mapLayout1,
		Players:  []*Player{player1, player2},
		Bullets:  []*Bullet{},
		tickRate: 40 * time.Millisecond,
	}
}

func (g *Game) Start() {
	SendGameInit(g)
	Player1Chan := make(chan Input)
	Player2Chan := make(chan Input)

	go g.Players[0].Listen(Player1Chan)
	go g.Players[1].Listen(Player2Chan)

	// tick := time.NewTicker(1 * time.Second)
	tick := time.NewTicker(g.tickRate)
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
			SendTick(g, count)

			if g.isGameOver() {
				fmt.Println("Game over!")
				return
			}

		}
	}
}

func (g *Game) SwapTile(xPos, yPos int) {
	if g.TileMap[yPos][xPos] == 1 {
		g.TileMap[yPos][xPos] = 2
	} else {
		g.TileMap[yPos][xPos] = 1
	}
}

func (g *Game) isGameOver() bool {

	for _, bullet := range g.Bullets {
		var dx, dy int
		switch bullet.Id {
		case 1:
			dx = int(bullet.X) - int(g.Players[1].X)
			dy = int(bullet.Y) - int(g.Players[1].Y)
			if dx == 0 && dy == 0 {
				//game over for player2
				SendGameEnd(g, 1)
				return true
			}
		case 2:
			dx = int(bullet.X) - int(g.Players[0].X)
			dy = int(bullet.Y) - int(g.Players[0].Y)
			if dx == 0 && dy == 0 {
				//game over for player1
				SendGameEnd(g, 2)
				return true
			}
		}
	}
	return false
}

func (g *Game) GetState() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Game) PrintState() {
	state, _ := g.GetState()
	fmt.Println(string(state))
}
