package game

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/vishruth555/MonoWars/game/utils"

	"github.com/gorilla/websocket"
)

type Lobby struct {
	mu      sync.Mutex
	Rooms   []*Game
	Players utils.Queue[*websocket.Conn]
}

func NewLobby() *Lobby {
	var queue utils.Queue[*websocket.Conn]
	var games []*Game
	return &Lobby{
		Rooms:   games,
		Players: queue,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (l *Lobby) AddPlayer(w http.ResponseWriter, r *http.Request) {
	player, _ := upgrader.Upgrade(w, r, nil)
	fmt.Println("Player connected: ", player.RemoteAddr())
	l.mu.Lock()
	if l.Players.IsEmpty() {
		l.Players.Enqueue(player)
		player.WriteJSON(Message{Type: "waiting"})
	} else {
		opponent := l.Players.Dequeue()
		fmt.Println("game start for", opponent.RemoteAddr(), " and ", player.RemoteAddr())
		game := NewGame(opponent, player)
		go game.Start()
		// game.AddPlayer(1, opponent)
		// game.AddPlayer(2, player)
		game.PrintState()
	}
	l.mu.Unlock()
}

// func (l *Lobby) AddRoom(id int) {
// 	l.Rooms[id] = NewGame()
// }
