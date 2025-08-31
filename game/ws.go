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
type GameEnd struct {
	Type   string `json:"type"`
	Winner int    `json:"winner"`
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
	Id      int     `json:"id"`
	XPos    float32 `json:"xPos"`
	YPos    float32 `json:"yPos"`
	Bullets int     `json:"bullets"`
}
type BulletData struct {
	Id       int     `json:"id"`
	BulletId int     `json:"bulletId"`
	XPos     float32 `json:"xPos"`
	YPos     float32 `json:"yPos"`
	State    string  `json:"state"`
}
type TileMapData struct {
	Id   int `json:"id"`
	XPos int `json:"xPos"`
	YPos int `json:"yPos"`
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
		Player1: PlayerData{Id: g.Players[0].Id, XPos: g.Players[0].X, YPos: g.Players[0].Y, Bullets: g.Players[0].Bullets},
		Player2: PlayerData{Id: g.Players[1].Id, XPos: g.Players[1].X, YPos: g.Players[1].Y, Bullets: g.Players[1].Bullets}}
	Broadcast(g, data)
}
func SendGameEnd(g *Game, id int) {
	data := GameEnd{
		Type:   "GameEnd",
		Winner: id,
	}
	Broadcast(g, data)
}

func SendTick(g *Game, tickCount int) {
	var diffs []Diff

	//player diffs
	for _, player := range g.Players {
		isChanged := player.Move(g)
		if isChanged {
			diffs = append(diffs, Diff{
				Entity: fmt.Sprintf("Player%dData", player.Id),
				Data:   PlayerData{Id: player.Id, XPos: player.X, YPos: player.Y, Bullets: player.Bullets},
			})
		}
	}
	//store old copy of tileMap
	tileMap := g.TileMap

	//bullet diffs
	for index, bullet := range g.Bullets {
		if !bullet.isActive {
			continue
		}
		isActive := bullet.Move(g)
		bulletData := BulletData{Id: bullet.Id, BulletId: index + 1, XPos: bullet.X, YPos: bullet.Y}
		if isActive {
			bulletData.State = "active"
		} else {
			bulletData.State = "expired"
		}
		diffs = append(diffs, Diff{
			Entity: "BulletData",
			Data:   bulletData,
		})
	}

	//tileMap diffs
	for i := 0; i < len(tileMap); i++ {
		for j := 0; j < len(tileMap[i]); j++ {
			if tileMap[i][j] != g.TileMap[i][j] {

				fmt.Printf("Changed at [%d][%d]: %d -> %d\n",
					i, j, tileMap[i][j], g.TileMap[i][j])

				diffs = append(diffs, Diff{
					Entity: "TileMapData",
					Data: TileMapData{
						Id:   g.TileMap[i][j],
						XPos: j,
						YPos: i,
					},
				})
			}
		}
	}

	data := Tick{
		Type: "Tick",
		Id:   tickCount,
		Diff: diffs,
	}
	Broadcast(g, data)
}
