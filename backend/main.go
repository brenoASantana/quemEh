package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	rooms     = make(map[string]*Room)
	roomMutex sync.RWMutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	loadQuestions()
	var staticDir string
	paths := []string{
		"frontend/dist",
		"../frontend/dist",
		"./frontend/dist",
		"../../frontend/dist",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			staticDir = path
			break
		}
	}
	if staticDir == "" {
		panic("Diretório frontend/dist não encontrado")
	}
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.Contains(r.URL.Path, ".") {
			http.ServeFile(w, r, filepath.Join(staticDir, r.URL.Path))
			return
		}
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})
	log.Printf("🎮 Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadQuestions() {
	paths := []string{"questions.json", "backend/questions.json", "./backend/questions.json"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			repo.LoadQuestions(path)
			return
		}
	}
	panic("Arquivo questions.json não encontrado")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	roomID := r.URL.Query().Get("room")
	playerName := r.URL.Query().Get("name")
	if roomID == "" || playerName == "" {
		conn.Close()
		return
	}
	player := &Player{
		ID:   fmt.Sprintf("%d_%s", time.Now().UnixNano(), playerName),
		Name: playerName,
		Conn: conn,
	}
	room := getOrCreateRoom(roomID)
	player.Room = room
	room.Mutex.Lock()
	if len(room.Players) == 0 {
		player.IsHost = true
	}
	room.Players[player.ID] = player
	room.Mutex.Unlock()
	go handlePlayerMessages(player)
	broadcastGameState(room)
}

func getOrCreateRoom(id string) *Room {
	roomMutex.RLock()
	r, exists := rooms[id]
	roomMutex.RUnlock()
	if !exists {
		r = &Room{
			ID:              id,
			Players:         make(map[string]*Player),
			State:           StateLobby,
			Broadcast:       make(chan Message, 100),
			Answers:         make(map[string]string),
			VotedGuesses:    make(map[string]map[int]string),
			ShuffledAnswers: []AnswerWithIndex{},
			QuestionDeck:    []string{},
		}
		roomMutex.Lock()
		rooms[id] = r
		roomMutex.Unlock()
		go runRoom(r)
	}
	return r
}

func handlePlayerMessages(p *Player) {
	defer func() {
		if p.Room != nil {
			p.Room.Mutex.Lock()
			delete(p.Room.Players, p.ID)
			p.Room.Mutex.Unlock()
		}
		p.Conn.Close()
	}()
	for {
		var msg Message
		if err := p.Conn.ReadJSON(&msg); err != nil {
			break
		}
		switch msg.Type {
		case "START_GAME":
			if p.Room != nil && p.IsHost {
				p.Room.StartGame()
			}
		case "SUBMIT_ANSWER":
			if p.Room != nil {
				if answer, ok := msg.Payload.(string); ok {
					p.Room.HandleAnswer(p, answer)
				}
			}
		case "SUBMIT_GUESS":
			if p.Room != nil {
				if guessData, ok := msg.Payload.(map[string]interface{}); ok {
					p.Room.HandleGuess(p, guessData)
				}
			}
		case "SHOW_RESULTS":
			if p.Room != nil && p.IsHost {
				p.Room.ShowResults()
			}
		case "NEXT_ROUND":
			if p.Room != nil && p.IsHost {
				p.Room.NextRound()
			}
		case "RESET_GAME":
			if p.Room != nil && p.IsHost {
				resetGame(p.Room)
			}
		}
	}
}

func runRoom(r *Room) {
	for msg := range r.Broadcast {
		r.Mutex.RLock()
		players := make([]*Player, 0, len(r.Players))
		for _, p := range r.Players {
			players = append(players, p)
		}
		r.Mutex.RUnlock()

		for _, p := range players {
			err := p.Conn.WriteJSON(msg)
			if err != nil {
				p.Conn.Close()
				r.Mutex.Lock()
				delete(r.Players, p.ID)
				r.Mutex.Unlock()
			}
		}
	}
}

func startGame(r *Room) {
	r.Mutex.Lock()
	r.QuestionDeck = repo.GetShuffledQuestions()
	r.QuestionIndex = 0
	r.State = StateAnswer
	if len(r.QuestionDeck) == 0 {
		r.Question = "Sem perguntas cadastradas"
	} else {
		r.Question = r.QuestionDeck[0]
	}
	r.Answers = make(map[string]string)
	r.VotedGuesses = make(map[string]map[int]string)
	r.ShuffledAnswers = []AnswerWithIndex{}
	r.AnswersShuffled = false
	r.Mutex.Unlock()

	broadcastGameState(r)
}

func handleAnswer(p *Player, answer string) {
	p.Room.Mutex.Lock()

	if _, alreadyAnswered := p.Room.Answers[p.ID]; alreadyAnswered {
		p.Room.Mutex.Unlock()
		return
	}

	if answer == "" {
		p.Room.Mutex.Unlock()
		return
	}

	p.Room.Answers[p.ID] = answer

	allAnswered := len(p.Room.Answers) == len(p.Room.Players)

	if allAnswered {
		p.Room.State = StateVoting
		shuffledAnswers := ShuffleAnswers(p.Room.Answers)
		p.Room.ShuffledAnswers = shuffledAnswers
		p.Room.VotedGuesses = make(map[string]map[int]string)
	}
}

func broadcastGameState(r *Room) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	playersList := make([]*Player, 0, len(r.Players))
	for _, p := range r.Players {
		playersList = append(playersList, p)
	}

	answersList := r.ShuffledAnswers
	if r.State != StateVoting && r.State != StateResults && r.State != "GAME_OVER" {
		answersList = []AnswerWithIndex{}
	}

	var lastRoundResult map[string]int
	if r.State == StateResults {
		lastRoundResult = make(map[string]int)

		// Criar mapa de DisplayID -> RealID para facilitar lookup
		displayToReal := make(map[int]string)
		for _, ans := range r.ShuffledAnswers {
			displayToReal[ans.DisplayID] = ans.RealID
		}

		// Contar votos corretos por jogador
		for playerId := range r.Players {
			correctVotes := 0
			if playerVotes, exists := r.VotedGuesses[playerId]; exists {
				for answerDisplayId, guessedPlayerId := range playerVotes {
					if correctPlayerId, ok := displayToReal[answerDisplayId]; ok && correctPlayerId == guessedPlayerId {
						correctVotes++
					}
				}
			}
			lastRoundResult[playerId] = correctVotes
		}
	}

	msg := Message{
		Type: "GAME_STATE",
		Payload: PublicGameState{
			Players:         playersList,
			State:           r.State,
			Question:        r.Question,
			AnswersCount:    len(r.Answers),
			Answers:         answersList,
			VotedGuesses:    r.VotedGuesses,
			LastRoundResult: lastRoundResult,
		},
	}

	select {
	case r.Broadcast <- msg:
	default:
	}
}

func resetGame(r *Room) {
	r.Mutex.Lock()
	for _, p := range r.Players {
		p.Score = 0
	}
	r.State = StateLobby
	r.QuestionDeck = repo.GetShuffledQuestions()
	r.QuestionIndex = 0
	r.Answers = make(map[string]string)
	r.VotedGuesses = make(map[string]map[int]string)
	r.ShuffledAnswers = []AnswerWithIndex{}
	r.AnswersShuffled = false
	r.Mutex.Unlock()
	broadcastGameState(r)
}
