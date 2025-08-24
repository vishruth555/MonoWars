package game

import (
	"fmt"
)

type Message struct {
	Type string `json:"type"`
}

type GameInit struct {
	Type    string      `json:"type"`
	TileMap [16][12]int `json:"tileMap"`
	Player1 PlayerData  `json:"player1Data"`
	Player2 PlayerData  `json:"player2Data"`
}

type Tick struct {
	Type string `json:"type"`
	Id   int    `json:"id"`
	Diff []Diff `json:"diff"`
}

type Diff struct {
	Entity string `json:"entity"`
	Data   any    `json:"data"`
}

type PlayerData struct {
	Id   int     `json:"id"`
	XPos float32 `json:"xPos"`
	YPos float32 `json:"yPos"`
}

func Broadcast(g *Game, data any) {
	for _, player := range g.Players {
		err := player.Conn.WriteJSON(data)
		if err != nil {
			fmt.Println("error sending data to ", player.Conn.RemoteAddr(), err)
		}
	}
}

func SendGameInit(g *Game) {
	data := GameInit{
		Type:    "GameStart",
		TileMap: g.TileMap,
		Player1: PlayerData{Id: g.Players[0].Id, XPos: g.Players[0].X, YPos: g.Players[0].Y},
		Player2: PlayerData{Id: g.Players[1].Id, XPos: g.Players[1].X, YPos: g.Players[1].Y}}
	Broadcast(g, data)
}

func SendTick(g *Game, tickCount int) {
	var diffs []Diff

	for _, player := range g.Players {
		if player.Vely == 0 && player.Velx == 0 {
			continue
		}
		player.X += player.Velx
		player.Y += player.Vely
		diffs = append(diffs, Diff{
			Entity: fmt.Sprintf("player%dData", player.Id),
			Data:   PlayerData{Id: player.Id, XPos: player.X, YPos: player.Y},
		})
	}

	data := Tick{
		Type: "Tick",
		Id:   tickCount,
		Diff: diffs,
	}
	Broadcast(g, data)
}
