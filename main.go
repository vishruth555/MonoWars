package main

import (
	"fmt"
	"log"
	"net/http"

	"MonoWars/game"
)

func main() {
	lobby := game.NewLobby()
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", lobby.AddPlayer)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
