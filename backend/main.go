package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)


var (
	rooms     = make(map[string]*Room)
	roomMutex sync.RWMutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	questionsPath string
)

func main() {
	loadQuestions()

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/ws", handleConnections)

	// Encontrar o diretório de arquivos estáticos
	// Primeiro tenta ./frontend/dist (quando rodado da raiz do projeto)
	// Depois tenta ../frontend/dist (quando rodado da pasta backend)
	var staticDir string
	if _, err := os.Stat("./frontend/dist"); err == nil {
		staticDir = "./frontend/dist"
	} else if _, err := os.Stat("../frontend/dist"); err == nil {
		staticDir = "../frontend/dist"
	} else {
		log.Fatal("Diretório frontend/dist não encontrado")
	}

	log.Printf("Servindo arquivos estáticos de: %s\n", staticDir)

	fileServer := http.FileServer(http.Dir(staticDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Se o caminho não for a raiz e tiver uma extensão, sirva o arquivo estático
		if r.URL.Path != "/" && filepath.Ext(r.URL.Path) != "" {
			fileServer.ServeHTTP(w, r)
			return
		}
		// Caso contrário, sirva o index.html (para o roteamento do React)
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadQuestions() {
	paths := []string{"./backend/questions.json", "./questions.json"}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			if err := repo.LoadQuestions(path); err != nil {
				log.Fatal("Erro ao carregar perguntas:", err)
			}
			questionsPath = path
			log.Printf("Perguntas carregadas: %d\n", len(repo.AllQuestions))
			return
		}
	}

	log.Fatal("Arquivo questions.json não encontrado")
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
		VotedGuesses:    make(map[string]map[int]string),
		ShuffledAnswers: []AnswerWithIndex{},
		QuestionIndex:   0,
		QuestionDeck:    []string{},
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
				handleGuess(p, guessData)
			}
		case "SHOW_RESULTS":
			if p.IsHost {
				showResults(p.Room)
			}
		case "NEXT_ROUND":
			if p.IsHost {
				nextRound(p.Room)
			}
		case "RESET_GAME":
			if p.IsHost {
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
		shuffledAnswers := shuffleAnswers(p.Room.Answers)
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
	if r.State != StateVoting && r.State != StateResults {
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

	go func() {
		select {
		case r.Broadcast <- msg:
		case <-time.After(1 * time.Second):
			log.Println("Timeout ao enviar mensagem")
		}
	}()
}

// Reinicia a partida e zera os pontos de todos
func resetGame(r *Room) {
	r.Mutex.Lock()
	for _, p := range r.Players {
		p.Score = 0
	}
	r.QuestionDeck = repo.GetShuffledQuestions()
	r.QuestionIndex = 0
	r.State = StateAnswer
	if len(r.QuestionDeck) == 0 {
		r.Question = "Sem perguntas cadastradas"
	} else {
		r.Question = r.QuestionDeck[0]
		       case "START_GAME":
			       if p.IsHost {
				       p.Room.StartGame()
			       }
		       case "SUBMIT_ANSWER":
			       if answer, ok := msg.Payload.(string); ok && answer != "" {
				       p.Room.HandleAnswer(p, answer)
			       }
		       case "SUBMIT_GUESS":
			       if guessData, ok := msg.Payload.(map[string]interface{}); ok {
				       p.Room.HandleGuess(p, guessData)
			       }
		       case "SHOW_RESULTS":
			       if p.IsHost {
				       p.Room.ShowResults()
			       }
		       case "NEXT_ROUND":
			       if p.IsHost {
				       p.Room.NextRound()
			       }
			Text:   text,
