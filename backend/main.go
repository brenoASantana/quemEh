package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed dist/*
var distFiles embed.FS

type GameState string

const (
	StateLobby   GameState = "LOBBY"
	StateAnswer  GameState = "ANSWERING"
	StateVoting  GameState = "VOTING"
	StateResults GameState = "RESULTS"
)

type Player struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	IsHost bool            `json:"isHost"`
	Score  int             `json:"score"`
	Conn   *websocket.Conn `json:"-"`
	Room   *Room           `json:"-"`
}

type Room struct {
	ID              string
	Players         map[string]*Player
	State           GameState
	Question        string
	Answers         map[string]string
	VotedGuesses    map[string]string
	ShuffledAnswers []AnswerWithIndex
	Broadcast       chan Message
	Mutex           sync.RWMutex
	QuestionIndex   int
}

type AnswerWithIndex struct {
	Text      string `json:"text"`
	RealID    string `json:"-"`
	DisplayID int    `json:"id"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type PublicGameState struct {
	Players         []*Player         `json:"players"`
	State           GameState         `json:"state"`
	Question        string            `json:"question"`
	Answers         []AnswerWithIndex `json:"answers,omitempty"`
	VotedGuesses    map[string]string `json:"votedGuesses,omitempty"`
	LastRoundResult map[string]int    `json:"lastRoundResult,omitempty"`
}

var (
	rooms     = make(map[string]*Room)
	roomMutex sync.RWMutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	questions = []string{
		"Se você fosse morrer em 5 dias e só você soubesse disso, o que você faria?",
		"Qual é o seu maior medo em relação ao seu futuro?",
		"Se tivesse um superpoder, qual seria e por quê?",
		"Qual é o segredo mais pesado que você nunca contou a ninguém?",
		"O que você faria com 1 milhão de reais?",
		"Qual pessoa famosa você mais gostaria de jantar?",
		"Se pudesse voltar no tempo, o que mudaria?",
		"Qual é a coisa mais estranha que você já fez?",
		"O que você mais ama na vida?",
		"Se fosse invisível por um dia, o que faria?",
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/ws", handleConnections)

	distFS, err := fs.Sub(distFiles, "dist")
	if err != nil {
		log.Fatal("Erro ao configurar arquivos embutidos:", err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFileFS(w, r, distFS, "index.html")
			return
		}

		file, err := distFS.Open(r.URL.Path[1:])
		if err == nil {
			file.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFileFS(w, r, distFS, "index.html")
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer upgrade:", err)
		return
	}

	roomID := r.URL.Query().Get("room")
	playerName := r.URL.Query().Get("name")

	if roomID == "" || playerName == "" {
		ws.Close()
		return
	}

	player := &Player{
		ID:    fmt.Sprintf("%d_%d", time.Now().UnixNano(), rand.Int63()),
		Name:  playerName,
		Conn:  ws,
		Score: 0,
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
	roomMutex.Lock()
	defer roomMutex.Unlock()

	if room, exists := rooms[id]; exists {
		return room
	}

	newRoom := &Room{
		ID:              id,
		Players:         make(map[string]*Player),
		State:           StateLobby,
		Broadcast:       make(chan Message, 100),
		Answers:         make(map[string]string),
		VotedGuesses:    make(map[string]string),
		ShuffledAnswers: []AnswerWithIndex{},
		QuestionIndex:   0,
	}

	rooms[id] = newRoom
	go runRoom(newRoom)

	log.Printf("Nova sala criada: %s\n", id)
	return newRoom
}

func handlePlayerMessages(p *Player) {
	defer func() {
		p.Conn.Close()
		p.Room.Mutex.Lock()
		delete(p.Room.Players, p.ID)
		p.Room.Mutex.Unlock()
		broadcastGameState(p.Room)
	}()

	for {
		var msg Message
		err := p.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		switch msg.Type {
		case "START_GAME":
			if p.IsHost {
				startGame(p.Room)
			}
		case "SUBMIT_ANSWER":
			if answer, ok := msg.Payload.(string); ok && answer != "" {
				handleAnswer(p, answer)
			}
		case "SUBMIT_GUESS":
			if guessData, ok := msg.Payload.(map[string]interface{}); ok {
				if answerId, exists := guessData["answerId"]; exists {
					handleGuess(p, answerId.(string))
				}
			}
		case "NEXT_ROUND":
			if p.IsHost {
				nextRound(p.Room)
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
	r.QuestionIndex = 0
	r.State = StateAnswer
	r.Question = questions[0]
	r.Answers = make(map[string]string)
	r.VotedGuesses = make(map[string]string)
	r.ShuffledAnswers = []AnswerWithIndex{}
	r.Mutex.Unlock()

	broadcastGameState(r)
}

func handleAnswer(p *Player, answer string) {
	p.Room.Mutex.Lock()
	p.Room.Answers[p.ID] = answer

	allAnswered := len(p.Room.Answers) == len(p.Room.Players)

	if allAnswered {
		p.Room.State = StateVoting
		// Embaralha as respostas
		shuffledAnswers := shuffleAnswers(p.Room.Answers)
		p.Room.ShuffledAnswers = shuffledAnswers
		p.Room.VotedGuesses = make(map[string]string) // Reset de votos
	}
	p.Room.Mutex.Unlock()

	broadcastGameState(p.Room)
}

func handleGuess(p *Player, answerId string) {
	p.Room.Mutex.Lock()
	p.Room.VotedGuesses[p.ID] = answerId
	allVoted := len(p.Room.VotedGuesses) == len(p.Room.Players)
	p.Room.Mutex.Unlock()

	if allVoted {
		// Calcular scores
		calculateScores(p.Room)
		p.Room.Mutex.Lock()
		p.Room.State = StateResults
		p.Room.Mutex.Unlock()
	}

	broadcastGameState(p.Room)
}

func nextRound(r *Room) {
	r.Mutex.Lock()
	r.QuestionIndex++
	if r.QuestionIndex >= len(questions) {
		r.QuestionIndex = 0
	}

	r.State = StateAnswer
	r.Question = questions[r.QuestionIndex]
	r.Answers = make(map[string]string)
	r.VotedGuesses = make(map[string]string)
	r.ShuffledAnswers = []AnswerWithIndex{}
	r.Mutex.Unlock()

	broadcastGameState(r)
}

func shuffleAnswers(answers map[string]string) []AnswerWithIndex {
	shuffled := []AnswerWithIndex{}
	ids := []string{}

	for id, text := range answers {
		ids = append(ids, id)
		shuffled = append(shuffled, AnswerWithIndex{
			Text:   text,
			RealID: id,
		})
	}

	// Fisher-Yates shuffle
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	// Atribui DisplayIDs sequenciais
	for i := range shuffled {
		shuffled[i].DisplayID = i
	}

	return shuffled
}

func calculateScores(r *Room) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	correct := make(map[int]string)
	for _, ans := range r.ShuffledAnswers {
		correct[ans.DisplayID] = ans.RealID
	}

	for playerId, guessDisplayId := range r.VotedGuesses {
		var displayId int
		fmt.Sscanf(guessDisplayId, "%d", &displayId)
		if guessId, exists := correct[displayId]; exists && guessId == correct[displayId] {
			if player, exists := r.Players[playerId]; exists {
				player.Score += 10
			}
		}
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
	if r.State != StateVoting && r.State != StateResults {
		answersList = []AnswerWithIndex{}
	}

	var lastRoundResult map[string]int
	if r.State == StateResults {
		lastRoundResult = make(map[string]int)
		for _, ans := range r.ShuffledAnswers {
			correct := 0
			if guess, exists := r.VotedGuesses[ans.RealID]; exists && guess == fmt.Sprintf("%d", ans.DisplayID) {
				correct = 1
			}
			lastRoundResult[ans.RealID] = correct
		}
	}

	msg := Message{
		Type: "GAME_STATE",
		Payload: PublicGameState{
			Players:         playersList,
			State:           r.State,
			Question:        r.Question,
			Answers:         answersList,
			VotedGuesses:    r.VotedGuesses,
			LastRoundResult: lastRoundResult,
		},
	}

	go func() {
		select {
		case r.Broadcast <- msg:
		case <-time.After(1 * time.Second):
			log.Println("Timeout ao enviar mensagem")
		}
	}()
}
