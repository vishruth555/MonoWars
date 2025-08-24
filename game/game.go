package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Game struct {
	Id      int         `json:"id"`
	TileMap [16][12]int `json:"tileMap"` // 0: empty, 1: white, 2: black , defined as [y][x]
	Players []*Player   `json:"players"`
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
	}
}

func (g *Game) Start() {
	SendGameInit(g)
	Player1Chan := make(chan Input)
	Player2Chan := make(chan Input)

	go g.Players[0].Listen(Player1Chan)
	go g.Players[1].Listen(Player2Chan)

	tick := time.NewTicker(1 * time.Second)
	count := 0

	for {
		select {

		case player1Inputs := <-Player1Chan:
			g.Players[0].HandleInput(player1Inputs)
			fmt.Println("player 1 sent: ", player1Inputs, " at ", time.Now().Format("15:04:05.000"))

			// g.Players[0].SetVelocity()

		case player2Inputs := <-Player2Chan:
			g.Players[1].HandleInput(player2Inputs)
			fmt.Println("player 2 sent: ", player2Inputs, " at ", time.Now().Format("15:04:05.000"))
			fmt.Println("Player 2 Velocity: ", g.Players[1].Velx, g.Players[1].Vely)
			fmt.Println("Player 2 Direction: ", g.Players[1].Dir)

		case <-tick.C:
			count++
			SendTick(g, count)

		}

	}
}

func (g *Game) GetState() ([]byte, error) {
	return json.Marshal(g)
}

func (g *Game) PrintState() {
	state, _ := g.GetState()
	fmt.Println(string(state))
}
